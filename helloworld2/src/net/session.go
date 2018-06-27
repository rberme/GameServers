package net

import (
	"fmt"
	"interfaces"

	"github.com/gorilla/websocket"
)

// Session .
type Session struct {
	ws *websocket.Conn
	ID int64
	// Network channels
	broadcast chan interface{}
	output    chan interface{}
	input     chan interface{}

	player interfaces.IPlayer
}

// type hub struct {
// 	sessions map[*Session]bool
// 	rooms    map[string]*netmsg.RoomMessage
// }

// func newHub() *hub {
// 	return &hub{
// 		sessions: make(map[*Session]bool),
// 		rooms:    make(map[string]*netmsg.RoomMessage),
// 	}
// }

// // Hub 用户中心
// var Hub = newHub()

// NewSession .
func NewSession(_ws *websocket.Conn) *Session {
	// The pointer allow us to modify session struct from outside
	session := &Session{
		ws:        _ws,
		broadcast: make(chan interface{}),
		output:    make(chan interface{}),
		input:     make(chan interface{}),
	}

	go session.Writer()
	go session.Reader()

	return session
}

// AssignToPlayer .
func (c *Session) AssignToPlayer(_player interfaces.IPlayer) {
	if _player == nil {
		panic("Session - Player interface can not be nil")
	}

	//Hub.sessions[c] = true

	c.player = _player
}

// Writer .
func (c *Session) Writer() {
	for {
		select {
		case netmessage := <-c.output:
			if netmessage == nil {
				fmt.Println("Session", "Writer", "Netmessage == nil, breaking loop")
				break
			}
			c.ws.WriteJSON(netmessage)
		case netmessage := <-c.broadcast:
			if netmessage == nil {
				fmt.Println("Session", "Writer", "Netmessage == nil, breaking loop")
				break
			}

			for hc := range Hub.sessions {
				hc.ws.WriteJSON(netmessage)
			}
		}
	}
}

// Reader .
func (c *Session) Reader() {
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}

		obj, err := simplejson.NewJson([]byte(message))
		if err != nil {
			fmt.Println("Can not parse the message:", string(message))
			continue
		}

		c.processPacket(obj)
	}
}

// Close .
func (c *Session) Close() {
	// Close channels
	close(c.output)

	// Close the websocket
	c.ws.Close()

	c.player = nil
}
