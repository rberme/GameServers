package Storage

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/vmihailenco/msgpack"
)

// DataUpdater 数据更新
type DataUpdater struct {
	m            *sync.RWMutex
	entityUpdate map[string][]byte
}

var dataUpdaterIns *DataUpdater

// DataSyncInst 返回实例
func DataSyncInst() *DataUpdater {
	if dataUpdaterIns == nil {
		dataUpdaterIns = newDataUpdater()
	}
	return dataUpdaterIns
}

// NewDataUpdater 创建
func newDataUpdater() *DataUpdater {
	return &DataUpdater{
		m:            new(sync.RWMutex),
		entityUpdate: make(map[string][]byte),
	}
}

// SyncEntity 添加要更新的数据的key 序列化的buff
func (me *DataUpdater) SyncEntity(typeName string, key string, personalID int) {
	var k string
	var data interface{}
	if personalID >= 0 {
		k = fmt.Sprintf("%s,%s,%d", typeName, key, personalID)
		data = cachePoolIns.FindPersonalData(typeName, strconv.Itoa(personalID), key)

	} else {
		k = fmt.Sprintf("%s,%s", typeName, key)
		data = cachePoolIns.FindSharedData(typeName, key)
	}
	buff, _ := msgpack.Marshal(data)
	me.m.Lock()
	defer me.m.Unlock()
	me.entityUpdate[k] = buff
}

// 把修改的数据存入redis
func (me *DataUpdater) entitySave() {
	for {
		time.Sleep(time.Millisecond * 500)
		me.m.RLock()
		leng := len(me.entityUpdate)
		me.m.RUnlock()
		if leng == 0 {
			continue
		}

		temp := me.entityUpdate
		me.m.Lock()
		me.entityUpdate = make(map[string][]byte)
		me.m.Unlock()
		for k, v := range temp {
			ks := strings.Split(k, ",")
			typeName := ks[0]
			dataKey := ks[1]
			personalID := 100000
			if len(ks) == 3 { //公共的
				personalID, _ = strconv.Atoi(ks[2])
			}
			RedisInst.Write(typeName, dataKey, personalID, v)
		}
	}
}
