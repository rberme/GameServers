package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"llog"
	"main/config"
	"main/model"
	"main/rmq"
	"serialize/protobuf"
	"strconv"
	"strings"
)

type gameMQ struct {
	exchange  string
	recvQueue string
	sendQueue string
}

func newGameMQ() *gameMQ {
	svrID := strconv.Itoa(int(config.Cfg.ServerID))
	return &gameMQ{
		exchange:  "GameChat",
		recvQueue: "ChatServer" + svrID,
		sendQueue: "GameServer" + svrID,
	}
}

func (me *gameMQ) start() error {
	mqAddr := fmt.Sprintf("amqp://%s:%s@%s/%s",
		config.Cfg.RabbitMQUser,
		config.Cfg.RabbitMQPassword,
		config.Cfg.RabbitMQIPPort,
		config.Cfg.RabbitMQvHost,
	) //"amqp://manager:0@192.168.0.91:5672/Game")
	err := rmq.SetupRMQ(mqAddr)
	if err != nil {
		return err
	}
	rmq.BindQueue(me.exchange, me.recvQueue)
	me.receive()
	return nil
}

func (me *gameMQ) publish(buff []byte) error {
	return rmq.PublishMQ(me.exchange, me.sendQueue, buff)
}

func (me *gameMQ) receive() {
	rmq.ReceiveMQ(me.exchange, me.recvQueue, rmqHandle)
}

////////////////////////////////////////////////////////////////////////////////////////////
//mq协议
const (
	// 从游戏服务器获取玩家基本数据
	MQ_GAMESERVER_PLAYERDATA     = 88001
	MQ_GAMESERVER_PLAYERDATA_RET = 88002

	MQ_SYSTEMMSG = 88003
)

func rmqHandle(msg []byte) {
	//llog.Debug(string(msg))
	msgcode := int32(binary.LittleEndian.Uint32(msg))
	llog.Release("MQ消息类型:%d", msgcode)
	switch msgcode {
	case MQ_GAMESERVER_PLAYERDATA_RET:
		bPlayerData := &model.BPlayerData{}
		err := protobuf.Decode(msg[4:], bPlayerData)
		if err == nil {
			uid := bPlayerData.Pid / 100000
			ag := globalAgents.Get(uid)
			if ag != nil {
				ag.writeData(func() {
					ag.userData.PlayerData = *bPlayerData
					ag.userData.OrganizeID = bPlayerData.Porganid
				})
			}
		}
	case MQ_SYSTEMMSG:
		//str := string(msg[4:])
		//llog.Release("系统消息内容:%s", str)
		var sysmsg = SystemMsg{} //json.RawMessage
		err := json.Unmarshal(msg[4:], &sysmsg)
		if err != nil {
			llog.Error(err.Error())
			break
		}
		params := strings.Split(sysmsg.Params, "\r")
		id, _ := strconv.Atoi(sysmsg.ID)

		sendSystemMsg(int32(id), params, nil)
	default:
	}
}

// SystemMsg .
type SystemMsg struct {
	ID       string `json:"sysmsg_id"`
	OperType int32  `json:"sysmsg_oper"`
	Params   string `json:"sysmsg_params"`
	//MsgCode  int32  `json:"msgcode"`
}
