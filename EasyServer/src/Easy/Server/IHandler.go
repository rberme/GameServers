package Server

// IHandler 消息处理
type IHandler interface {
	MainHandler(userID int, data []byte) (bool, string, string)
	ClientHandler(data []byte) (bool, string, string, int)
	UserLogin(userID int)
	UserLogout(userID int)
}
