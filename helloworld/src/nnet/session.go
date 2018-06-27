package nnet

import (
	"github.com/gorilla/websocket"
)

// Session .
type Session struct {
	ws *websocket.Conn
	ID uint64
	// Network channels
	//broadcast chan interface{}
	//output    chan interface{}
	//input chan interface{}
}

// NewSession .
func newSession(_ws *websocket.Conn) *Session {
	// The pointer allow us to modify session struct from outside
	session := &Session{
		ws: _ws,
		//broadcast: make(chan interface{}),
		//output:    make(chan interface{}),
		//input: make(chan interface{}),
	}

	//go session.Writer()
	go session.reader()

	return session
}

// WriteJSON .
func (me *Session) writeJSON(v interface{}) {
	err := me.ws.WriteJSON(v)
	if err != nil {
		Hub.delSession(me)
	}
}

// Writer .
func (me *Session) Write(buf []byte) {
	err := me.ws.WriteMessage(websocket.BinaryMessage, buf)
	if err != nil {
		Hub.delSession(me)
	}
}

// Reader .
func (me *Session) reader() {
	defer func() {
		Hub.delSession(me)
	}()

	for {
		_, message, err := me.ws.ReadMessage()
		if err != nil {
			break
		}

		// obj, err := simplejson.NewJson([]byte(message))
		// if err != nil {
		// 	fmt.Println("Can not parse the message:", string(message))
		// 	continue
		// }
		processPacket(me.ID, message)

	}
}

// Close .
func (me *Session) close() {
	// Close channels
	//close(c.output)

	// Close the websocket
	me.ws.Close()

	//c.player = nil
}
