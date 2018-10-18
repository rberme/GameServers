package server

import (
	"encoding/binary"
	"net"
	"server/protocol"
	"sync/atomic"
	"time"
)

// Socket 自定义套接字
type Socket struct {
	Used     int32
	ID       int64
	groupIdx int
	conn     net.Conn
	closed   int32
	startime int64
	//session     *Session
	TCPInPacket *protocol.TCPInPacket
}

// NewSocket .
func NewSocket(c net.Conn) *Socket {
	return &Socket{
		conn:     c,
		closed:   0,
		Used:     0,
		startime: time.Now().Unix(),
	}
}

// Close .
func (me *Socket) Close() {
	if atomic.CompareAndSwapInt32(&me.closed, 0, 1) {
		me.conn.Close()
	}
}

func (me *Socket) Write(buff []byte) error {
	msgLen := uint32(len(buff))
	msg := make([]byte, 4+msgLen)
	binary.LittleEndian.PutUint32(msg, msgLen+4)
	copy(msg[4:], buff)
	_, err := me.conn.Write(msg)
	return err
}
