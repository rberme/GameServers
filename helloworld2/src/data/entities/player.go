package entities

// Player 数据实体
type Player struct {
	ID       uint
	Username string `sql:"size(15),notnull"`
	Password string `sql:"size(15),notnull"`
	// Rank     int

	// Coins  int32
	// Clouds int32

	Token string
}
