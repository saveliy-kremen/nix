package fb

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	auth "../auth"
	db "../db"
	"gorm.io/gorm"

	"github.com/dghubble/oauth1"
	"github.com/dghubble/oauth1/twitter"
	"github.com/labstack/echo"
)

var (
	oauthConf = oauth1.Config{
		ConsumerKey:    "rYOjptLwc87cvBcCTF75Fj8DR",
		ConsumerSecret: "kSPPR59TQoC1LMsV6L3l1VGonjzvfF1ZlHoJoFed0fupVf9AYp",
		CallbackURL:    "https://localhost:3000/tw_oauth2callback",
		Endpoint:       twitter.AuthorizeEndpoint,
	}
	secret string
)

func HandleTwitterLogin(c echo.Context) error {
	requestToken, requestSecret, err := oauthConf.RequestToken()
	if err != nil {
		resp := fmt.Sprintf("Get: %s\n", err)
		return c.String(http.StatusBadRequest, resp)
	}
	secret = requestSecret
	authorizationURL, err := oauthConf.AuthorizationURL(requestToken)
	if err != nil {
		resp := fmt.Sprintf("Get: %s\n", err)
		return c.String(http.StatusBadRequest, resp)
	}
	http.Redirect(c.Response(), c.Request(), authorizationURL.String(), http.StatusFound)
	return nil
}

func HandleTwitterCallback(c echo.Context) error {
	requestToken, verifier, err := oauth1.ParseAuthorizationCallback(c.Request())
	if err != nil {
		resp := fmt.Sprintf("Get: %s\n", err)
		return c.String(http.StatusBadRequest, resp)
	}
	accessToken, accessSecret, err := oauthConf.AccessToken(requestToken, secret, verifier)
	if err != nil {
		resp := fmt.Sprintf("Get: %s\n", err)
		return c.String(http.StatusBadRequest, resp)
	}
	twToken := oauth1.NewToken(accessToken, accessSecret)

	token := oauth1.NewToken(twToken.Token, twToken.TokenSecret)
	httpClient := oauthConf.Client(oauth1.NoContext, token)

	path := "https://api.twitter.com/1.1/account/verify_credentials.json"
	resp, err := httpClient.Get(path)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		resp := fmt.Sprintf("Get: %s\n", err)
		return c.String(http.StatusBadRequest, resp)
	}
	var result struct {
		Id   string `json:"id_str"`
		Name string
	}
	json.Unmarshal(response, &result)

	user := db.User{}
	userResult := db.DB.Where("tw_id = ?", result.Id).First(&user)
	if errors.Is(userResult.Error, gorm.ErrRecordNotFound) {
		user.Name = result.Name
		user.TwID = result.Id
		db.DB.Create(&user)
	}

	userToken := auth.CreateToken(user.Id, 1)
	res := fmt.Sprintf("Name: %s\nToken: %s", result.Name, userToken)
	return c.String(http.StatusOK, res)
}
