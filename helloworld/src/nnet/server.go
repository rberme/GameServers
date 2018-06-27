package nnet

import (
	"encoding/binary"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

// type Server struct {
// }

// var origins = map[string]bool {
// 	"http://192.168.0.189:8338": true, // development
// 	//"http://playreapergame.com": true, // production

//   }

type iMsgHandler interface {
	Convert(int) string
}

var msgHandler iMsgHandler

// // SetMsgHandler .
// func SetMsgHandler(msghandler iMsgHandler) {
// 	msgHandler = msghandler
// }

type tempSessMap struct {
	sync.RWMutex
	sessions map[*Session]bool
}

func (me *tempSessMap) putSession(s *Session) {
	me.Lock()
	defer me.Unlock()
	me.sessions[s] = true
}

func (me *tempSessMap) delSession(s *Session) {
	me.Lock()
	defer me.Unlock()
	delete(me.sessions, s)
}

var tempMap = &tempSessMap{
	sessions: make(map[*Session]bool),
}

// Start 启动服务器
func Start(IPPort string, msghandler iMsgHandler) {
	//http.HandleFunc("/", serveHome)
	var clientID int32
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&clientID, 1)
		wsHandler(uint64(clientID), w, r)
	})
	
	err := http.ListenAndServe(IPPort, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func wsHandler(clientID uint64, w http.ResponseWriter, r *http.Request) {
	// if !origins[r.Header.Get("Origin")] {
	// 	http.Error(w, "Origin not allowed", 403)
	// 	return
	// }

	// Handshake
	//ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)

	ws, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}).Upgrade(w, r, nil)

	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		return
	}

	// Get possible session params
	// playerId := r.URL.Query().Get("kakid")
	// playerToken := r.URL.Query().Get("kakpa")

	// if len(playerId) == 0 || len(playerToken) == 0 {
	// 	ws.Close()
	// 	return
	// }

	// Add connection
	session := newSession(ws)
	//...
	tempMap.putSession(session)
	//Hub.putSession(session)

	// authMessage := &netmsg.AuthMessage{}
	// // player, err := helpers.AuthHelper.AuthenticateUsingCredentials(playerToken)
	// // if err != nil {
	// // 	authMessage.Status = "bad_credentials"
	// // } else {
	// authMessage.Status = "success"
	// authMessage.Keko = player.GetUsername()
	// authMessage.Id = player.GetPlayerId()
	// // authMessage.Creditos = player.GetCoins()
	// // authMessage.Fichas = player.GetClouds()

	// connection.AssignToPlayer(player)
	// //connection.output <- authMessage.WritePacket()

	// return
	// //}

	// connection.output <- authMessage.WritePacket()
	// connection.Close()
}

// AuthSuccess 验证成功
func AuthSuccess(sess *Session) {
	tempMap.delSession(sess)
	Hub.putSession(sess)
}

func processPacket(id uint64, msg []byte) {

	msgID := int(binary.LittleEndian.Uint32(msg[:4]))

	if msgHandler == nil {
		fmt.Println("没有消息处理对象")
		return
	}

	funName := msgHandler.Convert(msgID)
	([]byte(funName))[0] -= 32
	fun := reflect.ValueOf(msgHandler).MethodByName(funName)
	if fun.IsNil() {
		fmt.Println("没有处理的消息")
		return
	}

	args := []reflect.Value{reflect.ValueOf(id), reflect.ValueOf(msg[4:])}
	fun.Call(args)

}
