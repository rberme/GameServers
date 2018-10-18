package logic

import (
	"encoding/binary"
	"llog"
	"server"
)

const (
	RESULT_OK     = 0
	RESULT_FAILED = 1
	RESULT_DATA   = 2
)

//协议
const (
	MSG_CODE_HEARTBEAT     = 10095 //心跳
	MSG_CODE_HEARTBEAT_RET = 10096 //心跳返回

	MSG_CODE_CHAT_LOGIN     = 10101
	MSG_CODE_CHAT_LOGIN_RET = 10102
)

// func ProcessCmd(s *Socket) {
// }

func parseMsg(buff []byte) (msgcode int, state int8, data []byte) {
	msgcode = int(binary.LittleEndian.Uint32(buff[0:4]))
	state = int8(buff[4])
	data = buff[9:]
	datalen := int(binary.LittleEndian.Uint32(buff[5:9]))
	if datalen != len(data) {
		llog.Error("接收到错误数据 MsgCode:", msgcode)
	}
	return
}

//var testuid int64

// ProcessCmd .
func ProcessCmd(s *server.Socket, buff []byte) int {

	result := RESULT_OK
	uid := s.ID
	msgcode, _, buff := parseMsg(buff)
	//var client *GameClient
	if uid <= 0 {
		if msgcode != MSG_CODE_CHAT_LOGIN {
			ClientMgr.ClearClient(s.ID)
			server.TCPManager.CloseSocket(s)
			result = RESULT_FAILED
			goto CMDEND
		}
	} else {
		//client, _ = ClientMgr.GetClient(uid)
	}

	switch msgcode {

	// 玩家登录
	case MSG_CODE_CHAT_LOGIN:
		// _, token := cSharpBufToStr(buff)
		// //atomic.AddInt64(&testuid, 1)
		// uid := accountServerValidate(token) //testuid //
		// if uid == 0 {
		// 	result = RESULT_FAILED
		// 	goto CMDEND
		// }
		// s.ID = uid
		// atomic.StoreInt32(&s.Used, 1)
		// client, oldSocket := ClientMgr.SaveClient(uid, s, func() *GameClient {
		// 	return &GameClient{
		// 		userID:       uid,
		// 		accessTime:   time.Now(),
		// 		ClientSocket: s,
		// 		//State:        1,
		// 		userData: &model.ChatModel{},
		// 	}
		// })
		// if oldSocket != nil {
		// 	server.TCPManager.CloseSocket(oldSocket)
		// } else {
		// 	GameMq.publish(utils.MergeBytes(int32(MQ_GAMESERVER_PLAYERDATA), int32(uid)))
		// }
		// wc := chat.PutIntoWorldChannel(uid, 0)
		// client.WriteMsg(MSG_CODE_CHAT_LOGIN_RET, wc)
		// client.writeData(func() {
		// 	client.userData.WorldChannel = wc
		// })

	case MSG_CODE_HEARTBEAT:

	//无
	default:
	}
CMDEND:
	return result
}
