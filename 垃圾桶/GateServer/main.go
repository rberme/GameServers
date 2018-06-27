package main

import (
	"fmt"
)

var (
	gatewayIP    string
	gatewayPort  string
	transferIP   string
	transferPort string
)

func main() {

	if true {
		defer func() {
			fmt.Println("if end.")
		}()
	}
	fmt.Println("main end.")

	// // 获取端口号
	// getPort()

	// // //启动
	// // global.Startup(global.ServerName, "gateway_log", nil)

	// // //开启TransferServer，由GateServer充当中转服务器
	// // err := transferProxy.InitServer(transfer_port)
	// // checkError(err)
	// // INFO("Starting TransferServer")

	// //开启GateServer监听
	// startGateway()

	// // //保持进程
	// // global.Run()
}

// func getPort() {
// 	//端口号
// 	gatewayIP = cfg.GetValue("gateway_ip")
// 	gatewayPort = cfg.GetValue("gateway_port")
// 	global.ServerName = "GateServer[" + gatewayPort + "]"

// 	transferIP = cfg.GetValue("transfer_ip")
// 	transferPort = cfg.GetValue("transfer_port")
// }

// func startGateway() {
// 	// msgDispatch := dispatch.NewDispatch(
// 	// 	dispatch.HandleFunc{
// 	// 		H: transferProxy.SendToGameServer,
// 	// 	},
// 	// )

// 	addr := "0.0.0.0:" + gatewayPort
// 	err := global.Listen("tcp", addr,
// 		func(session *net.Conn) {
// 			//将此Session记录在缓存内，消息回传时使用
// 			global.AddSession(session)
// 			// //通知LoginServer用户上线
// 			// transferProxy.SetClientSessionOnline(session)
// 			// //添加session关闭时回调
// 			// session.AddCloseCallback(session, func() {
// 			// 	//通知LoginServer、GameServer用户下线
// 			// 	transferProxy.SetClientSessionOffline(session.Id())
// 			// })
// 		},
// 	)
// 	checkError(err)
// }

// func checkError(err error) {
// 	if err != nil {
// 		tools.ERR("Fatal error: %v", err)
// 		os.Exit(-1)
// 	}
// }
