package goffee

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	handlers map[string]HandlerFunc
}

// newRouter 构造一个新路由
func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc)}
}

// addRouter 添加路由到路由映射表里面
func (r *router) addRouter(method string, pattern string, handler HandlerFunc) {
	// 使用strings.builder拼接字符串，提升速度
	var key strings.Builder
	key.WriteString(method)
	key.WriteString("-")
	key.WriteString(pattern)
	// 成功添加后将其打在日志中
	defer log.Printf("Route %4s - %s", method, pattern)

	r.handlers[key.String()] = handler
}

func (r *router) handle(ctx *Context) {
	// 使用strings.builder拼接字符串，提升速度
	var key strings.Builder
	key.WriteString(ctx.Method)
	key.WriteString("-")
	key.WriteString(ctx.Path)
	// 去到路由映射表中检查是否有此路由，若有则进行处理，若无则报错404
	if handler, ok := r.handlers[key.String()]; ok {
		handler(ctx)
	} else {
		ctx.String(http.StatusNotFound, "404 Not Found: %s\n", ctx.Path)
	}
}
