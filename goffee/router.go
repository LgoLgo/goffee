package goffee

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

// parsePattern 进行解析，只允许一个“*”存在
func parsePattern(pattern string) []string {
	// 生成去除“/”的字符串数组
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// newRouter 构造一个新路由
func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// addRouter 添加路由到路由映射表里面
func (r *router) addRouter(method string, pattern string, handler HandlerFunc) {
	//先进行解析
	parts := parsePattern(pattern)
	// 使用strings.builder拼接字符串，提升速度
	var key strings.Builder
	key.WriteString(method)
	key.WriteString("-")
	key.WriteString(pattern)
	// 成功添加后将其打在日志中
	defer log.Printf("Route %4s - %s", method, pattern)
	//查看是否有此方法的路由前缀树
	_, ok := r.roots[method]
	if !ok {
		//没有就新增一个前缀树
		r.roots[method] = &node{}
	}
	//进行插入
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key.String()] = handler
}

// getRoute 得到路由，动态路由实现
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]
	//检查是否有此方法的前缀树
	if !ok {
		//没有就退出
		return nil, nil
	}
	//进行查找
	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			//进行动态路由实现
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

// handle 用于实现ServeHTTP
func (r *router) handle(ctx *Context) {
	//得到路由
	n, params := r.getRoute(ctx.Method, ctx.Path)
	if n != nil {
		ctx.Params = params
		// 使用strings.builder拼接字符串，提升速度
		var key strings.Builder
		key.WriteString(ctx.Method)
		key.WriteString("-")
		key.WriteString(n.pattern)
		r.handlers[key.String()](ctx)
	} else {
		// 去到路由映射表中检查是否有此路由，若有则进行处理，若无则报错404
		ctx.String(http.StatusNotFound, "404 Not Found: %s\n", ctx.Path)
	}
}
