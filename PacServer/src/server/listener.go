package server

import (
	"net"
	"sync"
	"sync/atomic"
	"time"
)

const (
	connGroup = 32
)

// Listener 服务器监听
type listener struct {
	ln       net.Listener
	connLock [connGroup]sync.RWMutex
	conns    [connGroup]map[*Socket]bool

	/// 断开成功通知函数
	SocketClosed func(s *Socket)
	/// 连接成功通知函数
	SocketConnected func(s *Socket)
	/// 接收数据通知函数
	SocketReceived func(s *Socket) error
	//processReceive func(s *Socket)
	/// 发送数据通知函数
	SocketSended func(s *Socket)
}

// NewListener .
func newListener() *listener {
	retval := &listener{}
	for i := 0; i < connGroup; i++ {
		retval.conns[i] = make(map[*Socket]bool)
	}
	return retval
}

// Init .
func (me *listener) Init() {

}

// Start 开始监听
func (me *listener) Start(ipport string) {
	ln, err := net.Listen("tcp", ipport)
	if err != nil {
		return
		//llog.Fatal("%v", err)
	}

	me.ln = ln
}

func (me *listener) startAccept() {
	connCount := 0
	var tempDelay time.Duration
	for {
		conn, err := me.ln.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				//llog.Release("accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return
		}
		tempDelay = 0
		go me.processAccept(conn, connCount)
		connCount++

	}
}

func (me *listener) processAccept(conn net.Conn, idx int) {
	s := NewSocket(conn)
	s.groupIdx = idx
	me.addSocket(s)

	if me.SocketConnected != nil {
		me.SocketConnected(s)
	}

	if me.SocketReceived != nil {
		me.SocketReceived(s)
	}

	me.CloseSocket(s)
}

// CloseSocket .
func (me *listener) CloseSocket(s *Socket) {
	if me.findSocket(s) == false {
		s.Close()
		return
	}
	me.removeSocket(s)

	// if me.SocketClosed != nil {
	// 	me.SocketClosed(s)
	// }
	s.Close()
}

func (me *listener) addSocket(s *Socket) {
	temp := s.groupIdx % connGroup
	me.connLock[temp].Lock()
	defer me.connLock[temp].Unlock()
	me.conns[temp][s] = true
}

func (me *listener) findSocket(s *Socket) bool {
	temp := s.groupIdx % connGroup
	me.connLock[temp].RLock()
	defer me.connLock[temp].RUnlock()
	_, ok := me.conns[temp][s]
	return ok
}

func (me *listener) removeSocket(s *Socket) {
	temp := s.groupIdx % connGroup
	me.connLock[temp].Lock()
	defer me.connLock[temp].Unlock()
	delete(me.conns[temp], s)
}

// func (me *listener) Length() int {
// 	length := 0
// 	for i := 0; i < connGroup; i++ {
// 		me.connLock[i].RLock()
// 		length += len(me.conns[i])
// 		me.connLock[i].RUnlock()
// 	}
// 	return length
// }

func (me *listener) ClearUnuseSocket() (length int) {
	length = 0
	now := time.Now().Unix()
	for i := 0; i < connGroup; i++ {
		me.connLock[i].RLock()
		for k := range me.conns[i] {
			used := atomic.LoadInt32(&k.Used)
			startime := k.startime
			if used == 0 && now-startime > 5 {
				k.Close()
				delete(me.conns[i], k)
			}
		}
		length += len(me.conns[i])
		me.connLock[i].RUnlock()
	}
	return
}
