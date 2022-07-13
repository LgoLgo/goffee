package goffee

import (
	"log"
	"net/http"
	"strings"
)

// HandlerFunc 将网络请求处理进行定义，给goffee使用
type HandlerFunc func(ctx *Context)

// Engine 用于实现ServeHTTP的接口
type (
	Engine struct {
		// Engine拥有RouterGroup所有的能力
		*RouterGroup
		// 添加路由映射表
		router *router
		//路由组数组
		groups        []*RouterGroup
		htmlTemplates *template.Template // for html render
		funcMap       template.FuncMap   // for html render
	}
	RouterGroup struct {
		prefix      string
		middlewares []HandlerFunc // 支持中间件
		parent      *RouterGroup  // 支持分组
		engine      *Engine       // 共用一个引擎
	}
)

// New 调用来构造一个goffee引擎
func New() *Engine {
	// 依次定义
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Group 实现一个新的分组，所有组共用同一个引擎
func (g *RouterGroup) Group(prefix string) *RouterGroup {
	engine := g.engine
	newGroup := &RouterGroup{
		prefix: g.prefix + prefix,
		parent: g,
		// 使用同一个引擎
		engine: engine,
	}
	// 加入到组群中
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// addRouter 添加路由到路由映射表里面
func (g *RouterGroup) addRouter(method string, comp string, handler HandlerFunc) {
	pattern := g.prefix + comp
	g.engine.router.addRouter(method, pattern, handler)
}

// GET 新增GET请求
func (g *RouterGroup) GET(pattern string, handler HandlerFunc) {
	g.addRouter("GET", pattern, handler)
}

// POST 新增POST请求
func (g *RouterGroup) POST(pattern string, handler HandlerFunc) {
	g.addRouter("POST", pattern, handler)
}

// PUT 新增PUT请求
func (g *RouterGroup) PUT(pattern string, handler HandlerFunc) {
	g.addRouter("PUT", pattern, handler)
}

// DELETE 新增DELETE请求
func (g *RouterGroup) DELETE(pattern string, handler HandlerFunc) {
	g.addRouter("DELETE", pattern, handler)
}

// Use 将中间件添加到中间件组中
func (g *RouterGroup) Use(middlewares ...HandlerFunc) {
	g.middlewares = append(g.middlewares, middlewares...)
}

// createStaticHandler
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		// Check if file exists and/or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// Static
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// Register GET handlers
	group.GET(urlPattern, handler)
}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

// ServeHTTP 实现方法ServeHTTP
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range e.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = engine
	e.router.handle(c)
}

// Run 定义了启动http服务器的方法
func (e *Engine) Run(addr string) error {
	log.Printf("The project is listening on port %s\n", addr)
	return http.ListenAndServe(addr, e)
}
