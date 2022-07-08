package main

import (
	"log"
	"net/http"

	"goffee"
)

func main() {
	r := goffee.New()

	r.GET("/hello", hello)
	r.POST("/login", login)

	err := r.Run(":9999")
	if err != nil {
		log.Println("run engine error, err:", err)
		return
	}
}

func hello(ctx *goffee.Context) {
	ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Query("name"), ctx.Path)
}

func login(ctx *goffee.Context) {
	ctx.JSON(http.StatusOK, goffee.H{
		"username": ctx.PostForm("username"),
		"password": ctx.PostForm("password"),
	})
}
