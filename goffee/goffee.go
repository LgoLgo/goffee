package goffee

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// HandlerFunc 将网络气请求处理进行定义，给goffee使用
type HandlerFunc func(w http.ResponseWriter, r *http.Request)

// Engine 用于实现ServeHTTP的接口
type Engine struct {
	// 添加路由映射表
	router map[string]HandlerFunc
}

// New 调用来构造一个goffee引擎
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

// addRouter 添加路由到路由映射表里面
func (e *Engine) addRouter(method string, pattern string, handler HandlerFunc) {
	// 使用strings.builder拼接字符串，提升速度
	var key strings.Builder
	key.WriteString(method)
	key.WriteString("-")
	key.WriteString(pattern)
	// 成功添加后将其打在日志中
	defer log.Printf("Route %4s - %s", method, pattern)

	e.router[key.String()] = handler
}

// GET 新增GET请求
func (e *Engine) GET(pattern string, handler HandlerFunc) {
	e.addRouter("GET", pattern, handler)
}

// POST 新增POST请求
func (e *Engine) POST(pattern string, handler HandlerFunc) {
	e.addRouter("POST", pattern, handler)
}

// PUT 新增PUT请求
func (e *Engine) PUT(pattern string, handler HandlerFunc) {
	e.addRouter("PUT", pattern, handler)
}

// DELETE 新增DELETE请求
func (e *Engine) DELETE(pattern string, handler HandlerFunc) {
	e.addRouter("DELETE", pattern, handler)
}

// ServeHTTP 实现方法ServeHTTP
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 使用strings.builder拼接字符串，提升速度
	var key strings.Builder
	key.WriteString(req.Method)
	key.WriteString("-")
	key.WriteString(req.URL.Path)
	// 去到路由映射表中检查是否有此路由，若有则进行处理，若无则报错404
	if handler, ok := e.router[key.String()]; ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 Not Found: %s\n", req.URL)
	}
}

// Run 定义了启动http服务器的方法
func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}
