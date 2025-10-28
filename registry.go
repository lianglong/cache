package cache

import (
	"fmt"
	"sync"
)

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]Constructor)
)

type Constructor func(config Config) (Cache, error)

// Register 由实现包的 init() 调用来注册自己
func Register(name string, c Constructor) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if _, exists := drivers[name]; exists {
		// 重复注册应该提示警告，但不应该 panic
		fmt.Printf("cache driver %q already registered\n", name)
		return
	}
	drivers[name] = c
}

// New 根据名字从注册表创建实例
func New(driver string, cfg Config) (Cache, error) {
	driversMu.RLock()
	defer driversMu.RUnlock()
	c, ok := drivers[driver]
	if !ok {
		return nil, fmt.Errorf("cache driver %q not found", driver)
	}
	return c(cfg)
}

// MustNew 创建缓存实例，失败则 panic（便捷方法）
func MustNew(driver string, cfg Config) Cache {
	cache, err := New(driver, cfg)
	if err != nil {
		panic(err)
	}
	return cache
}

// Drivers 返回已注册的驱动列表
func Drivers() []string {
	driversMu.RLock()
	defer driversMu.RUnlock()
	names := make([]string, 0, len(drivers))
	for name := range drivers {
		names = append(names, name)
	}
	return names
}
