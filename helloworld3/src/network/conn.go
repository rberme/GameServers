package network

import (
	"net"
)

// Conn .
type Conn interface {
	ReadMsg() ([]byte, error)
	WriteMsg([]byte) error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
	//Destroy()
}
