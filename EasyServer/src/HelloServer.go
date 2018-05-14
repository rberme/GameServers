package main

import (
	"Chat"
	"Easy/Serializer"
	"Easy/Server"
	"fmt"
	"runtime"
	"sync"
	"time"
)

// type myRouter struct {
// }

// func (me *myRouter) HandleConnection(sess *Easy.Session) {
// 	fmt.Println("有客户端连接上了", sess.Conn.GetRemoteAddr())
// }

// func (me *myRouter) HandleClose(sess *Easy.Session) {
// 	fmt.Println("有客户端断了", sess.Conn.GetRemoteAddr())
// }

// //ReqData 接收数据的结构
// type ReqData struct {
// 	Data0 int
// 	// Data1  int8
// 	// Data2  int16
// 	// Data3  uint16
// 	// Data4  int
// 	// Data5  uint
// 	// Data6  int64
// 	// Data7  uint64
// 	// Data8  float32
// 	// Data9  string
// 	// Data10 float64
// }

// //HandleMessage 消息处理
// func (me *myRouter) HandleMessage(sess *Easy.Session, msg []byte) (string, interface{}, error) {
// 	//content := string(msg)
// 	//fmt.Println("接收到的数据:", content)
// 	data := new(ProtoFile.ChatData)
// 	sess.Decode(msg, data)
// 	log.Printf(data.Content)

// 	return "", nil, errors.New("End")
// 	// return "", nil, nil
// }

// type hello struct {
// }

// func (me *hello) TestFun(sess *Easy.Session, req map[string]interface{}) {
// }

func main() {
	fmt.Println("HelloWorld.")
	maxProcs := runtime.NumCPU()
	runtime.GOMAXPROCS(maxProcs)

	rwm := sync.RWMutex{}
	rwm.RLock()
	rwm.Lock()
	rwm.Unlock()
	rwm.RUnlock()

	// redis := &Storage.RedisStorage{}
	// redis.OpenRedis("redis://127.0.0.1:6382", 0)
	// buff, _ := msgpack.Marshal(3)
	// redis.Write("abc", "123", "100000", buff)

	server := Server.NewServer()
	chatHandler := &Chat.NormalChatHandler{
		Server:     server,
		Serializer: Serializer.NewProtobuf(),
	}

	go func() {
		for {
			fmt.Println("当前总Go程数:", runtime.NumGoroutine())
			server.PrintStatus()
			time.Sleep(time.Second * 3)
		}
	}()

	server.Run("192.168.0.189:8887", chatHandler)

	fmt.Println("ByeWorld.")
	//fmt.Scanf("%d")

}

type testEntity struct {
	id   int
	name string
}
