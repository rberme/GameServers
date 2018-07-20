package network

import (
	"net"
)

// ConnSet 类型定义
type ConnSet map[net.Conn]struct{}

// TCPConn .
type TCPConn struct {
	//sync.Mutex
	conn net.Conn
	//writeChan chan []byte
	//closeFlag bool
	msgParser *MsgParser
}

func newTCPConn(conn net.Conn, msgParser *MsgParser) *TCPConn {
	tcpConn := new(TCPConn)
	tcpConn.conn = conn
	//tcpConn.writeChan = make(chan []byte, pendingWriteNum)
	tcpConn.msgParser = msgParser

	// go func() {
	// 	for b := range tcpConn.writeChan {
	// 		if b == nil {
	// 			break
	// 		}

	// 		_, err := conn.Write(b)
	// 		if err != nil {
	// 			break
	// 		}
	// 	}

	// 	conn.Close()
	// 	tcpConn.Lock()
	// 	tcpConn.closeFlag = true
	// 	tcpConn.Unlock()
	// }()

	return tcpConn
}

// func (me *TCPConn) doDestroy() {
// 	me.conn.(*net.TCPConn).SetLinger(0)
// 	me.conn.Close()

// 	// if !me.closeFlag {
// 	// 	close(me.writeChan)
// 	// 	me.closeFlag = true
// 	// }
// }

// // Destroy .
// func (me *TCPConn) Destroy() {
// 	me.Lock()
// 	defer me.Unlock()

// 	me.doDestroy()
// }

// Close .
func (me *TCPConn) Close() {
	me.conn.Close()
}

// b must not be modified by the others goroutines
func (me *TCPConn) Write(b []byte) {
	if b == nil {
		return
	}

	_, err := me.conn.Write(b)
	if err != nil {
		me.conn.Close()
	}
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
func (me *TCPConn) WriteMsg(buff []byte) error {
	return me.msgParser.Write(me, buff)
}

// // WriteMsg 通过写通道发送数据
// func (me *TCPConn) WriteMsg(msgcode int, bufs ...interface{}) error {
// 	//4+1+4+buf
// 	var bufLen uint32 //:= len(buf)

// 	for _, buf := range bufs {
// 		switch buf.(type) {
// 		case uint32:
// 			bufLen += 4
// 		case int32:
// 			bufLen += 4
// 		case int8:
// 			bufLen++
// 		case byte:
// 			bufLen++
// 		case string:
// 			b := []byte(buf.(string))
// 			//to C# binary
// 			l := len(b)
// 			if l < 128 {
// 				bufLen++
// 			} else if l < 65536 {
// 				bufLen += 2
// 			} else {
// 				llog.Error("WriteMsg:未处理的字符串长度")
// 			}
// 			bufLen += uint32(l)
// 		case []byte:
// 			bufLen += uint32(len(buf.([]byte)))
// 		case int64:
// 			bufLen += 8
// 		default:
// 			llog.Error("WriteMsg:未处理的类型")
// 		}
// 	}

// 	arg := make([]byte, 9+bufLen)
// 	binary.LittleEndian.PutUint32(arg[0:], uint32(msgcode))
// 	binary.LittleEndian.PutUint32(arg[5:], uint32(bufLen))

// 	l := 9
// 	for _, buf := range bufs {
// 		switch buf.(type) {
// 		case uint32:
// 			binary.LittleEndian.PutUint32(arg[l:], buf.(uint32))
// 			l += 4
// 		case int32:
// 			binary.LittleEndian.PutUint32(arg[l:], uint32(buf.(int32)))
// 			l += 4
// 		case int8:
// 			arg[l] = byte(buf.(int8))
// 			l++
// 		case byte:
// 			arg[l] = buf.(byte)
// 			l++
// 		case string:
// 			b := []byte(buf.(string))
// 			//to C# binary
// 			le := len(b)
// 			if le < 128 {
// 				arg[l] = byte(le)
// 				l++
// 				copy(arg[l:], b)
// 				l += le
// 			} else if le < 65536 {
// 				arg[l] = byte(le%128 + 128)
// 				arg[l+1] = byte(le / 128)
// 				l += 2
// 				copy(arg[l:], b)
// 				l += le
// 			} else {
// 				llog.Error("WriteMsg:未处理的字符串长度")
// 			}
// 		case []byte:
// 			b := buf.([]byte)
// 			copy(arg[l:], b)
// 			l += len(b)
// 		case int64:
// 			binary.LittleEndian.PutUint64(arg[l:], uint64(buf.(int64)))
// 			l += 8
// 		default:
// 		}
// 	}

// 	return me.msgParser.Write(me, arg)
// }
