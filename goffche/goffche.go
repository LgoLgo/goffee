package goffche

import (
	"fmt"
	"log"
	"sync"
)

// Group 缓存命名空间与相关数据加载
type Group struct {
	name      string
	getter    Getter // 缓存未命中时获取源数据的回调
	mainCache cache  // 开始的并发缓存
	peers     PeerPicker
}

// Getter 为key加载数据
type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc 回调函数，用函数实现Getter
type GetterFunc func(key string) ([]byte, error)

// Get 实现Getter接口的Get方法
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// NewGroup 用以实例化Group
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	// 将group存储在全局变量groups中
	groups[name] = g
	return g
}

// GetGroup 返回特定名称的Group，若没有则返回空
// 使用了只读锁，因不涉及任何写操作
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

// Get 获取键对应的缓存值
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}
	// 缓存不存在，调用load方法
	return g.load(key)
}

// load 用于调用getLocally
func (g *Group) load(key string) (value ByteView, err error) {
	if g.peers != nil {
		if peer, ok := g.peers.PickPeer(key); ok {
			if value, err = g.getFromPeer(peer, key); err == nil {
				return value, nil
			}
			log.Println("[goffche] Failed to get from peer", err)
		}
	}

	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	// 调用回调函数
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err

	}

	value := ByteView{b: cloneBytes(bytes)}
	// 将源数据添加到缓存mainCache中
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}

// RegisterPeers 注册一个 PeerPicker 来选择远端 peers
func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: bytes}, nil
}
