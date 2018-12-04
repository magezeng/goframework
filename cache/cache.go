package cache

import "sync"

var instance Cache

// Cache 线程安全的缓存
type Cache struct {
	*sync.RWMutex
	internalCache map[string]interface{}
}

// 只会初始化一次，线程安全，类似ClassLoader
func init() {
	instance = Cache{internalCache: make(map[string]interface{})}
}

// Put 向缓存中写入一条数据，也可以进行覆盖
func (ca Cache) Put(key string, value interface{}) {
	ca.Lock()
	defer ca.Unlock()
	ca.internalCache[key] = value
}

// Get 从缓存中读取一条数据
func (ca Cache) Get(key string) interface{} {
	ca.RLock()
	defer ca.RUnlock()
	return ca.internalCache[key]
}

// Delete 从缓存中删除一条数据
func (ca Cache) Delete(key string) {
	ca.Lock()
	defer ca.Unlock()
	delete(ca.internalCache, key)
}

// GetInstance 取得缓存的实例（只读）
func GetInstance() Cache {
	return instance
}
