package main

import (
	"Nano"
	"Nano/Component"
	"Nano/Serialize/Protobuf"
	"fmt"
	Console "fmt"
	"strings"
)

type (
	// Room 房间
	Room struct {
		Component.Base
	}
)

// NewRoom returns a new room
func NewRoom() *Room {
	return &Room{}
}

// AfterInit 初始化完成之后
func (me *Room) AfterInit() {
	fmt.Println("Room::AfterInit")
}

func main() {

	// s := "手游"
	// fmt.Println(reflect.TypeOf(s))
	// fmt.Println(reflect.ValueOf(s))
	// fmt.Println(reflect.Indirect(reflect.ValueOf(s)))
	// return
	Nano.SetSerializer(Protobuf.NewSerializer())

	room := NewRoom()
	Nano.Register(room,
		Component.WithName("room"),
		Component.WithNameFunc(strings.ToLower),
	)

	Nano.Listen(":8887")
	Console.Println("end")
}
