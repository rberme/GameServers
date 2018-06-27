package network

import (
	"llog"
	"net"
	"sync"
	"time"
)

// TCPServer .
type TCPServer struct {
	Addr            string
	MaxConnNum      int //最大连接数
	PendingWriteNum int //写通道缓存大小
	NewAgent        func(*TCPConn) Agent
	ln              net.Listener
	conns           ConnSet //map net.Conn
	mutexConns      sync.Mutex
	wgLn            sync.WaitGroup //监听的数量
	wgConns         sync.WaitGroup

	// msg parser 数据解析
	LenMsgLen    int
	MinMsgLen    uint32
	MaxMsgLen    uint32
	LittleEndian bool
	msgParser    *MsgParser
}

// Start .
func (me *TCPServer) Start() {
	me.init()
	go me.run()
}

// init 创建Listen 初始化消息解析器
func (me *TCPServer) init() {
	ln, err := net.Listen("tcp", me.Addr)
	if err != nil {
		llog.Fatal("%v", err)
	}

	if me.MaxConnNum <= 0 {
		me.MaxConnNum = 100
		llog.Release("invalid MaxConnNum, reset to %v", me.MaxConnNum)
	}
	if me.PendingWriteNum <= 0 {
		me.PendingWriteNum = 100
		llog.Release("invalid PendingWriteNum, reset to %v", me.PendingWriteNum)
	}
	if me.NewAgent == nil {
		llog.Fatal("NewAgent must not be nil")
	}

	me.ln = ln
	me.conns = make(ConnSet)

	// msg parser
	msgParser := NewMsgParser()
	msgParser.SetMsgLen(me.LenMsgLen, me.MinMsgLen, me.MaxMsgLen)
	msgParser.SetByteOrder(me.LittleEndian)
	me.msgParser = msgParser
}

// run 开始接收客户端连接
func (me *TCPServer) run() {
	me.wgLn.Add(1)
	defer me.wgLn.Done()

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
				llog.Release("accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return
		}
		tempDelay = 0

		me.mutexConns.Lock()
		if len(me.conns) >= me.MaxConnNum {
			me.mutexConns.Unlock()
			conn.Close()
			llog.Debug("too many connections")
			continue
		}
		me.conns[conn] = struct{}{}
		me.mutexConns.Unlock()

		me.wgConns.Add(1)

		tcpConn := newTCPConn(conn, me.PendingWriteNum, me.msgParser)
		agent := me.NewAgent(tcpConn)
		go func() {
			agent.Run()

			// cleanup
			tcpConn.Close()
			me.mutexConns.Lock()
			delete(me.conns, conn)
			me.mutexConns.Unlock()
			agent.OnClose()

			me.wgConns.Done()
		}()
	}
}

// Close .
func (me *TCPServer) Close() {
	me.ln.Close()
	me.wgLn.Wait()

	me.mutexConns.Lock()
	for conn := range me.conns {
		conn.Close()
	}
	me.conns = nil
	me.mutexConns.Unlock()
	me.wgConns.Wait()
}
