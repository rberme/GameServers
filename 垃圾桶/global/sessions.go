package global

import (
	"errors"
	"net"
	"sync"
	"sync/atomic"
)

// ErrClosed .
var ErrClosed = errors.New("link.Session closed")

// Session .
type Session struct {
	ID   uint64
	conn net.Conn
	// encoder         Encoder
	// decoder         Decoder
	//closeChan chan int
	closeFlag int32
	//closeEventMutex sync.Mutex
	//closeCallbacks  *list.List
	State interface{}
}

var globalSessionID uint64

// NewSession .
func NewSession(conn net.Conn, codecType CodecType) *Session {
	session := &Session{
		ID:   atomic.AddUint64(&globalSessionID, 1),
		conn: conn,
		// encoder:        codecType.NewEncoder(conn),
		// decoder:        codecType.NewDecoder(conn),
		// closeCallbacks: list.New(),
	}
	return session
}

// NewSessionByID .
func NewSessionByID(conn net.Conn, codecType CodecType, sessionID uint64) *Session {
	session := &Session{
		ID:   sessionID,
		conn: conn,
		// encoder:        codecType.NewEncoder(conn),
		// decoder:        codecType.NewDecoder(conn),
		//closeCallbacks: list.New(),
	}
	return session
}

// func (session *Session) Id() uint64     { return session.id }
// func (session *Session) Conn() net.Conn { return session.conn }

// IsClosed .
func (me *Session) IsClosed() bool { return atomic.LoadInt32(&me.closeFlag) == 1 }

// Close .
func (me *Session) Close() {
	if atomic.CompareAndSwapInt32(&me.closeFlag, 0, 1) {
		me.conn.Close()
	}
}

// Receive .
func (me *Session) Receive(msg interface{}) (err error) {
	if me.IsClosed() {
		return ErrClosed
	}
	// err = me.decoder.Decode(msg)
	// if err != nil {
	// 	me.Close()
	// }
	return
}

// Send .
func (me *Session) Send(msg interface{}) (err error) {
	if me.IsClosed() {
		return ErrClosed
	}
	// err = me.encoder.Encode(msg)
	// if err != nil {
	// 	me.Close()
	// }
	return
}

// type closeCallback struct {
// 	Handler interface{}
// 	Func    func()
// }

// func (session *Session) AddCloseCallback(handler interface{}, callback func()) {
// 	if session.IsClosed() {
// 		return
// 	}

// 	session.closeEventMutex.Lock()
// 	defer session.closeEventMutex.Unlock()

// 	session.closeCallbacks.PushBack(closeCallback{handler, callback})
// }

// func (session *Session) RemoveCloseCallback(handler interface{}) {
// 	if session.IsClosed() {
// 		return
// 	}

// 	session.closeEventMutex.Lock()
// 	defer session.closeEventMutex.Unlock()

// 	for i := session.closeCallbacks.Front(); i != nil; i = i.Next() {
// 		if i.Value.(closeCallback).Handler == handler {
// 			session.closeCallbacks.Remove(i)
// 			return
// 		}
// 	}
// }

// func (session *Session) invokeCloseCallbacks() {
// 	session.closeEventMutex.Lock()
// 	defer session.closeEventMutex.Unlock()

// 	for i := session.closeCallbacks.Front(); i != nil; i = i.Next() {
// 		callback := i.Value.(closeCallback)
// 		callback.Func()
// 	}
// }

var (
	sessions     = make(map[uint64]*Session)
	sessionMutex sync.Mutex
)

// AddSession .
func AddSession(sess *Session) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	sessions[sess.ID] = sess
	// me.AddCloseCallback(session, func() {
	// 	RemoveSession(session.Id())
	// })
}

// RemoveSession .
func RemoveSession(key uint64) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()
	delete(sessions, key)
}

// GetSession .
func GetSession(key uint64) *Session {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	session, _ := sessions[key]
	return session
}

// SessionLen .
func SessionLen() int {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	return len(sessions)
}

// FetchSession .
func FetchSession(callback func(*Session)) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	for _, session := range sessions {
		callback(session)
	}
}
