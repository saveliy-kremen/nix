package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Post struct {
	UserID int `json:"userId"`
	Id     int
	Title  string `form:"title"`
	Body   string `form:"body"`
}

type Comment struct {
	UserID int `json:"userId"`
	PostID int `json:"postId" form:"postId"`
	Id     int
	Name   string `form:"name"`
	Email  string `form:"email"`
	Body   string `form:"body"`
}

type User struct {
	Id       int
	Name     string
	FbID     string
	GoogleID string
	TwID     string
}

func init() {
	var err error

	// _, err = db.Exec("CREATE TABLE posts ( user_id integer, id integer, title text, body text, PRIMARY KEY (id))")
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = db.Exec("CREATE TABLE comments ( user_id integer, post_id integer, id integer, name text, email varchar(256), body text, PRIMARY KEY (id))")
	// if err != nil {
	// 	panic(err)
	// }

	//db.Exec("CREATE TABLE users ( id integer, name text, fb_id varchar(256), google_id varchar(256), tw_id varchar(256), PRIMARY KEY (id))")

	dsn := "Stas_nixuser:edUfw5nxpT@tcp(192.168.1.1:3306)/Stas_nix?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		panic("gorm error")
	}
}
