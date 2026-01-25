package hash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type HashFunc func(string) uint32

type HashCircleOpts struct {
	hashFunc HashFunc
	replicas int
}

type HashCircle struct {
	keys   []uint32
	keyMap map[uint32]string
	opts   HashCircleOpts
}

func NewHashCircle(opts HashCircleOpts) *HashCircle {
	if opts.hashFunc == nil {
		opts.hashFunc = func(key string) uint32 {
			return crc32.ChecksumIEEE([]byte(key))
		}
	}
	return &HashCircle{
		keys:   make([]uint32, 0),
		keyMap: make(map[uint32]string),
		opts:   opts,
	}
}

func (h *HashCircle) Add(keys ...string) {
	hashKeys := make([]uint32, len(keys)*h.opts.replicas)
	for i, k := range keys {
		for j, hk := range h.calculateKey(k) {
			hashKeys[i*h.opts.replicas+j] = hk
		}
	}

	h.keys = append(h.keys, hashKeys...)
	sort.Slice(h.keys, func(i, j int) bool {
		return h.keys[i] < h.keys[j]
	})
}

func (h *HashCircle) Get(key string) string {
	hKey := h.opts.hashFunc(key)
	idx := sort.Search(len(h.keys), func(i int) bool {
		return h.keys[i] >= hKey
	})
	if idx == len(h.keys) {
		idx = 0
	}
	return h.keyMap[h.keys[idx]]
}

func (h *HashCircle) calculateKey(key string) []uint32 {
	hashKeys := make([]uint32, h.opts.replicas)
	for i := range h.opts.replicas {
		hashKeys[i] = h.opts.hashFunc(key + strconv.Itoa(i))
		h.keyMap[hashKeys[i]] = key
	}
	return hashKeys
}
