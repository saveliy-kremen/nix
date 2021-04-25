package fb

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo"
	"gorm.io/gorm"

	auth "../auth"
	db "../db"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

var (
	oauthConf = &oauth2.Config{
		ClientID:     "436634540735904",
		ClientSecret: "8423e6763aa22e40525b783674c396eb",
		RedirectURL:  "https://localhost:3000/fb_oauth2callback",
		Scopes:       []string{},
		Endpoint:     facebook.Endpoint,
	}
	oauthStateString = "localhost"
)

func HandleFacebookLogin(c echo.Context) error {
	Url, err := url.Parse(oauthConf.Endpoint.AuthURL)
	if err != nil {
		log.Fatal("Parse: ", err)
	}
	parameters := url.Values{}
	parameters.Add("client_id", oauthConf.ClientID)
	parameters.Add("scope", strings.Join(oauthConf.Scopes, " "))
	parameters.Add("redirect_uri", oauthConf.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", oauthStateString)
	Url.RawQuery = parameters.Encode()
	url := Url.String()
	http.Redirect(c.Response(), c.Request(), url, http.StatusTemporaryRedirect)
	return nil
}

func HandleFacebookCallback(c echo.Context) error {
	state := c.Request().FormValue("state")
	if state != oauthStateString {
		resp := fmt.Sprintf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		return c.String(http.StatusBadRequest, resp)
	}

	code := c.Request().FormValue("code")

	token, err := oauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		resp := fmt.Sprintf("oauthConf.Exchange() failed with '%s'\n", err)
		return c.String(http.StatusBadRequest, resp)
	}

	resp, err := http.Get("https://graph.facebook.com/me?access_token=" +
		url.QueryEscape(token.AccessToken))
	if err != nil {
		resp := fmt.Sprintf("Get: %s\n", err)
		return c.String(http.StatusBadRequest, resp)
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		resp := fmt.Sprintf("ReadAll: %s\n", err)
		return c.String(http.StatusBadRequest, resp)
	}

	var result struct {
		Id   string
		Name string
	}
	json.Unmarshal(response, &result)

	user := db.User{}
	userResult := db.DB.Where("fb_id = ?", result.Id).First(&user)
	if errors.Is(userResult.Error, gorm.ErrRecordNotFound) {
		user.Name = result.Name
		user.FbID = result.Id
		db.DB.Create(&user)
	}

	userToken := auth.CreateToken(user.Id, 1)
	res := fmt.Sprintf("Name: %s\nToken: %s", result.Name, userToken)
	return c.String(http.StatusOK, res)
}
