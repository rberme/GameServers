package main

import (
	"encoding/binary"
	"llog"
	"main/channel"
	"main/model"
	"network"
	"sync"
	"time"
	"utils"
)

type agent struct {
	sync.RWMutex
	conn      network.Conn
	chat      *chat
	userData  model.ChatModel
	userID    int64
	loginTime time.Time
}

func (me *agent) readData(f func()) {
	if f == nil {
		return
	}
	me.RLock()
	defer me.RUnlock()
	f()
}
func (me *agent) writeData(f func()) {
	if f == nil {
		return
	}
	me.Lock()
	defer me.Unlock()
	f()
}

func (me *agent) Run() {
	for {
		data, err := me.conn.ReadMsg()
		if err != nil {
			llog.Debug("read message: %v", err)
			break
		}

		if me.chat.Processor != nil {
			msg, err := me.chat.Processor.Unmarshal(data)
			if err != nil {
				llog.Debug("unmarshal message error: %v", err)
				break
			}
			err = me.chat.Processor.Route(msg, me)
			if err != nil {
				llog.Debug("route message error: %v", err)
				break
			}
		}
	}
}

func (me *agent) OnClose() {
	var wc, tc, oc int32
	me.readData(func() {
		wc = me.userData.WorldChannel
		tc = me.userData.TeamID
		oc = me.userData.OrganizeID
	})
	if wc > 0 { //退出世界频道
		channel.GetoutWorldChannel(me.userID, wc)
	}
	if tc > 0 { //退出队伍频道
		channel.GetoutTeamChannel(me.userID, tc)
	}
	if oc > 0 { //退出社团频道
		channel.GetoutOrganizeChannel(me.userID, oc)
	}

	//me.conn = nil
	globalAgents.Del(me.userID)
}

func (me *agent) WriteMsg(msgcode int, args ...interface{}) {

	buff, bufLen := utils.AppendHead(utils.MergeBytes(args...))
	temp := make([]byte, 5+bufLen)

	binary.LittleEndian.PutUint32(temp[0:], uint32(msgcode))
	temp[4] = byte(0)
	copy(temp[5:], buff)
	// for i := len(args) - 1; i >= 0; i-- {
	// 	temp[i+2] = args[i]
	// }
	me.conn.WriteMsg(temp)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////

// type agentSet struct {
// 	agents sync.Map //map[int64]*agent
// }

// func newAgentSet() *agentSet {
// 	return &agentSet{agents: sync.Map{}}
// }

// func (me *agentSet) Add(a *agent) {
// 	me.agents.Store(a.userID, a)
// }

// func (me *agentSet) Del(id int64) {
// 	me.agents.Delete(id)
// }

// func (me *agentSet) Get(id int64) *agent {
// 	temp, ok := me.agents.Load(id)
// 	if ok {
// 		return temp.(*agent)
// 	}
// 	return nil
// }

// func (me *agentSet) Range(f func(id int64, ag *agent) bool) {
// 	if f == nil {
// 		return
// 	}
// 	me.agents.Range(func(key, value interface{}) bool {
// 		id := key.(int64)
// 		ag := value.(*agent)
// 		return f(id, ag)
// 	})
// }

const agentSetLength = 32

type agentSet struct {
	//sync.RWMutex
	agents [agentSetLength]sync.Map //map[int64]*agent
}

func newAgentSet() *agentSet {
	return &agentSet{}
}

func (me *agentSet) Add(id int64, a *agent) {
	// me.RLock()
	// defer me.RUnlock()
	idx := id % agentSetLength
	me.agents[idx].Store(id, a)
}

func (me *agentSet) Del(id int64) {
	idx := id % agentSetLength
	me.agents[idx].Delete(id)
}

func (me *agentSet) Get(id int64) *agent {
	idx := id % agentSetLength
	temp, ok := me.agents[idx].Load(id)
	if ok {
		return temp.(*agent)
	}
	return nil
}

func (me *agentSet) Range(f func(id int64, ag *agent) bool) {
	if f == nil {
		return
	}

	//for _, v := range me.agents { // range 是复制操作
	for i := 0; i < agentSetLength; i++ {
		me.agents[i].Range(func(key, value interface{}) bool {
			id := key.(int64)
			ag := value.(*agent)
			return f(id, ag)
		})
	}
}

/////////////////////////////////////////////////////////////

type tempSet struct {
	sync.RWMutex
	agents map[*agent]int8
}

func (me *tempSet) Add(a *agent) {
	me.Lock()
	defer me.Unlock()
	me.agents[a] = 0
}

func (me *tempSet) Del(ag *agent) {
	me.Lock()
	defer me.Unlock()
	delete(me.agents, ag)
}
