package goffee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// trace 获取触发 panic 的堆栈信息
func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // 跳过前三个Caller

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

// Recovery 使用 defer 挂载上错误恢复的函数，
// 在这个函数中调用 *recover()*，捕获 panic，并且将堆栈信息打印在日志中，向用户返回 Internal Server Error
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		c.Next()
	}
}
