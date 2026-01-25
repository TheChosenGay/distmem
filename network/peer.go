package network

type CacheOpt interface {
	Get(key string) (any, error)
	Set(key string, value any) error
}

type Peer interface {
	CacheOpt
	Addr() string
}
