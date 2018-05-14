package Nano

import (
	"Nano/Internal/Codec"
	"Nano/Internal/Message"
	"Nano/Internal/Packet"
	"Nano/Session"
	"errors"
	"fmt"
	"net"
	"reflect"
	"sync/atomic"
	"time"
)

const (
	//写的缓存
	agentWriteBacklog = 16
)

var (
	// ErrBrokenPipe 表示低层连接被破坏
	ErrBrokenPipe = errors.New("broken low-level pipe")
	// ErrBufferExceed 表示当前session的buf已经满了 不能再接收更多数据
	ErrBufferExceed = errors.New("session send buffer exceed")
)

type (
	agent struct {
		session *Session.Session    // session
		conn    net.Conn            // 低层 conn fd
		lastMid uint                // 上一个消息id
		state   int32               //当前agent状态
		chDie   chan struct{}       //等待关闭
		chSend  chan pendingMessage //推送消息的队列
		lastAt  int64
		decoder *Codec.Decoder
		options *options

		srv reflect.Value // cached session reflect.Value
	}

	pendingMessage struct {
		typ     Message.Type // message type
		route   string       // message route (push)
		mid     uint         // 回应的消息id(response)
		payload interface{}  // payload
	}
)

// newAgent 实例化agent
func newAgent(conn net.Conn, options *options) *agent {
	a := &agent{
		conn:    conn,
		state:   statusStart,
		chDie:   make(chan struct{}),
		lastAt:  time.Now().Unix(),
		chSend:  make(chan pendingMessage, agentWriteBacklog),
		decoder: Codec.NewDecoder(),
		options: options,
	}
	s := Session.New(a)
	a.session = s
	a.srv = reflect.ValueOf(s)
	return a
}

func (me *agent) send(m pendingMessage) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = ErrBrokenPipe
		}
	}()
	me.chSend <- m
	return
}

func (me *agent) MID() uint {
	return me.lastMid
}

// Push, implementation for session.NetworkEntity interface
func (me *agent) Push(route string, v interface{}) error {
	if me.status() == statusClosed {
		return ErrBrokenPipe
	}

	if len(me.chSend) >= agentWriteBacklog {
		return ErrBufferExceed
	}

	if env.debug {
		switch d := v.(type) {
		case []byte:
			logger.Println(fmt.Sprintf("Type=Push, ID=%d, UID=%d, Route=%s, Data=%dbytes",
				me.session.ID(), me.session.UID(), route, len(d)))
		default:
			logger.Println(fmt.Sprintf("Type=Push, ID=%d, UID=%d, Route=%s, Data=%+v",
				me.session.ID(), me.session.UID(), route, v))
		}
	}

	return me.send(pendingMessage{
		typ:     Message.Push,
		route:   route,
		payload: v,
	})
}

// Response, implementation for session.NetworkEntity interface
// Response message to session
func (me *agent) Response(v interface{}) error {
	return me.ResponseMID(me.lastMid, v)
}

// Response, implementation for session.NetworkEntity interface
// Response message to session
func (me *agent) ResponseMID(mid uint, v interface{}) error {
	if me.status() == statusClosed {
		return ErrBrokenPipe
	}

	if mid <= 0 {
		return ErrSessionOnNotify
	}

	if len(me.chSend) >= agentWriteBacklog {
		return ErrBufferExceed
	}

	if env.debug {
		switch d := v.(type) {
		case []byte:
			logger.Println(fmt.Sprintf("Type=Response, ID=%d, UID=%d, MID=%d, Data=%dbytes",
				me.session.ID(), me.session.UID(), mid, len(d)))
		default:
			logger.Println(fmt.Sprintf("Type=Response, ID=%d, UID=%d, MID=%d, Data=%+v",
				me.session.ID(), me.session.UID(), mid, v))
		}
	}

	return me.send(pendingMessage{
		typ:     Message.Response,
		mid:     mid,
		payload: v,
	})
}

// Close, implementation for session.NetworkEntity interface
// Close closes the agent, clean inner state and close low-level connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (me *agent) Close() error {
	if me.status() == statusClosed {
		return ErrCloseClosedSession
	}
	me.setStatus(statusClosed)

	if env.debug {
		logger.Println(fmt.Sprintf("Session closed, ID=%d, UID=%d, IP=%s",
			me.session.ID(), me.session.UID(), me.conn.RemoteAddr()))
	}

	// prevent closing closed channel
	select {
	case <-me.chDie:
		// expect
	default:
		close(me.chDie)
		handler.chCloseSession <- me.session
	}

	return me.conn.Close()
}

// RemoteAddr, implementation for session.NetworkEntity interface
// returns the remote network address.
func (me *agent) RemoteAddr() net.Addr {
	return me.conn.RemoteAddr()
}

// String, implementation for Stringer interface
func (me *agent) String() string {
	return fmt.Sprintf("Remote=%s, LastTime=%d", me.conn.RemoteAddr().String(), me.lastAt)
}

func (me *agent) status() int32 {
	return atomic.LoadInt32(&me.state)
}

func (me *agent) setStatus(state int32) {
	atomic.StoreInt32(&me.state, state)
}

func (me *agent) write() {
	ticker := time.NewTicker(env.heartbeat)         //心跳计时器
	chWrite := make(chan []byte, agentWriteBacklog) //发送通道
	// clean func
	defer func() {
		ticker.Stop() //停止心跳计时器
		close(me.chSend)
		close(chWrite) //关闭发送通道
		me.Close()
		if env.debug {
			logger.Println(fmt.Sprintf("Session write goroutine exit, SessionID=%d, UID=%d", me.session.ID(), me.session.UID()))
		}
	}()

	for {
		select {
		case <-ticker.C:
			deadline := time.Now().Add(-2 * env.heartbeat).Unix()
			if me.lastAt < deadline { //两次心跳内无回应则表示连接超时
				logger.Println(fmt.Sprintf("Session heartbeat timeout, LastTime=%d, Deadline=%d", me.lastAt, deadline))
				return
			}
			chWrite <- hbd

		case data := <-chWrite:
			if _, err := me.conn.Write(data); err != nil {
				// 当连接断开后关闭代理
				logger.Println(err.Error())
				return
			}

		case data := <-me.chSend:
			payload, err := serializeOrRaw(data.payload)
			if err != nil {
				logger.Println(err.Error())
				break
			}

			// construct message and encode
			m := &Message.Message{
				Type:  data.typ,
				Data:  payload,
				Route: data.route,
				ID:    data.mid,
			}
			if pipe := me.options.pipeline; pipe != nil {
				err := pipe.Outbound().Process(me.session, *m)

				if err != nil {
					logger.Println("broken pipeline", err.Error())
					break
				}
			}

			em, err := m.Encode()
			if err != nil {
				logger.Println(err.Error())
				break
			}

			// packet encode
			p, err := Codec.Encode(Packet.Data, em)
			if err != nil {
				logger.Println(err)
				break
			}
			chWrite <- p

		case <-me.chDie: // agent closed signal
			return

		case <-env.die: // application quit
			return
		}
	}
}
