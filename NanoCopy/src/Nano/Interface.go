package Nano

import (
	"Nano/Component"
	"time"
)

// Listen 监听
func Listen(addr string, opts ...Option) {
	listen(addr, opts...)
}

// Register 注册一个组件
func Register(c Component.IComponent, options ...Component.Option) {
	comps = append(comps, regComp{c, options})
}

// SetHeartbeatInterval 设置心跳时间间隔
func SetHeartbeatInterval(d time.Duration) {
	env.heartbeat = d
}
