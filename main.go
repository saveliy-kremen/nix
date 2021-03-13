package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"

	_ "github.com/go-sql-driver/mysql"
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

func main() {
	db, err := sql.Open("mysql", "Stas_nixuser:edUfw5nxpT@tcp(192.168.1.1:3306)/Stas_nix")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

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
		go savePost(post, db, &wgPosts)
	}
	wgPosts.Wait()
}

func savePost(post Post, db *sql.DB, wg *sync.WaitGroup) {
	defer wg.Done()
	_, err := db.Exec("INSERT INTO posts VALUES ( ?, ?, ?, ? )", post.UserID, post.Id, post.Title, post.Body)
	if err != nil {
		panic(err.Error())
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
		go saveComment(comment, db, &wgComments)
	}
	wgComments.Wait()
}

func saveComment(comment Comment, db *sql.DB, wg *sync.WaitGroup) {
	defer wg.Done()
	_, err := db.Exec("INSERT INTO comments VALUES ( ?, ?, ?, ?, ? )", comment.PostID, comment.Id, comment.Name, comment.Email, comment.Body)
	if err != nil {
		panic(err.Error())
	}
}
