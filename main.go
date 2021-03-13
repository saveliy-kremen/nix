package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func main() {
	directory := "./storage/posts/"
	if _, err := os.Stat(directory); err != nil {
		os.MkdirAll(directory, 0775)
	}

	ch := make(chan []byte)
	for i := 0; i < 100; i++ {
		go MakeRequest("https://jsonplaceholder.typicode.com/posts/"+strconv.Itoa(i), ch)
	}
	for i := 0; i < 100; i++ {
		ioutil.WriteFile("./storage/posts/"+strconv.Itoa(i+1)+".txt", <-ch, 0644)
	}
}

func MakeRequest(url string, ch chan []byte) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		panic("http error")
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	ch <- body
}
