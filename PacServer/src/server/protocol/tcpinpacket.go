package protocol

import "encoding/binary"

type iSocket interface {
	Close()
}

// TCPInPacket .
type TCPInPacket struct {
	CurrentSocket iSocket

	packetDataSize     int
	packetDataHaveSize int
	PacketBytes        []byte
	TCPCmdPacketEvent  TCPCmdPacketEventHandler
}

// func NewTCPInPacket() *TCPInPacket {
// 	return &TCPInPacket{}
// }

// WriteData 将收到的数据写入
func (me *TCPInPacket) WriteData(buffer []byte, bufLen int) bool {
	//bufLen := len(buffer)
	if me.PacketBytes == nil {
		datalen := int(binary.LittleEndian.Uint32(buffer)) - 4
		if datalen > 65536 {
			panic("数据长度异常")
		}
		me.packetDataSize = datalen
		me.PacketBytes = make([]byte, me.packetDataSize)

		//copy(buffer[4:], me.PacketBytes)
		copy(me.PacketBytes, buffer[4:bufLen])
		me.packetDataHaveSize = bufLen - 4

	} else {
		//copy(buffer, me.PacketBytes[me.packetDataHaveSize:])
		copy(me.PacketBytes[me.packetDataHaveSize:], buffer[:bufLen])
		me.packetDataHaveSize += bufLen
	}

	if me.packetDataSize <= me.packetDataHaveSize {
		if me.TCPCmdPacketEvent != nil {
			me.TCPCmdPacketEvent(me)
		}
		me.PacketBytes = nil
	}

	return true
}

// Reset .
func (me *TCPInPacket) Reset() {
	me.CurrentSocket.Close()
	me.PacketBytes = nil
	me.packetDataHaveSize = 0

	me.packetDataSize = 0
}
