package hash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type HashFunc func(string) uint32

type HashCircleOpts struct {
	HashFunc HashFunc
	Replicas int
}

type HashCircle struct {
	keys   []uint32
	keyMap map[uint32]string
	opts   HashCircleOpts
}

func NewHashCircle(opts HashCircleOpts) *HashCircle {
	if opts.HashFunc == nil {
		opts.HashFunc = func(key string) uint32 {
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
	hashKeys := make([]uint32, len(keys)*h.opts.Replicas)
	for i, k := range keys {
		for j, hk := range h.calculateKey(k) {
			hashKeys[i*h.opts.Replicas+j] = hk
		}
	}

	h.keys = append(h.keys, hashKeys...)
	sort.Slice(h.keys, func(i, j int) bool {
		return h.keys[i] < h.keys[j]
	})
}

func (h *HashCircle) Delete(keys ...string) {
	for _, k := range keys {
		for _, hk := range h.calculateKey(k) {
			idx := sort.Search(len(h.keys), func(i int) bool {
				return h.keys[i] == hk
			})
			if idx < len(h.keys) {
				h.keys = append(h.keys[:idx], h.keys[idx+1:]...)
				delete(h.keyMap, hk)
			}
		}
	}
}

// Get the peer address for the given key
func (h *HashCircle) Get(key string) string {
	hKey := h.opts.HashFunc(key)
	idx := sort.Search(len(h.keys), func(i int) bool {
		return h.keys[i] >= hKey
	})
	if idx == len(h.keys) {
		idx = 0
	}
	return h.keyMap[h.keys[idx]]
}

func (h *HashCircle) calculateKey(key string) []uint32 {
	hashKeys := make([]uint32, h.opts.Replicas)
	for i := range h.opts.Replicas {
		hashKeys[i] = h.opts.HashFunc(key + strconv.Itoa(i))
		h.keyMap[hashKeys[i]] = key
	}
	return hashKeys
}
