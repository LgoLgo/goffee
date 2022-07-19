package goffche

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/_goffche/"

// HTTPPool 作为承载节点间 HTTP 通信的核心数据结构
type HTTPPool struct {
	// e.g. "https://example.net/_goffee/"
	self     string // 用于记录自己的地址
	basePath string // 作为节点间通讯地址的前缀
}

// NewHTTPPool initializes an HTTP pool of peers.
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// Log 打印服务日志
func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// ServeHTTP 实现ServerHTTP方法
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 首先判断访问路径的前缀是否是 basePath
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("%s %s", r.Method, r.URL.Path)
	// /<basepath>/<groupname>/<key> 此为约定格式
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]
	// 获取缓存数据
	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	// 将缓存值作为 httpResponse 的 body 返回
	w.Write(view.ByteSlice())
}
