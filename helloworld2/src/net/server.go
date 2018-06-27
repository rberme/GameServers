package net

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// var origins = map[string]bool {
// 	"http://192.168.0.189:8338": true, // development
// 	//"http://playreapergame.com": true, // production

//   }

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// if !origins[r.Header.Get("Origin")] {
	// 	http.Error(w, "Origin not allowed", 403)
	// 	return
	// }

	// Handshake
	//ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)

	ws, err := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}).Upgrade(w, r, nil)

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
	connection := NewConnection(ws)

	authMessage := &netmsg.AuthMessage{}
	// player, err := helpers.AuthHelper.AuthenticateUsingCredentials(playerToken)
	// if err != nil { 
	// 	authMessage.Status = "bad_credentials"
	// } else {
		authMessage.Status = "success"
		authMessage.Keko = player.GetUsername()
		authMessage.Id = player.GetPlayerId()
		// authMessage.Creditos = player.GetCoins()
		// authMessage.Fichas = player.GetClouds()

		connection.AssignToPlayer(player)
		//connection.output <- authMessage.WritePacket()

		return
	//}

	connection.output <- authMessage.WritePacket()
	connection.Close()
}
