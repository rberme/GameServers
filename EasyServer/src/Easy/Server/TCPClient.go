package Server

import (
	"net"
	"sync"
)

/////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////////////////

//TCPConn Tcp客户端连接
type TCPConn struct {
	ID              int
	conn            net.Conn
	remoteAddr      string
	RecvChan        chan *chanMsg
	server          *TCPServer
	NoDealHeartbeat int
	IsClosed        bool
}

//Close 关闭客户端连接
func (me *TCPConn) Close() bool {
	if me.IsClosed == false {
		me.conn.Close()
		close(me.RecvChan)
		me.IsClosed = true
		return true
	}
	return false
}

//GetRemoteAddr 客户端地址
func (me *TCPConn) GetRemoteAddr() string {
	return me.remoteAddr
}

// ConnMap 客户端容器
type ConnMap struct {
	sync.RWMutex
	conns        map[int]*TCPConn
	connAddCount int
}

// NewConnMap 创建带锁的map
func NewConnMap() *ConnMap {
	return &ConnMap{
		conns:        make(map[int]*TCPConn),
		connAddCount: 0,
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Length 返回元素数量
func (me *ConnMap) Length() int {
	me.RLock()
	defer me.RUnlock()
	return len(me.conns)
}

// Add 添加元素
func (me *ConnMap) Add(conn *TCPConn) {
	me.Lock()
	defer me.Unlock()
	conn.ID = me.connAddCount
	me.connAddCount++
	value, exist := me.conns[conn.ID]
	if exist == true {
		value.Close()
	}
	me.conns[conn.ID] = conn
}

// Set 设置元素值
func (me *ConnMap) Set(key int, conn *TCPConn) {
	me.Lock()
	defer me.Unlock()
	me.conns[key] = conn
}

// // Remove 删除一个元素
// func (me *ConnMap) Remove(id int) {
// 	me.sessM.RLock()
// 	_, exist := me.conns[id]
// 	me.sessM.RUnlock()
// 	if exist == true {
// 		me.sessM.Lock()
// 		_, exist := me.conns[id]
// 		if exist == true {
// 			me.conns[id].Close()
// 			delete(me.conns, id)
// 		}
// 		me.sessM.Unlock()
// 	}
// }

// RemoveConn 删除某个连接
func (me *ConnMap) RemoveConn(conn *TCPConn) {
	id := conn.ID
	me.RLock()
	value, exist := me.conns[id]
	me.RUnlock()
	if exist == true && value == conn {
		me.Lock()
		value, exist = me.conns[id]
		if exist == true && value == conn {
			me.conns[id].Close()
			delete(me.conns, id)
		}
		me.Unlock()
	}
}

//ForeachRead 遍历
func (me *ConnMap) ForeachRead(action func(sess *TCPConn), ids ...int) {
	me.RLock()
	defer me.RUnlock()
	if ids == nil || len(ids) == 0 || len(ids) == len(me.conns) {
		for _, value := range me.conns {
			action(value)
		}
	} else {
		for _, id := range ids {
			action(me.conns[id])
		}
	}
}

// //ForeachRead 遍历
// func (me *ConnMap) ForeachRead(action func(sess *TCPConn)) {
// 	me.sessM.RLock()
// 	defer me.sessM.RUnlock()
// 	for _, value := range me.conns {
// 		action(value)
// 	}
// }

// // ForsubRead 遍历部分
// func (me *ConnMap) ForsubRead(action func(sess *TCPConn), ids ...int) {
// 	me.sessM.RLock()
// 	defer me.sessM.RUnlock()
// 	for _, id := range ids {
// 		action(me.conns[id])
// 	}
// }

// Get 获取一个元素
func (me *ConnMap) Get(id int) *TCPConn {
	me.RLock()
	defer me.RUnlock()
	result, exist := me.conns[id]
	if exist {
		return result
	}
	return nil
}

// ContainKey 是否包含某个id
func (me *ConnMap) ContainKey(id int) bool {
	me.RLock()
	defer me.RUnlock()
	_, exist := me.conns[id]
	return exist
}

// // ClearExpired 清除过期的连接
// func (me *ConnMap) ClearExpired(nowTime int64) {
// 	me.sessM.Lock()
// 	defer me.sessM.Unlock()
// 	for key, value := range me.conns {
// 		if (nowTime - value.aliveTime.Unix()) > int64(time.Second*120) {
// 			delete(me.conns, key)
// 			value.Close()
// 		}
// 	}
// }

// ClearConn 清除过期的连接
func (me *ConnMap) ClearConn(ids []int) {
	me.Lock()
	defer me.Unlock()
	for _, key := range ids {
		me.conns[key].Close()
		delete(me.conns, key)
	}
}
