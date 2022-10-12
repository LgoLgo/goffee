package Lgoffche

// PeerPicker 为必须实现定位的接口
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter 是 peer 必须实现的接口。
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}
