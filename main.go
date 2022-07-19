package main

import (
	"fmt"
	"goffche"
	"log"
	"net/http"
)

//
//import (
//	"log"
//	"net/http"
//
//	"goffee"
//)
//
//func main() {
//	r := goffee.Default()
//
//	hello := r.Group("/hello")
//	{
//		hello.GET("/test", helloTest)
//		hello.GET("/:name", helloById)
//	}
//
//	r.POST("/login", login)
//
//	r.GET("/assets/*filepath", func(ctx *goffee.Context) {
//		ctx.JSON(http.StatusOK, goffee.H{"filepath": ctx.Param("filepath")})
//	})
//
//	r.GET("/panic", func(c *goffee.Context) {
//		names := []string{"goffee"}
//		c.String(http.StatusOK, names[100])
//	})
//
//	err := r.Run(":9999")
//	if err != nil {
//		log.Println("run engine error, err:", err)
//		return
//	}
//}
//
//func helloTest(ctx *goffee.Context) {
//	ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Query("name"), ctx.Path)
//}
//
//func helloById(ctx *goffee.Context) {
//	ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Param("name"), ctx.Path)
//}
//
//func login(ctx *goffee.Context) {
//	ctx.JSON(http.StatusOK, goffee.H{
//		"username": ctx.PostForm("username"),
//		"password": ctx.PostForm("password"),
//	})
//}

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	goffche.NewGroup("scores", 2<<10, goffche.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := "localhost:9999"
	peers := goffche.NewHTTPPool(addr)
	log.Println("goffche is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
