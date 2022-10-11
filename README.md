# goffee
A lightweight Go framework/轻量级的Go语言开发框架

## 前言
It's still being built.

还在建造中

## 安装

1. 需要已经安装 [Go](https://golang.org/)

```sh
go get -u github.com/L2ncE/goffee
```

2. 将它 import 到你的代码中

```go
import "github.com/L2ncE/goffee"
```

## 快速开始

```go
package main


import (
	"log"
	"net/http"

	"goffee"
)

func main() {
	r := goffee.Default()

	hello := r.Group("/hello")
	{
		hello.GET("/test", helloTest)
		hello.GET("/:name", helloById)
	}

	r.POST("/login", login)

	r.GET("/assets/*filepath", func(ctx *goffee.Context) {
		ctx.JSON(http.StatusOK, goffee.H{"filepath": ctx.Param("filepath")})
	})

	r.GET("/panic", func(c *goffee.Context) {
		names := []string{"goffee"}
		c.String(http.StatusOK, names[100])
	})

	err := r.Run(":9999")
	if err != nil {
		log.Println("run engine error, err:", err)
		return
	}
}
```

## 开发日志

### HTTP

2022.7.9 实现了路由映射表，提供了用户注册静态路由的方法，包装了启动服务的函数

2022.7.10 提供了对 Method 和 Path 这两个常用属性的直接访问。
提供了访问Query和PostForm参数的方法。
提供了快速构造String/Data/JSON/HTML响应的方法。
将路由相关功能单独抽离了出来

2022.7.11 通过前缀树实现了动态路由的功能，并实现了分组路由

2022.7.12 实现了中间件功能并提供了日志中间件

2022.7.13 支持模板渲染

2022.7.15 实现错误恢复功能

### Cache

2022.7.17 完成LRU缓存淘汰算法

2022.7.18 使用互斥锁完成封装，支持并发

2022.7.19 完成缓存的HTTP服务端

2022.8.5 实现了一致性哈希