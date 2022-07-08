package main

import (
	"fmt"
	"net/http"

	"goffee"
)

func main() {
	r := goffee.New()

	r.GET("/", test)
	r.GET("/hello", hello)

	r.Run(":9999")
}

func test(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
}

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello, World %q\n", req.URL)
}
