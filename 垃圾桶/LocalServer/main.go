package main

import (
	"fmt"
	"global"
	"os"
	"tools"
	"tools/cfg"
)

var (
	localIP   string
	localPort string
)

func main() {
	fmt.Println("LocalServer")
}

func getPort() {
	//端口号
	localIP = cfg.GetValue("local_ip")
	localPort = cfg.GetValue("local_port")
	global.ServerName = "LocalServer[" + localPort + "]"
}

func startLocalServer() {
	defer func() {
		if x := recover(); x != nil {
			tools.ERR("caught panic in main()", x)
		}
	}()

	// 获取端口号
	getPort()

	//开启LocalServer监听
	startLocalServer()
	// //开启客户端监听
	// getPort()
	// json := codec.JSON()
	// // json.Register(AddReq{})
	// // json.Register(AddRsp{})
	// addr := "0.0.0.0:" + localPort
	// listener, err := net.Listen("tcp", addr)
	// if err != nil {
	// 	return
	// }
	// server := global.NewServer(listener, json, 0, global.HandlerFunc(func(session *global.Session) {

	// }))
	// server.Serve()

}

func startLocalServer() {
	// //连接Redis
	// redisProxyErr := redisProxy.InitClient(cfg.GetValue("redis_ip"), cfg.GetValue("redis_port"), cfg.GetValue("redis_pwd"))
	// checkError(redisProxyErr)

	// //开启DB
	// db.Init()

	// //开启同步DB数据到数据库
	// dbProxy.StartSysDB()

	//开启客户端监听
	addr := "0.0.0.0:" + localPort
	err := global.Listener("tcp", addr, global.PackCodecType_UnSafe,
		func(session *Session) {
			global.AddSession(session)
		},
		gameProxy.MsgDispatch,
	)
	checkError(err)
}

func stopLocalServer() {
	tools.INFO("Waiting SyncDB...")
	//dbProxy.SyncDB()
	tools.INFO("SyncDB Success")
}

func checkError(err error) {
	if err != nil {
		tools.ERR("Fatal error: %v", err)
		os.Exit(-1)
	}
}
