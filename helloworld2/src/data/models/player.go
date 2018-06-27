package models

import "data/entities"

// Player 玩家数据
type Player struct {
	PlayerEntity *entities.Player

	input <-chan interface{}
}

// NewPlayer .
func NewPlayer(_entity *entities.Player) *Player {
	p := &Player{}
	p.PlayerEntity = _entity

	return p
}

// GetUsername 使用方法获取数据
func (p *Player) GetUsername() string {
	return p.PlayerEntity.Username
}

// GetPlayerID .
func (p *Player) GetPlayerID() int64 {
	return int64(p.PlayerEntity.ID)
}

// func (p *Player) GetCoins() int32 {
// 	return p.PlayerEntity.Coins
// }

// func (p *Player) GetClouds() int32 {
// 	return p.PlayerEntity.Clouds
// }
