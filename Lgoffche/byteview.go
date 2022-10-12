package Lgoffche

// ByteView 只读数据结构，用来表示缓存值
type ByteView struct {
	// b 将会存储真实的缓存值，使用 byte 类型用以支持任意的数据类型存储
	b []byte
}

// Len 返回长度
func (v ByteView) Len() int {
	return len(v.b)
}

// ByteSlice 返回一个拷贝，防止缓存值被外部程序修改
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

// String 将数据用字符串返回，用以制作必要的副本
func (v ByteView) String() string {
	return string(v.b)
}

// cloneBytes 拷贝
func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
