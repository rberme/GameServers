package utils

import (
	"sync"
)

type SyncMap struct {
	sync.RWMutex
	Data map[int32]interface{}
}
