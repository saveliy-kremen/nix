package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Post struct {
	UserID int `json:"userId"`
	Id     int
	Title  string
	Body   string
}

type Comment struct {
	PostID int `json:"postId"`
	Id     int
	Name   string
	Email  string
	Body   string
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

	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts?userId=7")
	if err != nil {
		fmt.Println(err)
		panic("http error")
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	posts := []Post{}
	err = json.Unmarshal(body, &posts)
	if err != nil {
		panic(err.Error())
	}

	var wgPosts sync.WaitGroup
	for _, post := range posts {
		wgPosts.Add(1)
		go savePost(post, &wgPosts)
	}
	wgPosts.Wait()
}

func savePost(post Post, wg *sync.WaitGroup) {
	defer wg.Done()
	result := db.Create(&post)
	if result.Error != nil {
		panic(result.Error.Error())
	}
	resp, err := http.Get("https://jsonplaceholder.typicode.com/comments?postId=" + strconv.Itoa(post.Id))
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	comments := []Comment{}
	err = json.Unmarshal(body, &comments)
	if err != nil {
		panic(err.Error())
	}
	var wgComments sync.WaitGroup
	for _, comment := range comments {
		wgComments.Add(1)
		go saveComment(comment, &wgComments)
	}
	wgComments.Wait()
}

func saveComment(comment Comment, wg *sync.WaitGroup) {
	defer wg.Done()
	result := db.Create(&comment)
	if result.Error != nil {
		panic(result.Error.Error())
	}
}
