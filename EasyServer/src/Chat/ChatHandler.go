package Chat

import (
	"Chat/Control"
	"Chat/Model"
	"Easy/BufferUtils"
	"Easy/Serializer"
	"Easy/Server"
	"encoding/binary"
	"fmt"
	"strconv"
	"time"
)

// NormalChatHandler 聊天消息处理
type NormalChatHandler struct {
	Server     *Server.TCPServer
	Serializer Serializer.ISerializer
}

//UserLogin 用户登入
func (me *NormalChatHandler) UserLogin(userID int) {
	fmt.Printf("............ 用户%d登入.\n", userID)
}

// UserLogout 用户登出
func (me *NormalChatHandler) UserLogout(userID int) {
	fmt.Printf("............ 用户%d登出.\n", userID)
}

// MainHandler 处理共享数据,这个函数不可以被卡住(chan)
func (me *NormalChatHandler) MainHandler(userID int, buff []byte) (m bool, t string, k string) {
	msgCode := int(binary.LittleEndian.Uint32(buff))
	buff = buff[4:]
	switch msgCode {
	case PublicLogin:
		// data := &Model.LoginModel{}
		// me.Serializer.Decode(buff, data)
		playerData := &Model.PlayerData{
			NickName:     string(buff),
			UserID:       userID,
			LastChatTime: time.Now(),
		}
		Control.PlayerCtrl.PlayerInfos[playerData.UserID] = playerData
		result := &Model.AllChaterInfo{
			Data: make([]*Model.ChaterInfo, len(Control.PlayerCtrl.PlayerInfos)),
		}
		i := 0
		for _, v := range Control.PlayerCtrl.PlayerInfos {
			result.Data[i] = &Model.ChaterInfo{
				Id:   int32(v.UserID),
				Nick: v.NickName,
			}
			i++
		}

		resultBytes, err := me.Serializer.Encode(result)
		if err != nil {
			break
		}
		resultBytes = BufferUtils.AppendHeadBytes(resultBytes)
		resultBytes = BufferUtils.AppendNumBytes(PublicLoginRET, resultBytes)

		resultBytes = BufferUtils.AppendHeadBytes(resultBytes)
		//me.Server.SendChan(userID, resultBytes)
		go me.Server.SendAllChan(resultBytes, nil)

	case PublicLogout:

	case PublicChating:
		data := &Model.ChatMessage{}

		me.Serializer.Decode(buff[4:], data)
		data.Nick = strconv.Itoa(userID)
		buff, _ := me.Serializer.Encode(data)
		buff = BufferUtils.AppendHeadBytes(buff)

		buff = BufferUtils.AppendNumBytes(PublicChatingRET, buff)
		buff = BufferUtils.AppendHeadBytes(buff)

		go me.Server.SendAllChan(buff, nil)
	default:
	}
	return m, t, k
}

// ClientHandler 处理私有,这个函数不可以被卡住(chan)
func (me *NormalChatHandler) ClientHandler(data []byte) (m bool, t string, k string, p int) {

	return m, t, k, p
}
