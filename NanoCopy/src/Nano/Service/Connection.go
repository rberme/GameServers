package Service

import (
	"sync/atomic"
)

// Connections 供session使用的全局变量
var Connections = newConnectionService()

type connectionService struct {
	count int64
	sid   int64
}

func newConnectionService() *connectionService {
	return &connectionService{
		sid: 0,
	}
}

func (me *connectionService) Increment() {
	atomic.AddInt64(&me.count, 1)
}

func (me *connectionService) Decrement() {
	atomic.AddInt64(&me.count, -1)
}

func (me *connectionService) Count() int64 {
	return atomic.LoadInt64(&me.count)
}

func (me *connectionService) Reset() {
	atomic.StoreInt64(&me.count, 0)
	atomic.StoreInt64(&me.sid, 0)
}

func (me *connectionService) SessionID() int64 {
	return atomic.AddInt64(&me.sid, 1)
}
