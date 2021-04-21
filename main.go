package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"

	fb "./fb"
	google "./google"
	tw "./tw"

	"github.com/labstack/echo"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "./docs"
)

type Post struct {
	UserID int `json:"userId"`
	Id     int
	Title  string `form:"title"`
	Body   string `form:"body"`
}

type Comment struct {
	PostID int `json:"postId" form:"postId"`
	Id     int
	Name   string `form:"name"`
	Email  string `form:"email"`
	Body   string `form:"body"`
}

type Response struct {
	ID      int
	Message string
}

type ErrorResponse struct {
	Status  int
	Message string
}

var db *gorm.DB

// @title Echo Swagger API
// @version 1.0
// @description This is a echo post + comments server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:80
// @BasePath /
// @schemes http
func main() {
	var err error
	dsn := "Stas_nixuser:edUfw5nxpT@tcp(192.168.1.1:3306)/Stas_nix?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		panic("gorm error")
	}

	// _, err = db.Exec("CREATE TABLE posts ( user_id integer, id integer, title text, body text, PRIMARY KEY (id))")
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = db.Exec("CREATE TABLE comments ( post_id integer, id integer, name text, email varchar(256), body text, PRIMARY KEY (id))")
	// if err != nil {
	// 	panic(err)
	// }

	e := echo.New()

	// Auth
	e.GET("/", handleAuth)
	e.GET("/fb_login", fb.HandleFacebookLogin)
	e.GET("/fb_oauth2callback", fb.HandleFacebookCallback)
	e.GET("/google_login", google.HandleGoogleLogin)
	e.GET("/google_oauth2callback", google.HandleGoogleCallback)
	e.GET("/tw_login", tw.HandleTwitterLogin)
	e.GET("/tw_oauth2callback", tw.HandleTwitterCallback)

	e.GET("/posts/", getPosts)
	e.GET("/posts/:id", getPost)
	e.POST("/posts/", savePost)
	e.PUT("/posts/:id", updatePost)
	e.DELETE("/posts/:id", deletePost)
	e.GET("/comments/:id", getComments)
	e.POST("/comments/:id", saveComment)
	e.PUT("/comments/:id", updateComment)
	e.DELETE("/comments/:id", deleteComment)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	//e.Pre(middleware.HTTPSRedirect())
	e.Logger.Fatal(e.StartTLS(":3000", "cert/localhost.crt", "cert/localhost.key"))
	//https://golangexample.com/automatically-generate-restful-api-documentation-with-swagger-2-0-for-go/
}

func handleAuth(c echo.Context) error {
	var htmlIndex = `
	<html>
	  <body>
		 <a href="/fb_login">Facebook Log In</a><br>
		 <a href="/google_login">Google Log In</a><br>
		 <a href="/tw_login">Twitter Log In</a>
	  </body>
	</html>`
	return c.HTML(http.StatusOK, htmlIndex)
}

// getPosts godoc
// @Summary Get all posts.
// @Description Get posts.
// @Tags Posts
// @Param Accept-xml header string false "Header for xml output"
// @Produce json
// @Produce xml
// @Success 200 {array} []Post
// @Failure 500 {object} ErrorResponse
// @Router /posts [get]
func getPosts(c echo.Context) error {
	req := c.Request()
	headers := req.Header
	acceptXML := headers.Get("Accept-xml")
	resp := c.Response()
	var posts []Post
	db.Order("id desc").Find(&posts)
	if acceptXML == "" {
		resp.WriteHeader(http.StatusOK)
		resp.Header().Set("Content-Type", "application/json")
		json.NewEncoder(resp).Encode(posts)
	} else {
		xmlOut, err := xml.MarshalIndent(posts, "", "  ")
		if err != nil {
			echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			return nil
		}
		resp.Header().Set("Content-Type", "application/xml")
		c.String(http.StatusOK, string(xmlOut))
	}
	return nil
}

// getPost godoc
// @Summary Get post by id.
// @Description Get post based on given ID.
// @Tags Posts
// @Param id path integer true "Post ID"
// @Param Accept-xml header string false "Header for xml output"
// @Produce json
// @Produce xml
// @Success 200 {object} Post
// @Failure 500 {object} ErrorResponse
// @Router /posts/{id} [get]
func getPost(c echo.Context) error {
	req := c.Request()
	headers := req.Header
	acceptXML := headers.Get("Accept-xml")
	id := c.Param("id")
	resp := c.Response()
	post := Post{}
	result := db.First(&post, id)
	if result.Error != nil {
		echo.NewHTTPError(http.StatusBadRequest, result.Error.Error())
		return nil
	}
	if acceptXML == "" {
		resp.WriteHeader(http.StatusOK)
		resp.Header().Set("Content-Type", "application/json")
		json.NewEncoder(resp).Encode(post)
	} else {
		xmlOut, err := xml.MarshalIndent(post, "", "  ")
		if err != nil {
			echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			return nil
		}
		resp.WriteHeader(http.StatusOK)
		c.Response().Header().Set("Content-Type", "application/xml")
		resp.Write(xmlOut)
	}
	return nil
}

// savePost godoc
// @Summary Save post.
// @Description Save post.
// @Tags Posts
// @Param title formData string true "Post title"
// @Param body formData string true "Post body"
// @Produce json
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponse
// @Router /posts/ [post]
func savePost(c echo.Context) error {
	resp := c.Response()
	var post Post
	if err := c.Bind(&post); err != nil {
		return err
	}
	post.UserID = 7
	result := db.Create(&post)
	if result.Error == nil {
		resp.WriteHeader(http.StatusCreated)
		resp.Header().Set("Content-Type", "application/json")
		b, _ := json.Marshal(Response{ID: post.Id, Message: "post created"})
		resp.Write(b)
	} else {
		echo.NewHTTPError(http.StatusBadRequest, result.Error.Error())
		return nil
	}
	return nil
}

// updatePost godoc
// @Summary Update post.
// @Description Update post.
// @Tags Posts
// @Param id path integer true "Post ID"
// @Param title formData string true "Post title"
// @Param body formData string true "Post body"
// @Produce json
// @Success 200 {object} Post
// @Failure 400 {object} ErrorResponse
// @Router /posts/{id} [put]
func updatePost(c echo.Context) error {
	req := c.Request()
	headers := req.Header
	acceptXML := headers.Get("Accept-xml")
	resp := c.Response()
	post := Post{}
	id := c.Param("id")
	if id == "" {
		echo.NewHTTPError(http.StatusBadRequest, `{"message": "post not found"}`)
		return nil
	}
	result := db.First(&post, id)
	if result.Error != nil {
		echo.NewHTTPError(http.StatusBadRequest, `{"message": "post not found"}`)
		return nil
	}
	if err := c.Bind(&post); err != nil {
		return err
	}
	result = db.Save(&post)
	if result.Error != nil {
		echo.NewHTTPError(http.StatusBadRequest, `{"message": "post not found"}`)
		return nil
	}
	if acceptXML == "" {
		resp.WriteHeader(http.StatusOK)
		resp.Header().Set("Content-Type", "application/json")
		json.NewEncoder(resp).Encode(post)
	} else {
		xmlOut, err := xml.MarshalIndent(post, "", "  ")
		if err != nil {
			echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			return nil
		}
		resp.WriteHeader(http.StatusOK)
		c.Response().Header().Set("Content-Type", "application/xml")
		resp.Write(xmlOut)
	}
	return nil
}

// deletePost godoc
// @Summary Delete post.
// @Description Delete post.
// @Tags Posts
// @Param id path integer true "Post ID"
// @Produce json
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponse
// @Router /posts/{id} [delete]
func deletePost(c echo.Context) error {
	resp := c.Response()
	var post Post
	id := c.Param("id")
	if id == "" {
		echo.NewHTTPError(http.StatusBadRequest, `{"message": "post not found"}`)
		return nil
	}
	result := db.First(&post, id)
	if result.Error != nil {
		echo.NewHTTPError(http.StatusBadRequest, result.Error.Error())
		return nil
	} else {
		result = db.Delete(&post)
		if result.Error != nil {
			echo.NewHTTPError(http.StatusBadRequest, `{"message": "error post delete"}`)
			return nil
		}
		resp.WriteHeader(http.StatusOK)
		resp.Header().Set("Content-Type", "application/json")
		b, _ := json.Marshal(Response{ID: post.Id, Message: "post deleted"})
		resp.Write(b)
	}
	return nil
}

// getComments godoc
// @Summary Get comments based on given Post ID.
// @Description Get comments based on given Post ID.
// @Tags Comments
// @Param id path integer true "Post ID"
// @Param Accept-xml header string false "Header for xml output"
// @Produce json
// @Produce xml
// @Success 200 {array} []Comment
// @Failure 500 {object} ErrorResponse
// @Router /comments/{id} [get]
func getComments(c echo.Context) error {
	req := c.Request()
	headers := req.Header
	acceptXML := headers.Get("Accept-xml")
	postID := c.Param("id")
	resp := c.Response()
	var comments []Comment
	db.Where("post_id = ?", postID).Order("id desc").Find(&comments)
	if acceptXML == "" {
		resp.WriteHeader(http.StatusOK)
		resp.Header().Set("Content-Type", "application/json")
		json.NewEncoder(resp).Encode(comments)
	} else {
		xmlOut, err := xml.MarshalIndent(comments, "", "  ")
		if err != nil {
			echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			return nil
		}
		resp.WriteHeader(http.StatusOK)
		c.Response().Header().Set("Content-Type", "application/xml")
		resp.Write(xmlOut)
	}
	return nil
}

// saveComment godoc
// @Summary Save comment.
// @Description Save comment.
// @Tags Comments
// @Param id path number true "Post ID"
// @Param name formData string true "Comment name"
// @Param email formData string true "Comment email"
// @Param body formData string true "Comment body"
// @Produce json
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponse
// @Router /comments/{id} [post]
func saveComment(c echo.Context) error {
	postID := c.Param("id")
	resp := c.Response()
	var comment Comment
	if err := c.Bind(&comment); err != nil {
		return err
	}
	comment.PostID, _ = strconv.Atoi(postID)
	result := db.Create(&comment)
	if result.Error == nil {
		resp.WriteHeader(http.StatusCreated)
		resp.Header().Set("Content-Type", "application/json")
		b, _ := json.Marshal(Response{ID: comment.Id, Message: "comment created"})
		resp.Write(b)
	} else {
		echo.NewHTTPError(http.StatusBadRequest, result.Error.Error())
		return nil
	}
	return nil
}

// updateComment godoc
// @Summary Update comment.
// @Description Update comment.
// @Tags Comments
// @Param id path integer true "Comment ID"
// @Param name formData string true "Comment name"
// @Param email formData string true "Comment email"
// @Param body formData string true "Comment body"
// @Produce json
// @Produce xml
// @Success 200 {object} Comment
// @Failure 400 {object} ErrorResponse
// @Router /comments/{id} [put]
func updateComment(c echo.Context) error {
	req := c.Request()
	headers := req.Header
	acceptXML := headers.Get("Accept-xml")
	resp := c.Response()
	comment := Comment{}
	id := c.Param("id")
	if id == "" {
		echo.NewHTTPError(http.StatusBadRequest, `{"message": "comment not found"}`)
		return nil
	}
	result := db.First(&comment, id)
	if result.Error != nil {
		echo.NewHTTPError(http.StatusBadRequest, `{"message": "comment not found"}`)
		return nil
	}
	if err := c.Bind(&comment); err != nil {
		return err
	}
	result = db.Save(&comment)
	if result.Error != nil {
		echo.NewHTTPError(http.StatusBadRequest, result.Error.Error())
		return nil
	}
	if acceptXML == "" {
		resp.WriteHeader(http.StatusOK)
		resp.Header().Set("Content-Type", "application/json")
		json.NewEncoder(resp).Encode(comment)
	} else {
		xmlOut, err := xml.MarshalIndent(comment, "", "  ")
		if err != nil {
			echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			return nil
		}
		resp.WriteHeader(http.StatusOK)
		c.Response().Header().Set("Content-Type", "application/xml")
		resp.Write(xmlOut)
	}
	return nil
}

// deleteComment godoc
// @Summary Delete comment.
// @Description Delete comment.
// @Tags Comments
// @Param id path integer true "Comment ID"
// @Produce json
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponse
// @Router /comments/{id} [delete]
func deleteComment(c echo.Context) error {
	resp := c.Response()
	var comment Comment
	id := c.Param("id")
	if id == "" {
		echo.NewHTTPError(http.StatusBadRequest, `{"message": "comment not found"}`)
		return nil
	}
	result := db.First(&comment, id)
	if result.Error != nil {
		echo.NewHTTPError(http.StatusBadRequest, result.Error.Error())
		return nil
	} else {
		result = db.Delete(&comment)
		if result.Error != nil {
			echo.NewHTTPError(http.StatusBadRequest, `{"message": "error comment delete"}`)
			return nil
		}
		resp.WriteHeader(http.StatusOK)
		resp.Header().Set("Content-Type", "application/json")
		b, _ := json.Marshal(Response{ID: comment.Id, Message: "comment deleted"})
		resp.Write(b)
	}
	return nil
}
