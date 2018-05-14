package Packet

import (
	"errors"
	"fmt"
)

// Type 网络包类型: 比如 handshake 等等.
type Type byte

const (
	_ Type = iota

	// Handshake represents a handshake: request(client) <====> handshake response(server)
	Handshake = 0x01

	// HandshakeAck represents a handshake ack from client to server
	HandshakeAck = 0x02

	// Heartbeat represents a heartbeat
	Heartbeat = 0x03

	// Data represents a common data packet
	Data = 0x04

	// Kick represents a kick off packet
	Kick = 0x05 // disconnect message from server
)

// ErrWrongPacketType 表示一个错误的包类型
var ErrWrongPacketType = errors.New("wrong packet type")

// Packet 表示一个网络数据包
type Packet struct {
	Type   Type
	Length int
	Data   []byte
}

//New 创建一个包的实例
func New() *Packet {
	return &Packet{}
}

//String 包的文本
func (p *Packet) String() string {
	return fmt.Sprintf("Type: %d, Length: %d, Data: %s", p.Type, p.Length, string(p.Data))
}
