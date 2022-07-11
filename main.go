package main

import (
	"log"
	"net/http"

	"goffee"
)

func main() {
	r := goffee.New()
	r.Use(goffee.Logger())

	hello := r.Group("/hello")
	{
		hello.GET("/test", helloTest)
		hello.GET("/:name", helloById)
	}

	r.POST("/login", login)

	r.GET("/assets/*filepath", func(ctx *goffee.Context) {
		ctx.JSON(http.StatusOK, goffee.H{"filepath": ctx.Param("filepath")})
	})

	err := r.Run(":9999")
	if err != nil {
		log.Println("run engine error, err:", err)
		return
	}
}

func helloTest(ctx *goffee.Context) {
	ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Query("name"), ctx.Path)
}

func helloById(ctx *goffee.Context) {
	ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Param("name"), ctx.Path)
}

func login(ctx *goffee.Context) {
	ctx.JSON(http.StatusOK, goffee.H{
		"username": ctx.PostForm("username"),
		"password": ctx.PostForm("password"),
	})
}
