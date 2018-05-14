package Storage

import (
	"sync"
)

// KVData 键值对类型
//type KVData map[string]interface{}

//CachePool 缓存
type CachePool struct {
	rwm  *sync.RWMutex
	rwm2 *sync.RWMutex
	//            tableName     key-value
	cacheStruct map[string]map[string]interface{}
}

var cachePoolIns *CachePool

// CacheInst 返回实例
func CacheInst() *CachePool {
	if cachePoolIns == nil {
		cachePoolIns = newCachePool()
	}
	return cachePoolIns
}

// NewCachePool 初始化缓存池
func newCachePool() *CachePool {
	return &CachePool{
		rwm:         new(sync.RWMutex),
		rwm2:        new(sync.RWMutex),
		cacheStruct: make(map[string]map[string]interface{}),
	}
}

func (me *CachePool) getTypeMap(typeName string) map[string]interface{} {
	me.rwm.RLock()
	typeStruct, ok := me.cacheStruct[typeName]
	me.rwm.RUnlock()
	if ok == false {
		me.rwm.Lock()
		typeStruct, ok = me.cacheStruct[typeName]
		if ok == false {
			typeStruct = make(map[string]interface{})
			me.cacheStruct[typeName] = typeStruct
		}
		me.rwm.Unlock()
	}
	return typeStruct
}

func (me *CachePool) getPersonalMap(typeName string, personalID string) map[string]interface{} {
	typeStruct := me.getTypeMap(typeName)

	me.rwm2.RLock()
	personalStruct, ok := typeStruct[personalID]
	me.rwm2.RUnlock()
	if ok == false {
		personalStruct = make(map[string]interface{})
		me.rwm2.Lock()
		typeStruct[personalID] = personalStruct
		me.rwm2.Unlock()
	}
	return personalStruct.(map[string]interface{})
}

// AddSharedData 添加共享数据到缓存池
func (me *CachePool) AddSharedData(typeName string, key string, data interface{}) {
	typeStruct := me.getTypeMap(typeName)
	//操作共享数据通过chan,不加锁
	typeStruct[key] = data
}

// AddPersonalData 添加私有数据到缓存池
func (me *CachePool) AddPersonalData(typeName string, personalID string, key string, data interface{}) {

	personalStruct := me.getPersonalMap(typeName, personalID)
	if personalStruct != nil {
		personalStruct[key] = data
	}
}

// FindSharedData 返回表里公共数据
func (me *CachePool) FindSharedData(typeName string, key string) interface{} {
	shareStruct := me.getTypeMap(typeName)
	if shareStruct == nil {

	}
	return shareStruct[key]
}

// FindPersonalData 查询私有数据到缓存池
func (me *CachePool) FindPersonalData(typeName string, personalID string, key string) interface{} {
	personalStruct := me.getPersonalMap(typeName, personalID)
	if personalStruct == nil {
		//
		SelectData(typeName, nil, personalID)
	}
	return personalStruct[key]
}

//DisposeCache 释放冷数据
func (me *CachePool) DisposeCache() {
}
