package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type Post struct {
	UserID int `json:"userId"`
	Id     int
	Title  string
	Body   string
}

func main() {
	directory := "./storage/posts/"
	if _, err := os.Stat(directory); err != nil {
		os.MkdirAll(directory, 0775)
	}

	ch := make(chan Post)
	for i := 0; i < 100; i++ {
		go MakeRequest("https://jsonplaceholder.typicode.com/posts/"+strconv.Itoa(i), ch)
	}
	for i := 0; i < 100; i++ {
		ioutil.WriteFile("./storage/posts/"+strconv.Itoa(i+1)+".txt", []byte(fmt.Sprintf("%+v\n", <-ch)), 0644)
	}
}

func MakeRequest(url string, ch chan Post) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		panic("http error")
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	post := Post{}
	err = json.Unmarshal(body, &post)
	if err != nil {
		fmt.Println(err)
		panic("json decode error")
	}
	ch <- post
}
