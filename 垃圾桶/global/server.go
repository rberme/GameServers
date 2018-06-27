package global

import (
	"io"
	"net"
	"time"
)

// Codec 数据传输编码解码
type Codec interface {
	Receive() (interface{}, error)
	Send(interface{}) error
	Close() error
}

// Protocol .
type Protocol interface {
	NewCodec(rw io.ReadWriter) (Codec, error)
}

// // ClearSendChan .
// type ClearSendChan interface {
// 	ClearSendChan(<-chan interface{})
// }

// Server .
type Server struct {
	// manager      *Manager     //管理session
	// listener     net.Listener //监听
	// protocol     Protocol     //协议
	// handler      Handler      //处理消息
	// sendChanSize int

	listener  net.Listener
	codecType CodecType

	// About sessions
	// maxSessionId uint64
	// sessions     map[uint64]*Session
	// sessionMutex sync.RWMutex

	// About server start and stop
	// stopOnce sync.Once
	// stopWait sync.WaitGroup

	// Server state
	State interface{}
}

// CodecType .
type CodecType interface {
	NewEncoder(w io.Writer) Encoder
	NewDecoder(r io.Reader) Decoder
}

// Encoder .
type Encoder interface {
	Encode(msg interface{}) error
}

// Decoder .
type Decoder interface {
	Decode(msg interface{}) error
}

// NewServer .
func NewServer(listener net.Listener, codecType CodecType) *Server {
	return &Server{
		listener:  listener,
		codecType: codecType,
		//sessions:  make(map[uint64]*Session),
	}
}

// Accept .
func (me *Server) Accept() (*Session, error) {
	var tempDelay time.Duration
	for {
		conn, err := me.listener.Accept()
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
				time.Sleep(tempDelay)
				continue
			}
			return nil, err
		}
		tempDelay = 0
		return me.newSession(conn), nil
	}
}

// // Stop .
// func (me *Server) Stop() {
// 	me.stopOnce.Do(func() {
// 		me.listener.Close()
// 		me.closeSessions()
// 		me.stopWait.Wait()
// 	})
// }

// // GetSession .
// func (me *Server) GetSession(sessionID uint64) *Session {
// 	me.sessionMutex.RLock()
// 	defer me.sessionMutex.RUnlock()
// 	session, _ := me.sessions[sessionID]
// 	return session
// }

// Server .
func (me *Server) newSession(conn net.Conn) *Session {
	session := NewSession(conn, me.codecType)

	return session
}

// func (me *Server) putSession(session *Session) {
// 	me.sessionMutex.Lock()
// 	defer me.sessionMutex.Unlock()

// 	//session.AddCloseCallback(server, func() { server.delSession(session) })
// 	me.sessions[session.ID] = session
// 	me.stopWait.Add(1)
// }

// func (me *Server) delSession(session *Session) {
// 	me.sessionMutex.Lock()
// 	defer me.sessionMutex.Unlock()

// 	//session.RemoveCloseCallback(server)
// 	delete(me.sessions, session.ID)
// 	me.stopWait.Done()
// }

// func (me *Server) copySessions() []*Session {
// 	me.sessionMutex.Lock()
// 	defer me.sessionMutex.Unlock()

// 	sessions := make([]*Session, 0, len(me.sessions))
// 	for _, session := range me.sessions {
// 		sessions = append(sessions, session)
// 	}
// 	return sessions
// }

// func (me *Server) closeSessions() {
// 	// copy session to avoid deadlock
// 	sessions := me.copySessions()
// 	for _, session := range sessions {
// 		session.Close()
// 	}
// }

// Serve .
func Serve(network, address string, codecType CodecType) (*Server, error) {
	listener, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}
	return NewServer(listener, codecType), nil
}

// Listener 开启服务器监听
func Listener(network, address string, codecType CodecType) error {
	listener, err := Serve(network, address, codecType)
	if err != nil {
		return err
	}

	go func() {
		for {
			session, err := listener.Accept()
			if err != nil {
				break
			}

			AddSession(session)

			// if dispatch != nil {
			// 	go sessionReceive(session, dispatch)
			// }
			go sessionLoop(session)
		}
	}()
	return nil
}

func sessionLoop(session *Session) {
	defer func() {
		session.Close()
		RemoveSession(session.ID)
	}()

	var buff []byte
	for {
		n, err := session.conn.Read(buff)
		if err != nil {
			return
		}
		data := buff[:n]

	}
}

// func sessionReceive(session *Session, d dispatch.IDispatch) {
// 	for {
// 		var msg []byte
// 		if err := session.Receive(&msg); err != nil {
// 			break
// 		}
// 		d.Process(session, msg)
// 	}
// }
