package main

import (
	"fmt"
	"llog"
	"logic"
	"logic/config"
	"runtime"
	"server"
)

type temp struct {
}

func main() {
	fmt.Println("hello world.")
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
	fmt.Printf("     %s\n\n", config.NetCfg.IPPort)

	// tabledata.ReadAll()
	// llog.Release("读取XML表.")

	// logic.GameMq = *logic.NewGameMQ()
	// err := logic.GameMq.Start()
	// if err != nil {
	// 	llog.Error("rabbitmq 连接失败")
	// 	return
	// }
	// llog.Release("连接rabbitmq.")

	// logic.Filter.LoadWordDict("./txt/dict.txt")
	// llog.Release("加载屏蔽字库.")

	go logic.MainLoop()
	server.ProcessCmd = logic.ProcessCmd
	server.ReceiveEnd = logic.ClientMgr.ClearClient
	llog.Release("服务器开始.")
	server.TCPManager.Start(config.NetCfg.IPPort)

	// test()
	// fmt.Println(time.Now().Unix())
}

var arr [2]int

func test() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err) //
		}
	}()
	idx := 2
	arr[idx] = 1
}
