package network

import (
	"llog"
	"net"
	"sync"
)

// ConnSet 类型定义
type ConnSet map[net.Conn]struct{}

// TCPConn .
type TCPConn struct {
	sync.Mutex
	conn      net.Conn
	writeChan chan []byte
	closeFlag bool
	msgParser *MsgParser
}

func newTCPConn(conn net.Conn, pendingWriteNum int, msgParser *MsgParser) *TCPConn {
	tcpConn := new(TCPConn)
	tcpConn.conn = conn
	tcpConn.writeChan = make(chan []byte, pendingWriteNum)
	tcpConn.msgParser = msgParser

	go func() {
		for b := range tcpConn.writeChan {
			if b == nil {
				break
			}

			_, err := conn.Write(b)
			if err != nil {
				break
			}
		}

		conn.Close()
		tcpConn.Lock()
		tcpConn.closeFlag = true
		tcpConn.Unlock()
	}()

	return tcpConn
}

func (me *TCPConn) doDestroy() {
	me.conn.(*net.TCPConn).SetLinger(0)
	me.conn.Close()

	if !me.closeFlag {
		close(me.writeChan)
		me.closeFlag = true
	}
}

// Destroy .
func (me *TCPConn) Destroy() {
	me.Lock()
	defer me.Unlock()

	me.doDestroy()
}

// Close .
func (me *TCPConn) Close() {
	me.Lock()
	defer me.Unlock()
	if me.closeFlag {
		return
	}

	me.doWrite(nil)
	me.closeFlag = true
}

func (me *TCPConn) doWrite(b []byte) {
	if len(me.writeChan) == cap(me.writeChan) {
		llog.Debug("close conn: channel full")
		me.doDestroy()
		return
	}

	me.writeChan <- b
}

// b must not be modified by the others goroutines
func (me *TCPConn) Write(b []byte) {
	me.Lock()
	defer me.Unlock()
	if me.closeFlag || b == nil {
		return
	}

	me.doWrite(b)
}

func (me *TCPConn) Read(b []byte) (int, error) {
	return me.conn.Read(b)
}

// LocalAddr .
func (me *TCPConn) LocalAddr() net.Addr {
	return me.conn.LocalAddr()
}

// RemoteAddr .
func (me *TCPConn) RemoteAddr() net.Addr {
	return me.conn.RemoteAddr()
}

// ReadMsg .
func (me *TCPConn) ReadMsg() ([]byte, error) {
	return me.msgParser.Read(me)
}

// WriteMsg 通过写通道发送数据
func (me *TCPConn) WriteMsg(args ...[]byte) error {
	return me.msgParser.Write(me, args...)
}
