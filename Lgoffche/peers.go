package Lgoffche

import pb "github.com/LgoLgo/Lgoffee/Lgoffche/Lgoffchepb/gen"

// PeerPicker 为必须实现定位的接口
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter 是 peer 必须实现的接口。
type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}
