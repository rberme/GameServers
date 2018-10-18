package protocol

import (
	"sync"
)

// TCPCmdPacketEventHandler .
type TCPCmdPacketEventHandler func(this *TCPInPacket)

// TCPInPacketPool .
type TCPInPacketPool struct {
	//sync.RWMutex
	pool sync.Pool //[]*TCPInPacket
}

// NewTCPInPacketPool .
func NewTCPInPacketPool(capacity int) *TCPInPacketPool {
	retval := &TCPInPacketPool{}
	retval.pool.New = func() interface{} {
		return &TCPInPacket{
		//TCPCmdPacketEvent: TCPCmdPacketEvent,
		}
	}
	return retval
}

// Count .
// func (me *TCPInPacketPool) Count() int {
// 	return len(me.pool)
// }

// Pop .
func (me *TCPInPacketPool) Pop(s iSocket, TCPCmdPacketEvent TCPCmdPacketEventHandler) *TCPInPacket {

	var tcpInPacket *TCPInPacket
	tcpInPacket = me.pool.Get().(*TCPInPacket)
	tcpInPacket.CurrentSocket = s
	tcpInPacket.TCPCmdPacketEvent = TCPCmdPacketEvent
	return tcpInPacket
}

// Push .
func (me *TCPInPacketPool) Push(item *TCPInPacket) {
	if item != nil {
		item.Reset()
		me.pool.Put(item)
	}
}
