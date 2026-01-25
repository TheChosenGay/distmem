package cache

import (
	"errors"
	"sync"
)

var (
	ErrKeyNotFound = errors.New("key not found")
	ErrKeyExists   = errors.New("key already exists")
)

type Cache interface {
	Get(key string) (any, error)
	Set(key string, value any) error
	Delete(key string) error
	Clear() error
}

type defaultCache struct {
	mx   sync.RWMutex
	data map[string]any
}

var _ Cache = (*defaultCache)(nil)

func NewDefaultCache() *defaultCache {
	return &defaultCache{
		mx:   sync.RWMutex{},
		data: make(map[string]any),
	}
}

func (c *defaultCache) Get(key string) (any, error) {
	c.mx.RLock()
	defer c.mx.RUnlock()
	value, ok := c.data[key]
	if !ok {
		return nil, ErrKeyNotFound
	}
	return value, nil
}

func (c *defaultCache) Set(key string, value any) error {
	c.mx.Lock()
	defer c.mx.Unlock()
	if _, ok := c.data[key]; ok {
		return ErrKeyExists
	}
	c.data[key] = value
	return nil
}

func (c *defaultCache) Delete(key string) error {
	c.mx.Lock()
	defer c.mx.Unlock()
	delete(c.data, key)
	return nil
}

func (c *defaultCache) Clear() error {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.data = make(map[string]any)
	return nil
}
