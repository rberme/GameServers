package Server

import (
	"Easy/BufferUtils"
	"Easy/Storage"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"
)

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//TCPServer Tcp服务器
type TCPServer struct {
	headSize int
	lsn      net.Listener

	sessions *ConnMap
	mainChan chan *chanMsg

	msgHandler IHandler
}

//NewServer 新建服务器对象
func NewServer() *TCPServer {
	return &TCPServer{
		headSize: 4,
		sessions: NewConnMap(),
		mainChan: make(chan *chanMsg),
	}
}

// 心跳
func (me *TCPServer) heartbeat() {
	num := 1
	for {
		time.Sleep(time.Second * 20)
		now := time.Now().Unix()
		timebuff := make([]byte, 8)
		binary.LittleEndian.PutUint64(timebuff, uint64(now))
		heartBuff := BufferUtils.AppendNumBytes(heartBeat, timebuff)
		heartBuff = BufferUtils.AppendHeadBytes(heartBuff)

		var deadID []int
		me.SendAllChan(heartBuff, func(sess *TCPConn) bool {
			if sess.NoDealHeartbeat > 0 {
				deadID = append(deadID, sess.ID)
				return false
			}
			sess.NoDealHeartbeat++
			return true
		})
		if len(deadID) > 0 {
			me.sessions.ClearConn(deadID)
		}

		if num == 6 {
			num = 0
			//
		}
		num++
	}
}

//公共消息处理
func (me *TCPServer) mainLoop() {
	for {
		msg := <-me.mainChan
		if msg.msgCode == recvMsg {
			userID := msg.userID
			if m, t, k := me.msgHandler.MainHandler(userID, msg.data); m {
				Storage.DataSyncInst().SyncEntity(t, k, 0)
			}
		}
	}
}

//私有消息处理
func (me *TCPServer) sessLoop(sess *TCPConn) {
	defer func() {
		fmt.Println("处理数据线程结束")
		recover()
	}()
	for {
		msg := <-sess.RecvChan
		switch msg.msgCode {
		case recvMsg:
			if m, t, k, p := me.msgHandler.ClientHandler(msg.data); m {
				Storage.DataSyncInst().SyncEntity(t, k, p)
			}
		case m2CSend: //发送数据
			sendLen, err := sess.conn.Write(msg.data)
			fmt.Println("发送数据长度:", sendLen)
			if err != nil {
				fmt.Println("数据发送失败", err) //write tcp 127.0.0.1:8887->127.0.0.1:50604: wsasend: An existing connection was forcibly closed by the remote host.
			}
		// case DataReq: //请求读取数据

		// 	msg.respChan <- &DataResp{
		// 		UserID:   msg.userID,
		// 		TypeName: msg.typeName,
		// 		DataKey:  msg.dataKey,
		// 	}
		default:
		}
		msg = nil
	}
}

func (me *TCPServer) sessRecv(sess *TCPConn) {
	defer func() {
		me.sessions.RemoveConn(sess)
		me.msgHandler.UserLogout(sess.ID)
		fmt.Println("接收数据线程结束")
		recover()
	}()

	bufflen := []byte{0, 0, 0, 0}
	for { //接收客户端数据
		_, err := sess.conn.Read(bufflen)
		if err != nil {
			return
		}
		bodySize := int(binary.LittleEndian.Uint32(bufflen)) - 4
		buf := make([]byte, bodySize)
		_, err = sess.conn.Read(buf)
		if err != nil {
			return
		}
		msgCode := int(binary.LittleEndian.Uint32(buf))

		if msgCode == heartBeatRet { //处理返回的心跳
			sess.NoDealHeartbeat = 0
			continue
		}

		if msgCode < 1500000 { //私有操作
			sess.RecvChan <- &chanMsg{
				msgCode: recvMsg,
				data:    buf,
			}
		} else { //主循环操作
			me.mainChan <- &chanMsg{
				msgCode: recvMsg,
				data:    buf,
				userID:  sess.ID,
			}
		}
	}
}

//Run 服务器运行
func (me *TCPServer) Run(address string, handler IHandler) {
	ls, err := net.Listen("tcp", address)
	me.lsn = ls
	if err != nil {
		return
	}
	fmt.Println("服务器启动:", address)
	me.msgHandler = handler
	go me.mainLoop()
	go me.heartbeat()

	for {
		c, err := ls.Accept()
		if err != nil {
			log.Println("Accept failed", err)
			break
		}
		sess := &TCPConn{
			// ID:         me.sessions.Length(),
			conn:            c,
			RecvChan:        make(chan *chanMsg),
			remoteAddr:      c.RemoteAddr().String(),
			NoDealHeartbeat: 0,
		}

		me.sessions.Add(sess)
		handler.UserLogin(sess.ID)

		go me.sessLoop(sess)
		go me.sessRecv(sess)
	}
}

// SendChan 向某个玩家发送
func (me *TCPServer) SendChan(useID int, buff []byte) {
	msg := &chanMsg{
		msgCode: m2CSend,
		data:    buff,
	}
	me.sessions.Get(useID).RecvChan <- msg
}

// SendAllChan 向所有客户端chan发送消息
func (me *TCPServer) SendAllChan(buff []byte, predicate func(sess *TCPConn) bool) {
	me.sendAllChan(m2CSend, buff, predicate)
}

func (me *TCPServer) sendAllChan(msgCode int, buff []byte, predicate func(sess *TCPConn) bool) {
	msg := &chanMsg{
		msgCode: msgCode,
		data:    buff,
	}
	if predicate == nil {
		me.sessions.ForeachRead(func(sess *TCPConn) {
			sess.RecvChan <- msg
		})
	} else {
		me.sessions.ForeachRead(func(sess *TCPConn) {
			if predicate(sess) == true {
				sess.RecvChan <- msg
			}
		})
	}
}

func (me *TCPServer) sendChan(msgCode int, buff []byte, ids ...int) {
	leng := len(ids)
	if leng == 0 {
		return
	}
	msg := &chanMsg{
		msgCode: msgCode,
		data:    buff,
	}
	if leng == 1 {
		sess := me.sessions.Get(ids[0])
		sess.RecvChan <- msg
	} else {
		me.sessions.ForeachRead(func(sess *TCPConn) {
			sess.RecvChan <- msg
		}, ids...)
	}
}

// PrintStatus 服务器状态
func (me *TCPServer) PrintStatus() {
	fmt.Println("当前sess数量:", me.sessions.Length())
}
