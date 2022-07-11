package goffee

import (
	"log"
	"time"
)

// Logger 日志中间件，
func Logger() HandlerFunc {
	return func(ctx *Context) {
		// 开始计时
		t := time.Now()
		// 进行其他操作
		ctx.Next()
		// 最后进行日志打印
		log.Printf("[%d] %s in %v", ctx.StatusCode, ctx.Req.RequestURI, time.Since(t))
	}
}
