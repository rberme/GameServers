package Server

const ( //私有 共享通道 通信指令
	//commandStart = iota + 1000
	// PrivateSend 私有通道中的发送数据
	heartBeat    = 96
	heartBeatRet = 97

	recvMsg = iota + 100
	c2MMsg
	m2CMsg
	c2cMsg
	//主逻辑循环发送消息给客户端通过这个
	m2CSend
	//DataReq 获取数据
	//DataReq
)

// ChanMsg 发送个sess chan的指令
type chanMsg struct {
	msgCode int
	userID  int
	data    []byte
}

// // DataReq 读取私有数据请求
// type DataReq struct {
// 	TypeName string
// 	UserID   int
// 	Request  chan DataResp
// }

// DataResp 返回私有数据请求
type DataResp struct {
	TypeName string
	UserID   int
	DataKey  int
	Data     []byte
}
