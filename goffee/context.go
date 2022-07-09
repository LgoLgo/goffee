package goffee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// H 定义一个接口map，构建JSON数据时，显得更简洁
type H map[string]interface{}

// Context 上下文结构体
type Context struct {
	// 请求响应原始参数
	Writer http.ResponseWriter
	Req    *http.Request
	// 请求路径、方法
	Path   string
	Method string
	Params map[string]string
	// 响应状态码
	StatusCode int
}

// newContext 构造一个新的上下文
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

// Param 实现动态路由
func (ctx *Context) Param(key string) string {
	value, _ := ctx.Params[key]
	return value
}

// PostForm 访问PostForm参数
func (ctx *Context) PostForm(key string) string {
	return ctx.Req.FormValue(key)
}

// Query 访问Query参数
func (ctx *Context) Query(key string) string {
	return ctx.Req.URL.Query().Get(key)
}

// Status 设置状态
func (ctx *Context) Status(code int) {
	ctx.StatusCode = code
	ctx.Writer.WriteHeader(code)
}

// SetHeader 设置请求头
func (ctx *Context) SetHeader(key string, value string) {
	ctx.Writer.Header().Set(key, value)
}

// String 构造String响应
func (ctx *Context) String(code int, format string, values ...interface{}) {
	ctx.SetHeader("Content-Type", "text/plain")
	ctx.Status(code)
	ctx.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// JSON 构造JSON响应
func (ctx *Context) JSON(code int, obj interface{}) {
	ctx.SetHeader("Content-Type", "application/json")
	ctx.Status(code)
	encoder := json.NewEncoder(ctx.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(ctx.Writer, err.Error(), 500)
	}
}

// Data 构造Data响应（通过byte数组实现）
func (ctx *Context) Data(code int, data []byte) {
	ctx.Status(code)
	ctx.Writer.Write(data)
}

// HTML 构造HTML响应
func (ctx *Context) HTML(code int, html string) {
	ctx.SetHeader("Content-Type", "text/html")
	ctx.Status(code)
	ctx.Writer.Write([]byte(html))
}
