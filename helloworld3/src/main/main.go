package main

import (
	"fmt"
	"llog"
	"main/config"
	"main/tabledata"
	"network"
	"runtime"
	"time"
	"utils/sensitive"
)

type chat struct {
	MaxConnNum int
	MaxMsgLen  uint32
	Processor  network.Processor

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
		tcpServer.LenMsgLen = me.LenMsgLen
		tcpServer.MaxMsgLen = me.MaxMsgLen
		tcpServer.LittleEndian = me.LittleEndian
		tcpServer.NewAgent = func(conn *network.TCPConn) network.Agent {
			a := &agent{
				conn:      conn,
				chat:      me,
				loginTime: time.Now(),
			}
			tempAgents.Add(a)
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

var filter = sensitive.New()
var globalAgents *agentSet
var tempAgents *tempSet
var gameMq gameMQ

func main() {
	//使用所有CPU
	maxProcs := runtime.NumCPU()
	runtime.GOMAXPROCS(maxProcs)

	var logo = `
	       ___     ___                      |         
	     .i .-'   '-. i.                    |         
	   .'   '/     \'  _'.                  |         
	   |,-../ o   o \.' '|                  |         
	(| |   /  _\ /_  \   | |)               |         
	 \\\  (_.'.'"'.'._)  ///                |         
	  \\'._(..:   :..)_.'//                 |          __       __    
	   \'.__\ .:-:. /__.'/                  |         / <'     '> \   
	    '-i-->.___.<--i-'                   |        (  / @   @ \  )  
	    .'.-'/.=^=.\'-.'.                   |         \(_ _\_/_ _)/   
	   /.'  //     \\  '.\                  |       (\ '-/     \-' /) 
	  ||   ||       ||   ||                 |        "===\     /==="  
	  \)   ||       ||   (/                 |         .==')___('==.   
	       \)       (/                      |        ' .='     '=.    

	     `

	fmt.Println(logo)
	fmt.Printf("     %s\n\n", config.Cfg.IPPort)

	tabledata.ReadAll()
	llog.Release("读取XML表.")

	gameMq = *newGameMQ()
	err := gameMq.start()
	if err != nil {
		llog.Error("rabbitmq 连接失败")
		return
	}
	llog.Release("连接rabbitmq.")

	filter.LoadWordDict("./dict.txt")
	llog.Release("加载屏蔽字库.")
	var cmd string
	// 42 即 "*"
	//cmd = filter.Replace("赌博机", 42)

	globalAgents = newAgentSet() //agentSet{agents: sync.Map{}}
	tempAgents = &tempSet{agents: make(map[*agent]int8)}

	closeSig := make(chan bool)
	c := &chat{
		Processor:    &ChatProcessor{},
		TCPAddr:      config.Cfg.IPPort,
		MaxConnNum:   2000,
		LenMsgLen:    4, //表示消息长度字节
		MaxMsgLen:    8192,
		LittleEndian: true,
	}

	go mainLoop()
	go c.Run(closeSig)
	llog.Debug("服务器启动.")

	for {
		fmt.Scanf("%s", &cmd)
		switch cmd {
		case "quit":
			closeSig <- true
			return
		default:
			fmt.Printf("%s", cmd)
		}
	}
	// fmt.Println("end")
}

func mainLoop() {
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
		agcount := 0
		globalAgents.Range(func(id int64, ag *agent) bool {
			ag.WriteMsg(MSG_CODE_HEARTBEAT_RET, time.Now().Unix())
			agcount++
			return true
		})
		llog.Release("连接数量 : %d", agcount)
		//////////////////////////////////////////////
		//if idx%6 == 0 {
		tempAgents.Lock()
		llog.Debug("2..")
		for k, v := range tempAgents.agents {
			if v > 0 {
				k.conn.Close()
				delete(tempAgents.agents, k)
			} else {
				tempAgents.agents[k]++
			}
		}
		llog.Debug("3.. len = %d", len(tempAgents.agents))
		tempAgents.Unlock()
		fmt.Println("当前总Go程数:", runtime.NumGoroutine())
		//}
	}
}
