package Control

import (
	"Chat/Model"
)

// PlayerCtrl 对玩家数据的逻辑处理
var PlayerCtrl PlayerControl

func init() {
	PlayerCtrl = PlayerControl{
		PlayerInfos: make(map[int]*Model.PlayerData),
	}
}

// PlayerControl 管理玩家聊天信息
type PlayerControl struct {
	PlayerInfos map[int]*Model.PlayerData
}
