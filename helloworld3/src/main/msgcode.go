package main

import (
	"encoding/binary"
	"errors"
	"llog"
	"main/channel"
	"main/model"
	"main/tabledata"
	"serialize/protobuf"
	"strings"
	"time"
	"utils"
)

//协议
const (
	MSG_CODE_HEARTBEAT_RET = 10096 //心跳返回

	MSG_CODE_CHAT_LOGIN     = 10101
	MSG_CODE_CHAT_LOGIN_RET = 10102
)
const (
	MSG_CODE_CHAT = iota + 1901
	MSG_CODE_CHAT_RET
	MSG_CODE_CHAT_SETWORLDCHANNEL
	MSG_CODE_CHAT_SETWORLDCHANNEL_RET
)

const (
	/// <summary>
	/// 系统(统合)
	/// </summary>
	chatType_System = 1

	/// <summary>
	/// 世界
	/// </summary>
	chatType_World = 2

	/// <summary>
	/// 公会
	/// </summary>
	chatType_Organize = 4

	/// <summary>
	/// 队伍
	/// </summary>
	chatType_Team = 8

	/// <summary>
	/// 私聊
	/// </summary>
	chatType_Whisper = 16

	/// <summary>
	/// 所有频道
	/// </summary>
	chatType_All = 31
)

// const (
// 	MQ_GAMESERVER_PLAYERDATA = iota + 88001
// )

var errValueAssert = errors.New("Processor: 类型断言错误")
var errValueLogin = errors.New("Processor: 登录错误")

// ChatProcessor .
type ChatProcessor struct {
}

func (me *ChatProcessor) parseMsg(buff []byte) (msgcode int, state int8, data []byte) {
	msgcode = int(binary.LittleEndian.Uint32(buff[0:4]))
	state = int8(buff[4])
	data = buff[9:]
	datalen := int(binary.LittleEndian.Uint32(buff[5:9]))
	if datalen != len(data) {
		llog.Error("接收到错误数据 MsgCode:", msgcode)
	}

	return
}

// Route must goroutine safe
func (me *ChatProcessor) Route(msg interface{}, a interface{}) error {
	b, ok := msg.([]byte)
	if ok == false {
		return errValueAssert
	}
	ag, ok := a.(*agent)
	if ok == false {
		return errValueAssert
	}

	msgcode, _, buff := me.parseMsg(b)
	if ag.userID > 0 && globalAgents.Get(ag.userID) == nil {
		return errValueLogin
	} else if ag.userID == 0 && msgcode != MSG_CODE_CHAT_LOGIN {
		return errValueLogin
	}

	llog.Release("接收消息ID: %d", msgcode)
	switch msgcode {

	// 玩家登录
	case MSG_CODE_CHAT_LOGIN:
		_, token := bufToStr(buff)
		uid := accountServerValidate(token)
		if uid == 0 {
			return errValueLogin
		}
		ag.userID = uid

		wc := channel.PutIntoWorldChannel(uid, 0)
		ag.writeData(func() {
			ag.userData.WorldChannel = wc
		})
		ag.WriteMsg(MSG_CODE_CHAT_LOGIN_RET, wc)

		tempAgents.Del(ag)
		globalAgents.Add(uid, ag)
		gameMq.publish(utils.MergeBytes(int32(MQ_GAMESERVER_PLAYERDATA), int32(uid)))

	//发送聊天内容
	case MSG_CODE_CHAT:
		chatType := int32(binary.LittleEndian.Uint32(buff))
		_, content := bufToStr(buff[4:])
		content = filter.Replace(content, 42) //过滤非法字符
		sendNormalMsg(ag, chatType, content, 0)

	//设置世界频道
	case MSG_CODE_CHAT_SETWORLDCHANNEL:
		worldID := int32(binary.LittleEndian.Uint32(buff))
		temp := channel.PutIntoWorldChannel(ag.userID, worldID)
		if temp > 0 {
			var oldworldID int32
			ag.writeData(func() {
				oldworldID = ag.userData.WorldChannel
				ag.userData.WorldChannel = worldID
			})
			channel.GetoutWorldChannel(ag.userID, oldworldID)

		}
	default:
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////

func sendMessage(ag *agent, chatdata *model.GChatData, toID int32) {
	switch chatdata.ChatType {
	case chatType_All: //所有频道
		fallthrough
	case chatType_System: //系统频道
		result, _ := protobuf.Encode(chatdata)
		globalAgents.Range(func(id int64, ag *agent) bool {
			ag.WriteMsg(MSG_CODE_CHAT_RET, int32(len(result)), result)
			return true
		})
	case chatType_World: //世界频道
		var worldID int32
		ag.readData(func() {
			worldID = ag.userData.WorldChannel
		})
		channel.World.RLock()
		temp, ok := channel.World.Data[int32(worldID)]
		if ok {
			result, _ := protobuf.Encode(chatdata)
			temp.Foreach(func(id int64) {
				globalAgents.Get(id).WriteMsg(MSG_CODE_CHAT_RET, int32(len(result)), result)
			})
		}
		channel.World.RUnlock()
	case chatType_Team: //队伍频道
		var teamID int32
		ag.readData(func() {
			teamID = ag.userData.TeamID
		})
		if teamID <= 0 {
			break
		}
		channel.Team.RLock()
		temp, ok := channel.Team.Data[int32(teamID)]
		if ok {
			result, _ := protobuf.Encode(chatdata)
			temp.Foreach(func(id int64) {
				globalAgents.Get(id).WriteMsg(MSG_CODE_CHAT_RET, int32(len(result)), result)
			})
		}
		channel.Team.RUnlock()
	case chatType_Organize: //社团
		var orgaID int32
		ag.readData(func() {
			orgaID = ag.userData.OrganizeID
		})
		if orgaID <= 0 {
			break
		}
		channel.Organize.RLock()
		temp, ok := channel.Organize.Data[int32(orgaID)]
		if ok {
			result, _ := protobuf.Encode(chatdata)
			temp.Foreach(func(id int64) {
				globalAgents.Get(id).WriteMsg(MSG_CODE_CHAT_RET, int32(len(result)), result)
			})
		}
		channel.Organize.RUnlock()
	default:
		result, _ := protobuf.Encode(chatdata)
		globalAgents.Range(func(id int64, ag *agent) bool {
			ag.WriteMsg(MSG_CODE_CHAT_RET, int32(len(result)), result)
			return true
		})
	}
}

func sendNormalMsg(ag *agent, chatType int32, content string, toID int32) {

	chatdata := &model.GChatData{}
	chatdata.ChatType = chatType
	chatdata.FromUserID = ag.userID
	playerData := model.BPlayerData{}
	ag.readData(func() {
		playerData = ag.userData.PlayerData
	})
	chatdata.SendDate = time.Now().Unix()
	chatdata.FromNick = playerData.Pname
	chatdata.FromExp = playerData.Pexp
	chatdata.Icon = playerData.Picon
	chatdata.Border = playerData.Piconborder

	chatdata.Content = make([][]byte, 1)
	cont := &model.GChatContent{
		Content: make(map[string]string),
	}
	cont.Content["str"] = content
	chatdata.Content[0], _ = protobuf.Encode(cont)

	sendMessage(ag, chatdata, toID)
}

func sendSystemMsg(id int32, params []string, operParams []byte) {
	notifTableHash := tabledata.NotificationTable.Hash
	chatTable, ok := notifTableHash[id]
	if ok == false {
		return
	}
	tstr := chatTable.Content
	chatData := &model.GChatData{
		Content:    createCont(tstr, params),
		SendDate:   time.Now().Unix(),
		ChatType:   chatType_System,
		FromUserID: int64(chatTable.FromUserID),
		OperID:     chatTable.OperID,
		OperType:   chatTable.OperType,
		OperParams: operParams,
	}
	sendMessage(nil, chatData, 0)
}

// createCont 处理系统消息信息
func createCont(tstr string, pstr []string) [][]byte {
	pcount := 0
	var cont [][]byte
	tstrArr := strings.Split(tstr, "|")
	l := len(tstrArr)
	for i := 0; i < l; i++ {
		thash := &model.GChatContent{
			Content: make(map[string]string),
		}
		if tstrArr[i][0] == '#' {
			ttstrArr := strings.Split(tstrArr[i], ",")
			switch ttstrArr[0] {
			case "#str":
				thash.Content["str"] = pstr[pcount]
				pcount++
				thash.Content["color"] = ttstrArr[1]
			default:
				thash.Content["str"] = ttstrArr[0][1:]
				thash.Content["color"] = ttstrArr[1]
			}

		} else {
			thash.Content["str"] = tstrArr[i]
		}
		temp, _ := protobuf.Encode(thash)
		cont = append(cont, temp)
	}
	return cont
}

////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////

func bufToStr(buff []byte) (length int, str string) {
	// 只处理65536长度以内的
	length = 0
	if buff[0] < 128 {
		length = int(buff[0])
		str = string(buff[1:])
	} else {
		length = int(buff[1]*128 + buff[0] - 128)
		str = string(buff[2:])
	}

	return
}

////////////////////////////////////////////////////////////////////////////////////////////

// Unmarshal must goroutine safe
func (me *ChatProcessor) Unmarshal(data []byte) (interface{}, error) {
	// msg := &model.Basemsg{}
	// err := protobuf.Decode(data, msg)
	// if err != nil {
	// 	return nil, err
	// }

	// return msg, nil

	return data, nil
}

// Marshal must goroutine safe
func (me *ChatProcessor) Marshal(msg interface{}) ([][]byte, error) {
	return nil, nil
}
