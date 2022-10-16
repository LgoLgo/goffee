package singleflight

import "sync"

// 代表正在进行中，或已经结束的请求
type call struct {
	// 避免重入
	wg  sync.WaitGroup
	val interface{}
	err error
}

// Group 代表一个工作类并形成一个命名空间，其中
// 可以使用重复抑制来执行工作单元
// 管理不同 key 的请求(call)
type Group struct {
	mu sync.Mutex
	m  map[string]*call // 延迟初始化
}

// Do 执行并返回给定函数的结果，确保只有一次执行在进行中。如果出现重复，则重复调用者等待
// 针对相同的 key，无论 Do 被调用多少次，函数 fn 都只会被调用一次，等待 fn 调用结束了，返回返回值或错误
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}
