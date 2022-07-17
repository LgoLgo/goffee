package lru

import "container/list"

// Cache 并发访问不安全的LRU缓存
type Cache struct {
	maxBytes int64                    //允许使用的最大内存
	nbytes   int64                    //当前已使用的内存
	ll       *list.List               //Go标准库实现的双向链表 list.List
	cache    map[string]*list.Element //键是字符串，值是双向链表中对应节点的指针
	// 是可选择的，并在清理时执行
	OnEvicted func(key string, value Value)
}

// entry 键值对，为双向链表节点的数据类型
// 在链表中仍保存每个值对应的 key 的好处在于，淘汰队首节点时，需要用 key 从字典中删除对应的映射
type entry struct {
	key   string
	value Value
}

// Value 使用 Len 来计算它需要多少字节
type Value interface {
	Len() int
}

// New 缓存的构造函数，用于实例化一个Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Add 新增/修改缓存
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		//如果键存在，则更新对应节点的值，并将该节点移到队尾
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		//不存在则是新增场景，首先队尾添加新节点 &entry{key, value}, 并字典中添加 key 和节点的映射关系
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}

	//更新 c.nbytes，如果超过了设定的最大值 c.maxBytes，则移除最少访问的节点
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// Get 用于查找键的值
// 第一步是从字典中找到对应的双向链表的节点，第二步，将该节点移动到队尾
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele) //即将链表中的节点 ele 移动到队尾
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest 缓存淘汰，即移除最近最少访问的节点（队首）
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back() //取到队首节点
	if ele != nil {
		c.ll.Remove(ele) //从链表中删除
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)                                //从字典中删除该节点的映射关系。
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len()) //更新当前所用内存
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value) //调用回调函数
		}
	}
}

// Len 获取添加了多少条数据
func (c *Cache) Len() int {
	return c.ll.Len()
}
