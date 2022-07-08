package goffee

import (
	"log"
	"net/http"
)

// HandlerFunc 将网络气请求处理进行定义，给goffee使用
type HandlerFunc func(ctx *Context)

// Engine 用于实现ServeHTTP的接口
type Engine struct {
	// 添加路由映射表
	router *router
}

// New 调用来构造一个goffee引擎
func New() *Engine {
	return &Engine{router: newRouter()}
}

// addRouter 添加路由到路由映射表里面
func (e *Engine) addRouter(method string, pattern string, handler HandlerFunc) {
	e.router.addRouter(method, pattern, handler)
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
	ctx := newContext(w, req)
	e.router.handle(ctx)
}

// Run 定义了启动http服务器的方法
func (e *Engine) Run(addr string) error {
	log.Printf("The project is listening on port %s\n", addr)
	return http.ListenAndServe(addr, e)
}
