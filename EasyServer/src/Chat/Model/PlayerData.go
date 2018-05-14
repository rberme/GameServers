package Model

import (
	"time"
)

// PlayerData 玩家基本信息
type PlayerData struct {
	UserID       int
	NickName     string
	LastChatTime time.Time
}
