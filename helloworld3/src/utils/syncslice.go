package utils

import (
	"sync"
)

type SyncSlice struct {
	sync.RWMutex
	Data []int64
}

func (me *SyncSlice) Len() int {
	return len(me.Data)
}
