package tabledata

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"llog"
	"os"
)

// NotificationItem .
type NotificationItem struct {
	XMLName xml.Name `xml:"info"`

	ID         int32  `xml:"id,attr"`
	OperType   int32  `xml:"operType,attr"`
	OperID     string `xml:"operID,attr"`
	FromUserID int32  `xml:"FromUserID,attr"`
	Content    string `xml:"Content,attr"`
}

// Notification 系统消息的内容表
type Notification struct {
	XMLName xml.Name `xml:"root"`
	Version string   `xml:"version,attr"`

	Items []NotificationItem `xml:"info"`

	Hash map[int32]*NotificationItem
}

// NotificationTable .
var NotificationTable = Notification{}

// ReadAll 读取xml表
func ReadAll() {

	fmt.Println("1234")
	file, err := os.Open("./xml/notification.xml")
	if err != nil {
		llog.Error("%s.xml文件打开失败.", "notification")
		return
	}
	data, err := ioutil.ReadAll(file)
	file.Close()
	if err != nil {
		llog.Error("%s.xml读取内容失败.", "notification")
	}

	err = xml.Unmarshal(data, &NotificationTable)
	if err != nil {
		llog.Error("%s.xml解析失败:", "notification", err.Error())
	}

	NotificationTable.Hash = make(map[int32]*NotificationItem)
	items := NotificationTable.Items
	for i := len(items) - 1; i >= 0; i-- {
		NotificationTable.Hash[items[i].ID] = &items[i]
	}
}
