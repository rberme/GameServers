package config

import (
	"llog"
	"log"

	"github.com/go-ini/ini"
)

var filepath = "config.ini"

// NetConfig 配置信息
type NetConfig struct {
	ServerID int32  `ini:"ServerId"`
	IPPort   string `ini:"IPPort"`
}

// NetCfg 配置信息
var NetCfg NetConfig

func init() {
	conf, err := ini.Load(filepath)
	if err != nil {
		llog.Error("配置文件加载失败")
	}
	conf.BlockMode = false
	err = conf.MapTo(&NetCfg)
	if err != nil {
		log.Println("配置信息映射失败")
	}
}
