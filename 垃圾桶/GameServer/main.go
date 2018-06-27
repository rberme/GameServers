package main

import (
	"flag"
	"fmt"
	"global"
	"os"
	"strconv"
	"tools"
	"tools/cfg"
)

var (
	gameIP   string
	gamePort string
)

func main() {
	fmt.Println("GameServer")
}

func getPort() {
	var s int
	flag.IntVar(&s, "s", 0, "tcp listen port")
	flag.Parse()
	if s == 0 {
		tools.ERR("please set gameserver port")
		os.Exit(-1)
	}
	gameIP = cfg.GetValue("game_ip_" + strconv.Itoa(s))
	gamePort = cfg.GetValue("game_port_" + strconv.Itoa(s))
	global.ServerName = "GameServer[" + gamePort + "]"
	global.ServerID = uint32(s)
}

func checkError(err error) {
	if err != nil {
		tools.ERR("Fatal error: %v", err)
		os.Exit(-1)
	}
}
