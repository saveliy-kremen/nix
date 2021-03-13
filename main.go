package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func main() {
	ch := make(chan string)
	var res string
	for i := 0; i < 100; i++ {
		go MakeRequest("https://jsonplaceholder.typicode.com/posts/"+strconv.Itoa(i), ch)
	}
	for i := 0; i < 100; i++ {
		res += <-ch
	}
	fmt.Print(res)
}

func MakeRequest(url string, ch chan string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		panic("http error")
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	ch <- fmt.Sprint(string(body))
}
