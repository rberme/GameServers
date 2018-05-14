package Codec

import (
	"Nano/Internal/Packet"
	"bytes"
	"errors"
)

const (
	// HeadLength 数据包长度
	HeadLength = 4
	// MaxPacketSize 数据包最大长度
	MaxPacketSize = 64 * 1024
)

// ErrPacketSizeExcced  此错误用于编码解码
var ErrPacketSizeExcced = errors.New("codec: packet size exceed")

// Decoder 读取并解码网络数据片
type Decoder struct {
	buf  *bytes.Buffer
	size int
	typ  byte
}

// NewDecoder 返回一个新的网络数据解码器
func NewDecoder() *Decoder {
	return &Decoder{
		buf:  bytes.NewBuffer(nil),
		size: -1,
	}
}

func (me *Decoder) forward() error {
	header := me.buf.Next(HeadLength)
	me.typ = header[0]
	if me.typ < Packet.Handshake || me.typ > Packet.Kick {
		return Packet.ErrWrongPacketType
	}
	me.size = bytesToInt(header[1:])

	// packet length limitation
	if me.size > MaxPacketSize {
		return ErrPacketSizeExcced
	}
	return nil
}

// Decode packet data length byte to int(Big end)
func bytesToInt(b []byte) int {
	result := 0
	for _, v := range b {
		result = result<<8 + int(v)
	}
	return result
	//return int(binary.BigEndian.Uint32(b))
}

// Decode 把网络数据流解码到 Packet.Packet(s)
// TODO(Warning): shared slice
func (me *Decoder) Decode(data []byte) ([]*Packet.Packet, error) {
	me.buf.Write(data)
	var (
		packets []*Packet.Packet
		err     error
	)
	// check length
	if me.buf.Len() < HeadLength {
		return nil, err
	}
	// first time
	if me.size < 0 {
		if err = me.forward(); err != nil {
			return nil, err
		}
	}

	for me.size <= me.buf.Len() {
		p := &Packet.Packet{
			Type:   Packet.Type(me.typ),
			Length: me.size,
			Data:   me.buf.Next(me.size),
		}
		packets = append(packets, p)

		// more packet
		if me.buf.Len() < HeadLength {
			me.size = -1
			break
		}

		if err = me.forward(); err != nil {
			return nil, err
		}
	}
	return packets, nil
}

// Encode create a packet.Packet from  the raw bytes slice and then encode to network bytes slice
// Protocol refs: https://github.com/NetEase/pomelo/wiki/Communication-Protocol
//
// -<type>-|--------<length>--------|-<data>-
// --------|------------------------|--------
// 1 byte packet type, 3 bytes packet data length(big end), and data segment
func Encode(typ Packet.Type, data []byte) ([]byte, error) {
	if typ < Packet.Handshake || typ > Packet.Kick {
		return nil, Packet.ErrWrongPacketType
	}

	p := &Packet.Packet{Type: typ, Length: len(data)}
	buf := make([]byte, p.Length+HeadLength)
	buf[0] = byte(p.Type)

	copy(buf[1:HeadLength], intToBytes(p.Length))
	copy(buf[HeadLength:], data)

	return buf, nil
}

// Encode packet data length to bytes(Big end)
func intToBytes(n int) []byte {
	buf := make([]byte, 3)
	buf[0] = byte((n >> 16) & 0xFF)
	buf[1] = byte((n >> 8) & 0xFF)
	buf[2] = byte(n & 0xFF)
	return buf
}
