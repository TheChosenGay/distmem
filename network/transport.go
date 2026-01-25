package network

type OnKeyFunc func(key string) (any, error)
type OnSetKeyFunc func(key string, value any) error

type Transport interface {
	Listen(addr string) error
	OnKey(keyFunc OnKeyFunc)
	OnSetKey(setKeyFunc OnSetKeyFunc)
	GetPeer(addr string) (Peer, error)
	Close() error
}
