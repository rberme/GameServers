package main

import (
	"sync"
)

type channelSet struct {
	sync.RWMutex
	data map[int64]bool
}

func (me *channelSet) length() int {
	me.RLock()
	defer me.RUnlock()
	return len(me.data)
}

func (me *channelSet) add(id int64) {
	me.Lock()
	defer me.Unlock()
	me.data[id] = true
}

func (me *channelSet) del(id int64) {
	me.Lock()
	defer me.Unlock()
	delete(me.data, id)
}

////////////////////////////////////////////////////////////////

type channelMap struct {
	sync.RWMutex
	Data map[int32]*channelSet
}

func newChannelMap() *channelMap {
	return &channelMap{
		Data: make(map[int32]*channelSet),
	}
}

// func (me *channelMap) length() int {
// 	me.RLock()
// 	defer me.RUnlock()
// 	return len(me.Data)
// }

func (me *channelMap) add(channel int32, id int64) {
	temp := &channelSet{
		data: make(map[int64]bool),
	}
	me.Data[channel] = temp
	temp.add(id)
}

// func (me *channelMap) del(channel int32) {
// 	me.Lock()
// 	defer me.Unlock()
// 	delete(me.Data, channel)
// }
