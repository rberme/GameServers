package server

import (
	"fmt"
	"server/protocol"
)

// TCPManager .
var TCPManager = newTCPManager()
var tcpInPacketPool = protocol.NewTCPInPacketPool(100)

// ProcessCmd 消息处理
var ProcessCmd func(s *Socket, buff []byte) int

// ReceiveEnd .
var ReceiveEnd func(id int64)

// TCPManager .
type tcpManager struct {
	socketListener *listener
}

func newTCPManager() *tcpManager {
	me := &tcpManager{}
	me.socketListener = newListener()
	me.socketListener.SocketReceived = socketReceived

	return me
}

func (me *tcpManager) CloseSocket(s *Socket) {
	if s == nil {
		return
	}
	me.socketListener.CloseSocket(s)
}

// Start 启动服务器
func (me *tcpManager) Start(ipport string) {
	me.socketListener.Start(ipport)
	me.socketListener.startAccept()
}

// func (me *tcpManager) Status() (socketNum int) {
// 	socketNum = me.socketListener.Length()
// 	return
// }

func (me *tcpManager) ClearUselessSocket() (socketNum int) {
	socketNum = me.socketListener.ClearUnuseSocket()
	return
}

func socketReceived(s *Socket) error {
	defer func() {
		if ReceiveEnd != nil {
			ReceiveEnd(s.ID)
		}

		if err := recover(); err != nil {
			fmt.Println(err) //
		}
	}()

	buff := make([]byte, 2048)
	for {
		len, err := s.conn.Read(buff)
		if err != nil || len <= 0 {
			return err
		}

		tcpInPacket := s.TCPInPacket
		if tcpInPacket == nil {
			//从池里面获取
			tcpInPacket = tcpInPacketPool.Pop(s, tcpCmdPacketEvent)
		}
		tcpInPacket.WriteData(buff, len)
	}

}

// tcpCmdPacketEvent 命令包接收完毕后的回调事件
func tcpCmdPacketEvent(packet *protocol.TCPInPacket) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err) //
		}
	}()
	socket := packet.CurrentSocket.(*Socket)
	if ProcessCmd != nil {
		ProcessCmd(socket, packet.PacketBytes)
	}
}
