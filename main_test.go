package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNCIsInVzZXJfcm9sZSI6IjEiLCJleHAiOjE2MTk3MTY3MDV9._9HYljbrrN6LUikG6DycmYcsMpCrDlCg3e4XHqkc7qw"
var userID = 4

var postID int
var commentID int

func Test_getPosts(t *testing.T) {
	resp, err := http.Get("http://localhost:3000/posts/")
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
		"body":  {"testBody"},
	}

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPost, "http://localhost:3000/posts/", strings.NewReader(formData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
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
		"body":  {"testUpdateBody"},
	}

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPut, "http://localhost:3000/posts/"+strconv.Itoa(postID), strings.NewReader(formData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Bearer "+token)
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
		Id:     postID,
		UserID: userID,
		Title:  "testUpdateTitle",
		Body:   "testUpdateBody",
	}
	if !cmp.Equal(post, expectedPost) {
		t.Errorf("Handler returned wrong post")
	}

}

func Test_getPost(t *testing.T) {
	resp, err := http.Get("http://localhost:3000/posts/" + strconv.Itoa(postID))
	if err != nil {
		t.Errorf("Handler returned %v", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Handler returned %v", resp.StatusCode)
	}

	var post Post
	json.NewDecoder(resp.Body).Decode(&post)
	expectedPost := Post{
		Id:     postID,
		UserID: userID,
		Title:  "testUpdateTitle",
		Body:   "testUpdateBody",
	}
	if !cmp.Equal(post, expectedPost) {
		t.Errorf("Handler returned wrong post")
	}
}

func Test_postComment(t *testing.T) {

	formData := url.Values{
		"name":  {"nameTitle"},
		"email": {"testEmail"},
		"body":  {"testBody"},
	}

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPost, "http://localhost:3000/comments/"+strconv.Itoa(postID), strings.NewReader(formData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Handler returned %v", err.Error())
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Handler returned %v", resp.StatusCode)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	commentID = int(result["ID"].(float64))
	if result["Message"].(string) != "comment created" {
		t.Errorf("Error create comment")
	}
}

func Test_putComment(t *testing.T) {
	formData := url.Values{
		"name":  {"testUpdateName"},
		"email": {"testUpdateEmail"},
		"body":  {"testUpdateBody"},
	}

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPut, "http://localhost:3000/comments/"+strconv.Itoa(commentID), strings.NewReader(formData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Handler returned %v", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Handler returned %v", resp.StatusCode)
	}

	var comment Comment
	json.NewDecoder(resp.Body).Decode(&comment)
	expectedComment := Comment{
		Id:     commentID,
		UserID: userID,
		PostID: postID,
		Name:   "testUpdateName",
		Email:  "testUpdateEmail",
		Body:   "testUpdateBody",
	}
	if !cmp.Equal(comment, expectedComment) {
		t.Errorf("Handler returned wrong comment")
	}

}

func Test_getComments(t *testing.T) {
	resp, err := http.Get("http://localhost:3000/comments/" + strconv.Itoa(postID))
	if err != nil {
		t.Errorf("Handler returned %v", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Handler returned %v", resp.StatusCode)
	}

	var comments []Comment
	json.NewDecoder(resp.Body).Decode(&comments)
	if len(comments) == 0 {
		t.Errorf("Handler returned empty comments")
	}
}

func Test_deleteComment(t *testing.T) {
	client := &http.Client{}

	req, err := http.NewRequest("DELETE", "http://localhost:3000/comments/"+strconv.Itoa(commentID), nil)
	req.Header.Add("Authorization", "Bearer "+token)
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

func Test_deletePost(t *testing.T) {
	client := &http.Client{}

	req, err := http.NewRequest("DELETE", "http://localhost:3000/posts/"+strconv.Itoa(postID), nil)
	req.Header.Add("Authorization", "Bearer "+token)
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
