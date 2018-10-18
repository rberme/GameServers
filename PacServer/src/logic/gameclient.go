package logic

import (
	"encoding/binary"
	"server"
	"sync"
	"time"
	"utils"
)

// GameClient .
type GameClient struct {
	sync.RWMutex
	accessTime   time.Time
	ClientSocket *server.Socket

	//userData *model.ChatModel
	userID int64
	//State    int //0不在线 1在线
}

func (me *GameClient) readData(f func()) {
	if f == nil {
		return
	}
	me.RLock()
	defer me.RUnlock()
	f()
}
func (me *GameClient) writeData(f func()) {
	if f == nil {
		return
	}
	me.Lock()
	defer me.Unlock()
	f()
}

// WriteMsg .
func (me *GameClient) WriteMsg(msgcode int, args ...interface{}) error {
	if me.ClientSocket == nil {
		return nil
	}
	buff, bufLen := utils.AppendHead(utils.MergeBytes(args...))
	temp := make([]byte, 5+bufLen)

	binary.LittleEndian.PutUint32(temp[0:], uint32(msgcode))
	temp[4] = byte(0)
	copy(temp[5:], buff)
	// for i := len(args) - 1; i >= 0; i-- {
	// 	temp[i+2] = args[i]
	// }
	return me.ClientSocket.Write(temp)
}
