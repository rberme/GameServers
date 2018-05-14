package Message

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// Type 消息的类型 Request/Notify/Response/Push
type Type byte

// 消息类型
const (
	Request  Type = 0x00
	Notify   Type = 0x01
	Response Type = 0x02
	Push     Type = 0x03
)

const (
	msgRouteCompressMask = 0x01
	msgTypeMask          = 0x07
	msgRouteLengthMask   = 0xFF
	msgHeadLength        = 0x02
)

var types = map[Type]string{
	Request:  "Request",
	Notify:   "Notify",
	Response: "Response",
	Push:     "Push",
}

var (
	routes = make(map[string]uint16) // route map to code
	codes  = make(map[uint16]string) // code map to route
)

// Errors 消息编码中可能会出现的错误
var (
	ErrWrongMessageType  = errors.New("wrong message type")
	ErrInvalidMessage    = errors.New("invalid message")
	ErrRouteInfoNotFound = errors.New("route info not found in dictionary")
)

// Message 表示一个未被序列化,将要被序列化的消息
type Message struct {
	Type       Type   //message type
	ID         uint   //unique id, 0为notify
	Route      string //route for locating service
	Data       []byte //payload
	compressed bool   //是否要压缩
}

// New 返回一个新的消息实例
func New() *Message {
	return &Message{}
}

func (me *Message) String() string {
	return fmt.Sprintf("Type: %s, ID: %d, Route: %s, Compressed: %t, BodyLength: %d",
		types[me.Type],
		me.ID,
		me.Route,
		me.compressed,
		len(me.Data))
}

// Encode marshals message to binary format.
func (me *Message) Encode() ([]byte, error) {
	return Encode(me)
}

func routable(t Type) bool {
	return t == Request || t == Notify || t == Push
}

func invalidType(t Type) bool {
	return t < Request || t > Push
}

// Encode marshals message to binary format. Different message types is corresponding to
// different message header, message types is identified by 2-4 bit of flag field. The
// relationship between message types and message header is presented as follows:
// ------------------------------------------
// |   type   |  flag  |       other        |
// |----------|--------|--------------------|
// | request  |----000-|<message id>|<route>|
// | notify   |----001-|<route>             |
// | response |----010-|<message id>        |
// | push     |----011-|<route>             |
// ------------------------------------------
// The figure above indicates that the bit does not affect the type of message.
// See ref: https://github.com/lonnng/nano/blob/master/docs/communication_protocol.md
func Encode(m *Message) ([]byte, error) {
	if invalidType(m.Type) {
		return nil, ErrWrongMessageType
	}

	buf := make([]byte, 0)
	flag := byte(m.Type) << 1

	code, compressed := routes[m.Route]
	if compressed {
		flag |= msgRouteCompressMask
	}
	buf = append(buf, flag)

	if m.Type == Request || m.Type == Response {
		n := m.ID
		// variant length encode
		for {
			b := byte(n % 128)
			n >>= 7
			if n != 0 {
				buf = append(buf, b+128)
			} else {
				buf = append(buf, b)
				break
			}
		}
	}

	if routable(m.Type) {
		if compressed {
			buf = append(buf, byte((code>>8)&0xFF))
			buf = append(buf, byte(code&0xFF))
		} else {
			buf = append(buf, byte(len(m.Route)))
			buf = append(buf, []byte(m.Route)...)
		}
	}

	buf = append(buf, m.Data...)
	return buf, nil
}

// Decode unmarshal the bytes slice to a message
// See ref: https://github.com/lonnng/nano/blob/master/docs/communication_protocol.md
func Decode(data []byte) (*Message, error) {
	if len(data) < msgHeadLength {
		return nil, ErrInvalidMessage
	}
	m := New()
	flag := data[0]
	offset := 1
	m.Type = Type((flag >> 1) & msgTypeMask)

	if invalidType(m.Type) {
		return nil, ErrWrongMessageType
	}

	if m.Type == Request || m.Type == Response {
		id := uint(0)
		// little end byte order
		// WARNING: must can be stored in 64 bits integer
		// variant length encode
		for i := offset; i < len(data); i++ {
			b := data[i]
			id += uint(b&0x7F) << uint(7*(i-offset))
			if b < 128 {
				offset = i + 1
				break
			}
		}
		m.ID = id
	}

	if routable(m.Type) {
		if flag&msgRouteCompressMask == 1 {
			m.compressed = true
			code := binary.BigEndian.Uint16(data[offset:(offset + 2)])
			route, ok := codes[code]
			if !ok {
				return nil, ErrRouteInfoNotFound
			}
			m.Route = route
			offset += 2
		} else {
			m.compressed = false
			rl := data[offset]
			offset++
			m.Route = string(data[offset:(offset + int(rl))])
			offset += int(rl)
		}
	}

	m.Data = data[offset:]
	return m, nil
}
