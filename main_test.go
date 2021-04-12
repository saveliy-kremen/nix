package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/go-cmp/cmp"
)

var postID int
var commentID int

func Test_getPosts(t *testing.T) {
	resp, err := http.Get("http://localhost/posts/")
	if err != nil {
		t.Errorf("Handler returned %v", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Handler returned %v", resp.StatusCode)
	}

	var posts []Post
	json.NewDecoder(resp.Body).Decode(&posts)
	if len(posts) == 0 {
		t.Errorf("Handler returned empty posts")
	}
}

func Test_postPost(t *testing.T) {

	formData := url.Values{
		"title": {"testTitle"},
		"body": {"testBody"},
	}

	resp, err := http.PostForm("http://localhost/posts/", formData)
	if err != nil {
		t.Errorf("Handler returned %v", err.Error())
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Handler returned %v", resp.StatusCode)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	postID = int(result["ID"].(float64))
	if result["Message"].(string) != "post created" {
		t.Errorf("Error create post")
	}
}

func Test_putPost(t *testing.T) {
	formData := url.Values{
		"title": {"testUpdateTitle"},
		"body": {"testUpdateBody"},
	}

	client := &http.Client{}
	spew.Dump(formData.Encode())
	req, err := http.NewRequest(http.MethodPut, "http://localhost/posts/"+ strconv.Itoa(postID), strings.NewReader(formData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Handler returned %v", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Handler returned %v", resp.StatusCode)
	}

	var post Post
	json.NewDecoder(resp.Body).Decode(&post)
	expectedPost := Post{
		Id: postID,
		UserID: 7,
		Title: "testUpdateTitle",
		Body: "testUpdateBody",
	}
	if !cmp.Equal(post, expectedPost) {
		t.Errorf("Handler returned wrong post")
	}

}

func Test_getPost(t *testing.T) {
	resp, err := http.Get("http://localhost/posts/"+ strconv.Itoa(postID))
	if err != nil {
		t.Errorf("Handler returned %v", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Handler returned %v", resp.StatusCode)
	}

	var post Post
	json.NewDecoder(resp.Body).Decode(&post)
	expectedPost := Post{
		Id: postID,
		UserID: 7,
		Title: "testUpdateTitle",
		Body: "testUpdateBody",
	}
	if !cmp.Equal(post, expectedPost) {
		t.Errorf("Handler returned wrong post")
	}
}

func Test_deletePost(t *testing.T) {
	client := &http.Client{}

	req, err := http.NewRequest("DELETE", "https://httpbin.org/delete/"+ strconv.Itoa(postID), nil)
	if err != nil {
		t.Errorf("Handler returned %v", err.Error())
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Handler returned %v", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Handler returned %v", resp.StatusCode)
	}
}
