package cache

import (
	"fmt"
	"sync"
)

type Cache interface {
	Get(string) ([]byte, error)
	SetCachePrefix(string)
	Save(string, []byte, ...int) bool
	Delete(string) bool
	Exists(string) bool
	SaveAll(map[string][]byte, ...int) bool
}

var adapters = make(map[string]AdapterBuilder)

var mu sync.RWMutex

type AdapterBuilder func(option Option) Cache

// 注册适配器
func Register(adapterName string, builder AdapterBuilder) {
	if builder == nil {
		panic("register cache adapter builder is nil")
	}
	if _, ok := adapters[adapterName]; ok {
		panic("cache adapter already exists")
	}
	mu.Lock()
	adapters[adapterName] = builder
	mu.Unlock()
}

func NewCache(adapterName string, option Option) (adapter Cache, err error) {
	mu.RLock()
	builder, ok := adapters[adapterName]
	mu.RUnlock()
	if !ok {
		err = fmt.Errorf("cache: unknown adapter name %q (forgot to import?)", adapterName)
		return
	}
	adapter = builder(option)
	return
}