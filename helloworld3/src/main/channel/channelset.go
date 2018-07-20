package channel

import (
	"sync"
)

type channelSet struct {
	sync.RWMutex
	data map[int64]bool
}

func (me *channelSet) Length() int {
	me.RLock()
	defer me.RUnlock()
	return len(me.data)
}

func (me *channelSet) Add(id int64) {
	me.Lock()
	defer me.Unlock()
	me.data[id] = true
}

func (me *channelSet) Foreach(action func(id int64)) {
	me.RLock()
	defer me.RUnlock()
	if action != nil {
		for id := range me.data {
			action(id)
		}
	}
}

func (me *channelSet) Del(id int64) {
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

func (me *channelMap) Add(channel int32, id int64) {
	temp := &channelSet{
		data: make(map[int64]bool),
	}
	me.Data[channel] = temp
	temp.Add(id)
}

func (me *channelMap) Del(channel int32) {
	delete(me.Data, channel)
}
