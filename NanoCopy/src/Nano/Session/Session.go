package Session

import (
	"Nano/Service"
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// INetworkEntity 低层的网络实例
type INetworkEntity interface {
	Push(route string, v interface{}) error
	MID() uint
	Response(v interface{}) error
	ResponseMID(mid uint, v interface{}) error
	Close() error
	RemoteAddr() net.Addr
}

var (
	// ErrIllegalUID 表示一个无效的uid
	ErrIllegalUID = errors.New("illegal uid")
)

// Session 在低层连接状态时可以保存临时数据的客户端session, 当低层连接中断后所有的数据会被释放
// 与客户端关联的session会传递给handler方法的第一个参数
type Session struct {
	sync.RWMutex
	id       int64
	uid      int64
	lastTime int64
	entity   INetworkEntity
	data     map[string]interface{}
}

// New 返回一个session实例
// INetworkEntity 是一个低层的网络实例
func New(entity INetworkEntity) *Session {
	return &Session{
		id:       Service.Connections.SessionID(),
		entity:   entity,
		data:     make(map[string]interface{}),
		lastTime: time.Now().Unix(),
	}
}

// Push 发送一个消息个客户端
func (me *Session) Push(route string, v interface{}) error {
	return me.entity.Push(route, v)
}

// Response 回复消息给客户端
func (me *Session) Response(v interface{}) error {
	return me.entity.Response(v)
}

// ResponseMID 返回消息给客户端，mid是消息id(message id)
func (me *Session) ResponseMID(mid uint, v interface{}) error {
	return me.entity.ResponseMID(mid, v)
}

// ID 返回sesssion id
func (me *Session) ID() int64 {
	return me.id
}

// UID 返回绑定到当前session的uid
func (me *Session) UID() int64 {
	return atomic.LoadInt64(&me.uid)
}

// MID 返回上一个消息的MID
func (me *Session) MID() uint {
	return me.entity.MID()
}

// Bind 绑定uid到当前session
func (me *Session) Bind(uid int64) error {
	if uid < 1 {
		return ErrIllegalUID
	}
	atomic.StoreInt64(&me.uid, uid)
	return nil
}

// Close 中断当前session,session关联的数据不会被释放
// all related data should be Clear explicitly in Session closed callback
func (me *Session) Close() {
	me.entity.Close()
}

// RemoteAddr 返回客户端网络地址
func (me *Session) RemoteAddr() net.Addr {
	return me.entity.RemoteAddr()
}
