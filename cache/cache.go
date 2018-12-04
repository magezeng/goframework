package cache

import "sync"

var instance Cache

// Cache 线程安全的缓存
type Cache struct {
	lock          *sync.RWMutex
	internalCache map[string]interface{}
}

// 只会初始化一次，线程安全，类似ClassLoader
func init() {
	lock := new(sync.RWMutex)
	instance = Cache{internalCache: make(map[string]interface{}), lock: lock}
}

// Put 向缓存中写入一条数据，也可以进行覆盖
func (ca Cache) Put(key string, value interface{}) {
	ca.lock.Lock()
	defer ca.lock.Unlock()
	ca.internalCache[key] = value
}

// Get 从缓存中读取一条数据
func (ca Cache) Get(key string) interface{} {
	ca.lock.RLock()
	defer ca.lock.RUnlock()
	return ca.internalCache[key]
}

// Delete 从缓存中删除一条数据
func (ca Cache) Delete(key string) {
	ca.lock.Lock()
	defer ca.lock.Unlock()
	delete(ca.internalCache, key)
}

// GetInstance 取得缓存的实例（只读）
func GetInstance() Cache {
	return instance
}
