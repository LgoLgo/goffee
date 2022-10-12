# goffee
A lightweight Go framework/轻量级的Go语言开发框架

## 前言
It's still being built.

还在建造中

## 安装

1. 需要已经安装 [Go](https://golang.org/)

```sh
go get -u github.com/LgoLgo/Lgoffee
```

2. 将它 import 到你的代码中

```go
import "github.com/LgoLgo/Lgoffee"
```

## 快速开始

### goffee

```go
package main


import (
	"log"
	"net/http"

	"github.com/LgoLgo/Lgoffee"
)

func main() {
	r := Lgoffee.Default()

	hello := r.Group("/hello")
	{
		hello.GET("/test", helloTest)
		hello.GET("/:name", helloById)
	}

	r.POST("/login", login)

	r.GET("/assets/*filepath", func(ctx *Lgoffee.Context) {
		ctx.JSON(http.StatusOK, Lgoffee.H{"filepath": ctx.Param("filepath")})
	})

	r.GET("/panic", func(c *Lgoffee.Context) {
		names := []string{"Lgoffee"}
		c.String(http.StatusOK, names[100])
	})

	err := r.Run(":9999")
	if err != nil {
		log.Println("run engine error, err:", err)
		return
	}
}
```