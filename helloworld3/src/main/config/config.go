package config

import (
	"llog"
	"log"

	"github.com/go-ini/ini"
)

var filepath = "config.ini"

// Config 配置信息
type Config struct {
	ServerID         int32  `ini:"ServerId"`
	IPPort           string `ini:"IPPort"`
	AccountServer    string `ini:"AccountServer"`
	RabbitMQIPPort   string `ini:"RabbitMQIPPort"`
	RabbitMQvHost    string `ini:"RabbitMQvHost"`
	RabbitMQUser     string `ini:"RabbitMQUser"`
	RabbitMQPassword string `ini:"RabbitMQPassword"`
}

// Cfg 配置信息
var Cfg Config

func init() {
	conf, err := ini.Load(filepath)
	if err != nil {
		llog.Error("配置文件加载失败")
	}
	conf.BlockMode = false
	err = conf.MapTo(&Cfg)
	if err != nil {
		log.Println("配置信息映射失败")
	}
}
