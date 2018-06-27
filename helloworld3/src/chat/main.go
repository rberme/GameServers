package main

import (
	"chat/model"
	"fmt"
	"llog"
	"network"
	"sync"
)

type chat struct {
	MaxConnNum      int
	PendingWriteNum int
	MaxMsgLen       uint32
	Processor       network.Processor
	// AgentChanRPC    *chanrpc.Server

	// tcp
	TCPAddr      string
	LenMsgLen    int
	LittleEndian bool
}

// Run 服务器启动,开始监听
func (me *chat) Run(closeSig chan bool) {

	var tcpServer *network.TCPServer
	if me.TCPAddr != "" {
		tcpServer = new(network.TCPServer)
		tcpServer.Addr = me.TCPAddr
		tcpServer.MaxConnNum = me.MaxConnNum
		tcpServer.PendingWriteNum = me.PendingWriteNum
		tcpServer.LenMsgLen = me.LenMsgLen
		tcpServer.MaxMsgLen = me.MaxMsgLen
		tcpServer.LittleEndian = me.LittleEndian
		tcpServer.NewAgent = func(conn *network.TCPConn) network.Agent {
			a := &agent{conn: conn}
			agents.Add(a)
			return a
		}
	}

	if tcpServer != nil {
		tcpServer.Start()
	}
	<-closeSig
	if tcpServer != nil {
		tcpServer.Close()
	}
}

func (me *chat) OnDestroy() {

}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type agentSet struct {
	sync.RWMutex
	agents map[int64]*agent
}

func (me *agentSet) Add(a *agent) {
	me.Lock()
	defer me.Unlock()
	me.agents[a.userID] = a
}

func (me *agentSet) Del(id int64) {
	me.Lock()
	defer me.Unlock()
	delete(me.agents, id)
}

var agents agentSet

func main() {
	// fmt.Println(accountServerValidate("55ce4822af33100491a0c7a0109e2671"))
	// return

	var cmd string

	agents = agentSet{agents: make(map[int64]*agent)}

	closeSig := make(chan bool)
	c := &chat{
		Processor:       &ChatProcessor{},
		TCPAddr:         "192.168.0.189:19001",
		MaxConnNum:      2000,
		PendingWriteNum: 100,
	}
	go c.Run(closeSig)

	for {
		fmt.Scanf("%s", &cmd)
		switch cmd {
		case "quit":
			closeSig <- true
			//c.OnDestroy()
			return
		default:
			fmt.Printf("%s", cmd)
		}
	}
	// fmt.Println("end")
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type agent struct {
	conn     network.Conn
	chat     *chat
	userData model.ChatModel
	userID   int64
}

func (me *agent) Run() {
	for {
		data, err := me.conn.ReadMsg()
		if err != nil {
			llog.Debug("read message: %v", err)
			break
		}

		if me.chat.Processor != nil {
			msg, err := me.chat.Processor.Unmarshal(data)
			if err != nil {
				llog.Debug("unmarshal message error: %v", err)
				break
			}
			err = me.chat.Processor.Route(msg, me)
			if err != nil {
				llog.Debug("route message error: %v", err)
				break
			}
		}
	}
}

func (me *agent) OnClose() {
	agents.Del(me.userID)
}
