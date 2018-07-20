package main

import (
	"fmt"
	"net"
	"time"
)

var buff = []byte{46, 0, 0, 0, 117, 39, 0, 0, 0, 33, 0, 0, 0, 32, 98, 98, 51, 54, 55, 97, 98, 50, 50, 102, 98, 50, 99, 52, 99, 49, 51, 101, 56, 53, 48, 102, 98, 98, 54, 50, 56, 56, 54, 51, 51, 51}

func main() {
	fmt.Println("hello world.")
	for i := 0; i < 1000; i++ {
		//time.Sleep(1 * time.Second)
		go tcpClientTest()
	}

	time.Sleep(40 * time.Second)
}

func tcpClientTest() {
	conn, err := net.Dial("tcp", "192.168.0.189:19001")
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	fmt.Println(conn.LocalAddr())
	conn.Write(buff)
	time.Sleep(30 * time.Second)

}
