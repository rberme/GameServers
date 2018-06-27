package main

import (
	"chat/model"
	"encoding/binary"
	"errors"
	"serialize/protobuf"
)

//协议
const (
	MSG_CODE_CHAT = iota + 1901
	MSG_CODE_CHAT_RET
	MSG_CODE_CHAT_SETWORLDCHANNEL
	MSG_CODE_CHAT_SETWORLDCHANNEL_RET
	MSG_CODE_CHAT_LOGIN
	MSG_CODE_CHAT_LOGIN_RET
)

var errValueAssert = errors.New("Processor: 类型断言错误")
var errValueLogin = errors.New("Processor: 登录错误")

var worldChannel = newChannelMap()

// ChatProcessor .
type ChatProcessor struct {
}

// Route must goroutine safe
func (me *ChatProcessor) Route(msg interface{}, a interface{}) error {
	basemsg, ok := msg.(*model.Basemsg)
	if ok == false {
		return errValueAssert
	}
	ag, ok := a.(*agent)
	if ok == false {
		return errValueAssert
	}

	switch basemsg.Msgcode {
	case MSG_CODE_CHAT_LOGIN:
		token := string(basemsg.Buff)
		uid := accountServerValidate(token)
		if uid == 0 {
			return errValueLogin
		}
		ag.userID = uid
		ag.userData.WorldChannel = PutIntoWorldChannel(uid, 0)
		//getoutWorldChannel(uid,)
		basemsg.Msgcode = MSG_CODE_CHAT_LOGIN_RET
		basemsg.Buff = make([]byte, 4)
		binary.LittleEndian.PutUint32(basemsg.Buff, uint32(ag.userData.WorldChannel))
		buff, err := protobuf.Encode(basemsg)
		if err == nil {
			ag.conn.WriteMsg(buff)
		}
	case MSG_CODE_CHAT:

	case MSG_CODE_CHAT_SETWORLDCHANNEL:
		worldID := int32(binary.LittleEndian.Uint32(basemsg.Buff))
		temp := PutIntoWorldChannel(ag.userID, worldID)
		if temp > 0 {
			getoutWorldChannel(ag.userID, ag.userData.WorldChannel)
			ag.userData.WorldChannel = worldID
		}
	default:
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////

func getoutWorldChannel(pid int64, worldID int32) {
	worldChannel.RLock()
	v, ok := worldChannel.Data[worldID]
	worldChannel.RUnlock()
	if ok {
		v.del(pid)
	}
}

// PutIntoWorldChannel 加入世界频道
func PutIntoWorldChannel(pid int64, targetID int32) int32 {
	if targetID == 0 { //系统分配
		worldChannel.Lock()
		length := len(worldChannel.Data)
		var retval int32
		for i := 1; i <= length+1; i++ {
			v, ok := worldChannel.Data[int32(i)]
			if ok {
				if v.length() < maxNumPerChannel {
					v.add(pid)
					retval = int32(i)
					break
				}
			} else {
				worldChannel.add(int32(i), pid)
				retval = int32(i)
				break
			}
		}
		worldChannel.Unlock()
		return retval
	} else if targetID > 0 && targetID < maxWorldChannel { //玩家分配
		worldChannel.RLock()
		v, ok := worldChannel.Data[targetID]
		worldChannel.RUnlock()
		if ok == true {
			v.add(pid)
		} else {
			worldChannel.Lock()
			worldChannel.add(targetID, pid)
			worldChannel.Unlock()
		}
		return targetID
	}
	return 0
}

////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////

// Unmarshal must goroutine safe
func (me *ChatProcessor) Unmarshal(data []byte) (interface{}, error) {
	msg := &model.Basemsg{}
	err := protobuf.Decode(data, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// Marshal must goroutine safe
func (me *ChatProcessor) Marshal(msg interface{}) ([][]byte, error) {
	return nil, nil
}
