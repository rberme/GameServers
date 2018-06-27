package concurrentmap

import (
	"sync"
)

// ConcurrentMap 线程安全的字典
type ConcurrentMap struct {
	sync.RWMutex
	values map[interface{}]interface{}
}

// New .
func New() *ConcurrentMap {
	return &ConcurrentMap{values: make(map[interface{}]interface{})}
}

// Clear .
func (me *ConcurrentMap) Clear() {
	me.Lock()
	defer me.Unlock()

	for key := range me.values {
		delete(me.values, key)
	}
}

// Count .
func (me *ConcurrentMap) Count() int {
	me.RLock()
	defer me.RUnlock()

	return len(me.values)
}

// Put .
func (me *ConcurrentMap) Put(_key, _value interface{}) {
	me.Lock()
	defer me.Unlock()

	me.values[_key] = _value
}

// Remove .
func (me *ConcurrentMap) Remove(_key interface{}) {
	me.Lock()
	defer me.Unlock()

	delete(me.values, _key)
}

// Contains .
func (me *ConcurrentMap) Contains(_key interface{}) bool {
	me.RLock()
	defer me.RUnlock()

	_, ok := me.values[_key]

	return ok
}

// Get .
func (me *ConcurrentMap) Get(_key interface{}) (interface{}, bool) {
	me.RLock()
	defer me.RUnlock()

	value, ok := me.values[_key]
	return value, ok
}
