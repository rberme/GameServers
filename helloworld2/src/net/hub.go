package net

import "sync"

const sessionMapNum = 32

// Hub 用户中心
var Hub = newHub()

type sessionMap struct {
	sync.RWMutex
	sessions map[uint64]*Session
	disposed bool
}

// Hub .
type Hub struct {
	sessionMaps [sessionMapNum]sessionMap
	disposeOnce sync.Once
	disposeWait sync.WaitGroup
}

// NewHub .
func newHub() *Hub {
	hub := &Hub{}
	for i := 0; i < len(hub.sessionMaps); i++ {
		hub.sessionMaps[i].sessions = make(map[uint64]*Session)
	}
	return hub
}

// Dispose 释放所有session
func (me *Hub) Dispose() {
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

func (me *Hub) putSession(session *Session) {
	smap := &me.sessionMaps[session.ID%sessionMapNum]

	smap.Lock()
	defer smap.Unlock()

	if smap.disposed {
		session.Close()
		return
	}

	smap.sessions[session.ID] = session
	me.disposeWait.Add(1)
}

func (me *Hub) delSession(session *Session) {
	smap := &me.sessionMaps[session.ID%sessionMapNum]

	smap.Lock()
	defer smap.Unlock()

	delete(smap.sessions, session.ID)
	me.disposeWait.Done()
}
