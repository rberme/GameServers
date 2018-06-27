package main

import "serializer"

//go:generate F:/Golang/GoPath/bin/stringer -type=MsgCodeType

var ser = serializer.ProtobufSerializer{}

// MsgCodeType 消息类型
type MsgCodeType int

// Convert .
func (me MsgCodeType) Convert(id int) string {
	return (MsgCodeType)(id).String()
}

// 协议
const (
	msgCodeLogin MsgCodeType = iota + 101 //登录协议请求
	msgCodeLoginRET
)

// 协议
const (
	msgCodePlayer MsgCodeType = iota + 1101
	msgCodePlayerRET
)

// MsgHandler 消息处理函数
type MsgHandler struct {
	MsgCodeType
}

// MsgCodeLogin .
func (me *MsgHandler) MsgCodeLogin(id uint64, bytes []byte) {
	//sess := nnet.Hub.Get(id)
}
