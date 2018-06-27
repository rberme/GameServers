package protos

import (
	"encoding/binary"
	"reflect"
	"tools"

	"github.com/golang/protobuf/proto"
)

// ISession .
type ISession interface {
	Send([]byte)
}

// ProtoMsg .
type ProtoMsg struct {
	ID             uint16
	Body           interface{}
	Identification uint64
}

// XXX
var (
	// NullProtoMsg 空消息
	NullProtoMsg = ProtoMsg{0, nil, 0}
	MsgObjectMap = make(map[uint16]reflect.Type)
	MsgIDMap     = make(map[reflect.Type]uint16)
)

// SetMsg 设置消息类型和消息ID的对应关系
func SetMsg(msgID uint16, data interface{}) {
	msgType := reflect.TypeOf(data)

	MsgObjectMap[msgID] = msgType
	MsgIDMap[reflect.TypeOf(reflect.New(msgType).Interface())] = msgID
}

// GetMsgObject 根据消息ID获取消息实体
func GetMsgObject(msgID uint16) proto.Message {
	if msgType, exists := MsgObjectMap[msgID]; exists {
		return reflect.New(msgType).Interface().(proto.Message)
	}
	tools.ERR("No MsgID:", msgID)
	return nil
}

// GetMsgID 根据一条消息获取消息ID
func GetMsgID(msg interface{}) uint16 {
	msgType := reflect.TypeOf(msg)
	if msgID, exists := MsgIDMap[msgType]; exists {
		return msgID
	} else {
		tools.ERR("No MsgType:", msgType)
	}
	return 0
}

// MarshalProtoMsg 序列化
func MarshalProtoMsg(args proto.Message) []byte {
	msgID := GetMsgID(args)
	msgBody, _ := proto.Marshal(args)

	result := make([]byte, 2+len(msgBody))
	binary.LittleEndian.PutUint16(result[:2], msgID)
	copy(result[2:], msgBody)

	return result
}

// UnmarshalProtoMsg 反序列化
func UnmarshalProtoMsg(msg []byte) ProtoMsg {
	if len(msg) < 2 {
		return NullProtoMsg
	}

	msgID := binary.LittleEndian.Uint16(msg[:2])
	msgBody := GetMsgObject(msgID)
	if msgBody == nil {
		return NullProtoMsg
	}

	err := proto.Unmarshal(msg[2:], msgBody)
	if err != nil {
		return NullProtoMsg
	}

	return ProtoMsg{
		ID:   msgID,
		Body: msgBody,
	}
}

// Send 发送消息
func Send(session *ISession, msgBody []byte) {
	(*session).Send(msgBody)
	//session.Send(msgBody)
}

// String 封装消息String类型字段
func String(param string) *string {
	return proto.String(param)
}

// Uint64 封装消息Uint64类型字段
func Uint64(param uint64) *uint64 {
	return proto.Uint64(param)
}

// Int64 封装消息Int64类型字段
func Int64(param int64) *int64 {
	return proto.Int64(param)
}

// Int32 封装消息Int32类型字段
func Int32(param int32) *int32 {
	return proto.Int32(param)
}

// Uint32 封装消息Uint32类型字段
func Uint32(param uint32) *uint32 {
	return proto.Uint32(param)
}
