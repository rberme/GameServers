package main

import (
	"fmt"
	"nnet"
)

// func (me msgCodeType) String() string {

// }
type base struct {
}

func (me base) Name() {
	fmt.Println("base")
}

type sub struct {
	base
}

func (me sub) Name() {
	fmt.Println("sub")
}

func test(b base) {
	b.Name()
}

type customInt int

func main() {

	nnet.Start("192.168.0.189:8338", MsgHandler{})

	//fmt.Println(msgCodeLogin.String())
	fmt.Println("helloworld")
}
