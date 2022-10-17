# Lgoffee

A lightweight Go framework/轻量级的Go语言开发框架

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

### Lgoffee

```go
package main

import (
	"log"
	"net/http"

	"github.com/LgoLgo/Lgoffee/Lgoffee"
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

### Lgoffche

```go
package main

/*
$ curl "http://localhost:9999/api?key=Tom"
630

$ curl "http://localhost:9999/api?key=kkk"
kkk not exist
*/

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/LgoLgo/Lgoffee/Lgoffche"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *Lgoffche.Group {
	return Lgoffche.NewGroup("scores", 2<<10, Lgoffche.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

func startCacheServer(addr string, addrs []string, Lgo *Lgoffche.Group) {
	peers := Lgoffche.NewHTTPPool(addr)
	peers.Set(addrs...)
	Lgo.RegisterPeers(peers)
	log.Println("Lgoffche is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startAPIServer(apiAddr string, Lgo *Lgoffche.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := Lgo.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())

		}))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))

}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "Lgoffche server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	Lgo := createGroup()
	if api {
		go startAPIServer(apiAddr, Lgo)
	}
	startCacheServer(addrMap[port], addrs, Lgo)
}

```