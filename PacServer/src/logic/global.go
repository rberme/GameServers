package logic

import (
	"errors"
	"fmt"
	"llog"
	"runtime"
	"server"
	"time"
)

var errValueAssert = errors.New("Processor: 类型断言错误")
var errValueLogin = errors.New("Processor: 登录错误")

var (
	// //rabbitmq
	// GameMq gameMQ
	// // Filter 屏蔽字
	// Filter = sensitive.New()
	// ClientMgr .
	ClientMgr = NewClientManager()
)

// const (
// 	/// <summary>
// 	/// 系统(统合)
// 	/// </summary>
// 	chatType_System = 1

// 	/// <summary>
// 	/// 世界
// 	/// </summary>
// 	chatType_World = 2

// 	/// <summary>
// 	/// 公会
// 	/// </summary>
// 	chatType_Organize = 4

// 	/// <summary>
// 	/// 队伍
// 	/// </summary>
// 	chatType_Team = 8

// 	/// <summary>
// 	/// 私聊
// 	/// </summary>
// 	chatType_Whisper = 16

// 	/// <summary>
// 	/// 所有频道
// 	/// </summary>
// 	chatType_All = 31
// )

// const (
// 	MSG_CODE_CHAT = iota + 1901
// 	MSG_CODE_CHAT_RET
// 	MSG_CODE_CHAT_SETWORLDCHANNEL
// 	MSG_CODE_CHAT_SETWORLDCHANNEL_RET
// )

// // // ForceCloseSocket 强制关闭连接
// // func ForceCloseSocket(s *server.Socket, clear bool) {
// // 	if clear {
// // 		ClientMgr.ClearClient(s.ID)
// // 	}
// // 	server.TCPManager.CloseSocket(s)

// // }

// // ProcessClientHeart 处理角色的心跳时间, 如果超时，则执行清除工作
// func ProcessClientHeart(client *GameClient) {
// 	//long nowTicks = DateTime.Now.Ticks / 10000;
// 	//nowTicks := time.Now().Unix()

// }

// func sendMessage(ag *GameClient, chatdata *model.GChatData, toID int32) {
// 	switch chatdata.ChatType {
// 	case chatType_All: //所有频道
// 		fallthrough
// 	case chatType_System: //系统频道
// 		result, _ := protobuf.Encode(chatdata)
// 		ClientMgr.Range(func(id int64, ag *GameClient) bool {
// 			ag.WriteMsg(MSG_CODE_CHAT_RET, int32(len(result)), result)
// 			return true
// 		})
// 	case chatType_World: //世界频道
// 		var worldID int32
// 		ag.readData(func() {
// 			worldID = ag.userData.WorldChannel
// 		})
// 		chat.World.RLock()
// 		temp, ok := chat.World.Data[int32(worldID)]
// 		if ok {
// 			result, _ := protobuf.Encode(chatdata)
// 			temp.Foreach(func(id int64) {
// 				client, ok := ClientMgr.GetClient(id)
// 				if ok {
// 					client.WriteMsg(MSG_CODE_CHAT_RET, int32(len(result)), result)
// 				}
// 			})
// 		}
// 		chat.World.RUnlock()
// 	case chatType_Team: //队伍频道
// 		var teamID int32
// 		ag.readData(func() {
// 			teamID = ag.userData.TeamID
// 		})
// 		if teamID <= 0 {
// 			break
// 		}
// 		chat.Team.RLock()
// 		temp, ok := chat.Team.Data[int32(teamID)]
// 		if ok {
// 			result, _ := protobuf.Encode(chatdata)
// 			temp.Foreach(func(id int64) {
// 				client, ok := ClientMgr.GetClient(id)
// 				if ok {
// 					client.WriteMsg(MSG_CODE_CHAT_RET, int32(len(result)), result)
// 				}
// 			})
// 		}
// 		chat.Team.RUnlock()
// 	case chatType_Organize: //社团
// 		var orgaID int32
// 		ag.readData(func() {
// 			orgaID = ag.userData.OrganizeID
// 		})
// 		if orgaID <= 0 {
// 			break
// 		}
// 		chat.Organize.RLock()
// 		temp, ok := chat.Organize.Data[int32(orgaID)]
// 		if ok {
// 			result, _ := protobuf.Encode(chatdata)
// 			temp.Foreach(func(id int64) {
// 				client, ok := ClientMgr.GetClient(id)
// 				if ok {
// 					client.WriteMsg(MSG_CODE_CHAT_RET, int32(len(result)), result)
// 				}
// 			})
// 		}
// 		chat.Organize.RUnlock()
// 	default:
// 		result, _ := protobuf.Encode(chatdata)
// 		ClientMgr.Range(func(id int64, ag *GameClient) bool {
// 			ag.WriteMsg(MSG_CODE_CHAT_RET, int32(len(result)), result)
// 			return true
// 		})
// 	}
// }

// func sendNormalMsg(ag *GameClient, chatType int32, content string, toID int32) {

// 	chatdata := &model.GChatData{}
// 	chatdata.ChatType = chatType
// 	chatdata.FromUserID = ag.userID
// 	playerData := model.BPlayerData{}
// 	ag.readData(func() {
// 		playerData = ag.userData.PlayerData
// 	})
// 	chatdata.SendDate = time.Now().Unix()
// 	chatdata.FromNick = playerData.Pname
// 	chatdata.FromExp = playerData.Pexp
// 	chatdata.Icon = playerData.Picon
// 	chatdata.Border = playerData.Piconborder

// 	chatdata.Content = make([][]byte, 1)
// 	cont := &model.GChatContent{
// 		Content: make(map[string]string),
// 	}
// 	cont.Content["str"] = content
// 	chatdata.Content[0], _ = protobuf.Encode(cont)

// 	sendMessage(ag, chatdata, toID)
// }

// func sendSystemMsg(id int32, params []string, operParams []byte) {
// 	notifTableHash := tabledata.NotificationTable.Hash
// 	chatTable, ok := notifTableHash[id]
// 	if ok == false {
// 		return
// 	}
// 	tstr := chatTable.Content
// 	chatData := &model.GChatData{
// 		Content:    createCont(tstr, params),
// 		SendDate:   time.Now().Unix(),
// 		ChatType:   chatType_System,
// 		FromExp:    chatTable.FromUserID, //int32(chatTable.FromExp),
// 		FromUserID: int64(chatTable.FromUserID),
// 		OperID:     chatTable.OperID,
// 		OperType:   chatTable.OperType,
// 		OperParams: operParams,
// 	}
// 	sendMessage(nil, chatData, 0)
// }

// // createCont 处理系统消息信息
// func createCont(tstr string, pstr []string) [][]byte {
// 	pcount := 0
// 	var cont [][]byte
// 	tstrArr := strings.Split(tstr, "|")
// 	l := len(tstrArr)
// 	for i := 0; i < l; i++ {
// 		thash := &model.GChatContent{
// 			Content: make(map[string]string),
// 		}
// 		if tstrArr[i][0] == '#' {
// 			ttstrArr := strings.Split(tstrArr[i], ",")
// 			if len(ttstrArr) == 2 {
// 				switch ttstrArr[0] {
// 				case "#str":
// 					thash.Content["str"] = pstr[pcount]
// 					pcount++
// 					thash.Content["color"] = ttstrArr[1]
// 				default:
// 					thash.Content["str"] = ttstrArr[0][1:]
// 					thash.Content["color"] = ttstrArr[1]
// 				}
// 			}
// 		} else {
// 			thash.Content["str"] = tstrArr[i]
// 		}
// 		temp, _ := protobuf.Encode(thash)
// 		cont = append(cont, temp)
// 	}
// 	return cont
// }

// ////////////////////////////////////////////////////////////////////////////////////////////
// ////////////////////////////////////////////////////////////////////////////////////////////

// ////////////////////////////////////////////////////////////////////////////////////////////
// ////////////////////////////////////////////////////////////////////////////////////////////

// func cSharpBufToStr(buff []byte) (length int, str string) {
// 	// 只处理65536长度以内的
// 	length = 0
// 	if buff[0] < 128 {
// 		length = int(buff[0])
// 		str = string(buff[1:])
// 	} else {
// 		length = int(buff[1]*128 + buff[0] - 128)
// 		str = string(buff[2:])
// 	}

// 	return
// }

// MainLoop .
func MainLoop() {
	defer func() {
		llog.Release("123")
	}()

	var idx uint64
	for {
		time.Sleep(10 * time.Second)
		idx++
		//////////////////////////////////////////////

		// for _, v := range agents.agents {
		// 	v.conn.WriteMsg(MSG_CODE_HEARTBEAT_RET, time.Now().Unix())
		// }

		agcount := ClientMgr.ClearExpired()
		llog.Release("Client数量 : %d", agcount)
		llog.Release("Socket数量 : %d", server.TCPManager.ClearUselessSocket())
		//////////////////////////////////////////////
		//if idx%6 == 0 {
		// tempAgents.Lock()
		// //llog.Debug("2..")
		// for k, v := range tempAgents.agents {
		// 	if v > 0 {
		// 		k.conn.Close()
		// 		delete(tempAgents.agents, k)
		// 	} else {
		// 		tempAgents.agents[k]++
		// 	}
		// }
		// llog.Debug("临时连接数 : %d", len(tempAgents.agents))
		// tempAgents.Unlock()
		fmt.Println("当前总Go程数:", runtime.NumGoroutine())
		//}
	}
}
