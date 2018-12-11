package logic

import (
	"server"
	"sync"
	"time"
)

const mapGroup = 32

// ClientManager .
type ClientManager struct {
	//clientMaps [mapGroup]sync.Map
	clientMaps  [mapGroup]map[int64]*GameClient
	clientLocks [mapGroup]sync.RWMutex
}

// NewClientManager .
func NewClientManager() *ClientManager {
	temp := &ClientManager{}
	for i := 0; i < mapGroup; i++ {
		temp.clientMaps[i] = make(map[int64]*GameClient)
	}
	return temp
}

// GetClient 获取客户端
func (me *ClientManager) GetClient(uid int64) (*GameClient, bool) {
	groudID := uid % mapGroup
	clientMap := me.clientMaps[groudID]

	me.clientLocks[groudID].RLock()
	defer me.clientLocks[groudID].RUnlock()

	client, ok := clientMap[uid]
	if ok {
		client.accessTime = time.Now()
		return client, true
	}
	return nil, false
}

// SaveClient .
func (me *ClientManager) SaveClient(uid int64, s *server.Socket, factory func() *GameClient) (*GameClient, *server.Socket) {
	groudID := uid % mapGroup
	clientMap := me.clientMaps[groudID]

	me.clientLocks[groudID].Lock()
	defer me.clientLocks[groudID].Unlock()

	oldclient, ok := clientMap[uid]
	if ok {
		oldclient.accessTime = time.Now()
		oldSocket := oldclient.ClientSocket
		oldclient.ClientSocket = s
		//oldclient.State = 1
		return oldclient, oldSocket
	}
	newClient := factory()
	clientMap[uid] = newClient
	return newClient, nil

}

// ClearClient .
func (me *ClientManager) ClearClient(uid int64) {
	if uid <= 0 {
		return
	}
	groudID := uid % mapGroup

	clientMap := me.clientMaps[groudID]

	me.clientLocks[groudID].RLock()
	defer me.clientLocks[groudID].RUnlock()

	client, ok := clientMap[uid]
	if ok {
		client.accessTime = time.Now()
		client.ClientSocket = nil
		//client.State = 0
	}
}

// ClearExpired 遍历发送心跳并删除过期的玩家数据
func (me *ClientManager) ClearExpired() (clientNum int) {
	clientNum = 0
	now := time.Now().Unix()
	for i := 0; i < mapGroup; i++ {
		me.clientLocks[i].RLock()
		for k, v := range me.clientMaps[i] {
			del := false
			v.readData(func() {
				deltatime := now - v.accessTime.Unix()

				if deltatime > 60 { //超过1分钟心跳无响应强制断线 超过10分钟清理玩家数据
					if v.ClientSocket != nil {
						server.TCPManager.CloseSocket(v.ClientSocket)
						v.ClientSocket = nil
					}
					if deltatime > 600 {
						del = true
					} else {
						clientNum++
					}
					// if deltatime > 3600 {

				} else {
					if v.ClientSocket != nil {
						v.WriteMsg(MSG_CODE_HEARTBEAT_RET, now)
					}
					clientNum++
				}
			})
			if del {
				delete(me.clientMaps[i], k)
			}
		}
		me.clientLocks[i].RUnlock()
	}
	return
}

// Range 遍历
func (me *ClientManager) Range(f func(id int64, ag *GameClient) bool) {
	if f == nil {
		return
	}

	//for _, v := range me.agents { // range 是复制操作
	for i := 0; i < mapGroup; i++ {
		me.clientLocks[i].RLock()
		for _, v := range me.clientMaps[i] {
			v.readData(func() {
				if v.ClientSocket != nil {
					f(v.userID, v)
				}
			})
		}
		me.clientLocks[i].RUnlock()
	}
}

// // DoClientBackgroundWork .
// func (me *ClientManager) DoClientBackgroundWork() {
// 	// index := 0
// 	// for {
// 	// 	if client, ok := me.GetNextClient(index); ok {

// 	// 	}
// 	// }
// }

// // DoClientWorks .
// func (me *ClientManager) DoClientWorks() {

// }
