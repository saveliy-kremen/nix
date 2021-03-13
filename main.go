package main

import (
    "fmt"
    "net/http"
)

func main() {
	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts") 
	if err != nil { 
		fmt.Println(err) 
		return 
	} 
	defer resp.Body.Close()
	for true {
		bs := make([]byte, 1024)
		n, err := resp.Body.Read(bs)
		fmt.Print(string(bs[:n]))
		if n == 0 || err != nil{
			break
		}
	}
}
