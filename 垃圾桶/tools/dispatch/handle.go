package dispatch

import (
	"encoding/binary"
	"protos"
	"tools"
)

// IHandle .
type IHandle interface {
	DealMsg(session *ISession, msg []byte)
}

// Handle 通用Handle
type Handle map[uint16]func(*ISession, protos.ProtoMsg)

// DealMsg 处理消息
func (me Handle) DealMsg(session *ISession, msg []byte) {
	msgID := binary.LittleEndian.Uint16(msg[:2])
	var protoMsg protos.ProtoMsg
	// if systemProto.IsValidID(msgID) || logProto.IsValidID(msgID) || gameProto.IsValidID(msgID) {
	protoMsg = protos.UnmarshalProtoMsg(msg)
	// } else if dbProto.IsValidID(msgID) {
	// 	protoMsg = dbProto.UnmarshalProtoMsg(msg)
	// }

	if protoMsg == protos.NullProtoMsg {
		tools.ERR("收到Proto未处理消息：", msgID)
		return
	}

	if f, exists := me[msgID]; exists {
		f(session, protoMsg)
	} else {
		tools.ERR("收到Handle未处理消息：", msgID)
	}
}
