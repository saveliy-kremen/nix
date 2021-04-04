package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/labstack/echo"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strconv"
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

var db *gorm.DB

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
	e.GET("/posts/", getPosts)
	e.GET("/posts/:id", getPost)
	e.POST("/posts/", savePost)
	e.PUT("/posts/:id", updatePost)
	e.DELETE("/posts/:id", deletePost)
	e.GET("/comments/:id", getComments)
	e.POST("/comments/:id", saveComment)
	e.PUT("/comments/:id", updateComment)
	e.DELETE("/comments/:id", deleteComment)

	e.Logger.Fatal(e.Start(":80"))
}

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
		resp.Write([]byte(`{"message": "post created"}`))
	} else {
		echo.NewHTTPError(http.StatusBadRequest, result.Error.Error())
		return nil
	}
	return nil
}

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
		resp.Write([]byte(`{"message": "post deleted"}`))
	}
	return nil
}

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
		resp.Write([]byte(`{"message": "post created"}`))
	} else {
		echo.NewHTTPError(http.StatusBadRequest, result.Error.Error())
		return nil
	}
	return nil
}

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
		resp.Write([]byte(`{"message": "comment deleted"}`))
	}
	return nil
}
