package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
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

	http.HandleFunc("/posts/", posts)
	http.HandleFunc("/comments/", comments)
	fmt.Printf("Starting server at port 80\n")
	http.ListenAndServe(":80", nil)
}

func posts(w http.ResponseWriter, r *http.Request) {
	acceptXML := r.Header.Get("Accept-xml")
	switch r.Method {
	case "GET":
		id := strings.TrimPrefix(r.URL.Path, "/posts/")
		if id == "" {
			var posts []Post
			db.Order("id desc").Find(&posts)
			if acceptXML == "" {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(posts)
			} else {
				xmlOut, err := xml.MarshalIndent(posts, "", "  ")
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/xml")
				w.Write(xmlOut)
			}
		} else {
			post := Post{}
			result := db.First(&post, id)
			if result.Error != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"message": "post not found"}`))
				return
			}
			if acceptXML == "" {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(post)
			} else {
				xmlOut, err := xml.MarshalIndent(post, "", "  ")
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/xml")
				w.Write(xmlOut)
			}
		}
	case "POST":
		r.ParseForm() // Parses the request body
		title := r.FormValue("title")
		body := r.FormValue("body")
		userId := 7
		post := Post{Title: title, Body: body, UserID: userId}
		result := db.Create(&post)
		if result.Error == nil {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"message": "post created"}`))
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "error post create"}`))
		}
	case "PUT":
		id := strings.TrimPrefix(r.URL.Path, "/posts/")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "post not found"}`))
			return
		}
		post := Post{}
		result := db.First(&post, id)
		if result.Error != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "post not found"}`))
			return
		}
		r.ParseForm() // Parses the request body
		post.Title = r.FormValue("title")
		post.Body = r.FormValue("body")
		result = db.Save(&post)
		if result.Error != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "error post update"}`))
			return
		}
		if acceptXML == "" {
			w.WriteHeader(http.StatusAccepted)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(post)
		} else {
			xmlOut, err := xml.MarshalIndent(post, "", "  ")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/xml")
			w.Write(xmlOut)
		}
	case "DELETE":
		id := strings.TrimPrefix(r.URL.Path, "/posts/")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "post not found"}`))
			return
		}
		post := Post{}
		result := db.First(&post, id)
		if result.Error != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "post not found"}`))
			return
		}
		result = db.Delete(&post)
		if result.Error != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "error post delete"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "post deleted"}`))
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "not found"}`))
	}
}

func comments(w http.ResponseWriter, r *http.Request) {
	acceptXML := r.Header.Get("Accept-xml")
	switch r.Method {
	case "GET":
		postID := strings.TrimPrefix(r.URL.Path, "/comments/")
		var comments []Comment
		db.Where("post_id = ?", postID).Order("id desc").Find(&comments)
		if acceptXML == "" {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(comments)
		} else {
			xmlOut, err := xml.MarshalIndent(comments, "", "  ")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/xml")
			w.Write(xmlOut)
		}
	case "POST":
		postID, _ := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/comments/"))
		r.ParseForm()
		name := r.FormValue("name")
		email:= r.FormValue("email")
		body := r.FormValue("body")
		comment := Comment{PostID: postID, Name: name, Email: email, Body: body}
		spew.Dump(comment)
		result := db.Create(&comment)
		if result.Error == nil {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"message": "comment created"}`))
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "error comment create"}`))
		}
	case "PUT":
		id := strings.TrimPrefix(r.URL.Path, "/comments/")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "comment not found"}`))
			return
		}
		comment := Comment{}
		result := db.First(&comment, id)
		if result.Error != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "comment not found"}`))
			return
		}
		r.ParseForm()
		comment.Name = r.FormValue("name")
		comment.Email = r.FormValue("email")
		comment.Body  = r.FormValue("body")// Parses the request body
		result = db.Save(&comment)
		if result.Error != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "error comment update"}`))
			return
		}
		if acceptXML == "" {
			w.WriteHeader(http.StatusAccepted)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(comment)
		} else {
			xmlOut, err := xml.MarshalIndent(comment, "", "  ")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/xml")
			w.Write(xmlOut)
		}
	case "DELETE":
		id := strings.TrimPrefix(r.URL.Path, "/comments/")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "comment not found"}`))
			return
		}
		comment := Comment{}
		result := db.First(&comment, id)
		if result.Error != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "comment not found"}`))
			return
		}
		result = db.Delete(&comment)
		if result.Error != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "error comment delete"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "comment deleted"}`))
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "not found"}`))
	}
}


