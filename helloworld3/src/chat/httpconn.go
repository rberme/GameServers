package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	signKey = "qwertyuiopasdfgh"
)

type accountValidate struct {
	State  int    `json:"state"`
	UserID string `json:"UserID"`
	Sign   string `json:"sign"`
}

func accountServerValidate(token string) int64 {
	params := fmt.Sprintf("Token=%s", token)

	md5Ctx := md5.New()
	md5Ctx.Write([]byte(params + signKey))
	cipherStr := md5Ctx.Sum(nil)

	//temp := fmt.Sprintf("http://accountserver.com/AccountServer/Gateway/Validation.php?%s&Sign=%s", params, hex.EncodeToString(cipherStr))
	//fmt.Println(temp)
	resp, err := http.Get(fmt.Sprintf("http://accountserver.com/AccountServer/Gateway/Validation.php?%s&Sign=%s", params, hex.EncodeToString(cipherStr)))
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0
	}

	data := accountValidate{}
	//llog.Debug(string(body))

	// data.State = 0
	// data.UserID = 2
	// data.Sign = "815a5884399bb8eb46a09803c38c2d3c"
	// buf, err := json.Marshal(data)
	// fmt.Println(string(buf))

	err = json.Unmarshal(body, &data)
	if err != nil {
		return 0
	}
	retval, err := strconv.ParseInt(data.UserID, 10, 64) //data.UserID
	if err != nil {
		return 0
	}
	return retval
}

// func httpGet() {
// 	resp, err := http.Get("http://accountserver.com/AccountServer/Gateway/Validation.php?Token=55ce4822af33100491a0c7a0109e2671&Sign=80d88a8ca22937e34ee21d8218ec1bcd")
// 	if err != nil {

// 	}
// 	defer resp.Body.Close()
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {

// 	}
// 	llog.Debug(string(body))
// }
