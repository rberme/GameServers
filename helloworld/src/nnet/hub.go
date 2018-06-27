package nnet

import (
	"sync"
)

const sessionMapNum = 32

// Hub 用户中心
var Hub = newHub()

type sessionMap struct {
	sync.RWMutex
	sessions map[uint64]*Session
	disposed bool
}

// Hub .
type hub struct {
	sessionMaps [sessionMapNum]sessionMap
	disposeOnce sync.Once
	disposeWait sync.WaitGroup
}

// NewHub .
func newHub() *hub {
	h := &hub{}
	for i := 0; i < len(h.sessionMaps); i++ {
		h.sessionMaps[i].sessions = make(map[uint64]*Session)
	}
	return h
}

// Dispose 释放所有session
func (me *hub) Dispose() {
	me.disposeOnce.Do(func() {
		for i := 0; i < sessionMapNum; i++ {
			smap := &me.sessionMaps[i]
			smap.Lock()
			smap.disposed = true
			for _, session := range smap.sessions {
				//session.Close()
				me.delSession(session)
			}
			smap.Unlock()
		}
		me.disposeWait.Wait()
	})
}

func (me *hub) putSession(session *Session) {
	smap := &me.sessionMaps[session.ID%sessionMapNum]

	smap.Lock()
	defer smap.Unlock()

	if smap.disposed {
		session.close()
		return
	}

	smap.sessions[session.ID] = session
	me.disposeWait.Add(1)
}

func (me *hub) delSession(session *Session) {
	smap := &me.sessionMaps[session.ID%sessionMapNum]

	smap.Lock()
	defer smap.Unlock()
	session.close()
	delete(smap.sessions, session.ID)
	me.disposeWait.Done()
}

func (me *hub) Get(id uint64) *Session {
	smap := &me.sessionMaps[id%sessionMapNum]
	smap.RLock()
	defer smap.RUnlock()
	return smap.sessions[id]
}
