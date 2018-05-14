package Codec_test

import (
	"Nano/Internal/Codec"
	"Nano/Internal/Packet"
	"reflect"
	"testing"
)

func TestPack(t *testing.T) {
	data := []byte("hello world")
	p1 := &Packet.Packet{Type: Packet.Handshake, Data: data, Length: len(data)}
	pp1, err := Codec.Encode(Packet.Handshake, data)
	if err != nil {
		t.Error(err.Error())
	}

	d1 := Codec.NewDecoder()
	packets, err := d1.Decode(pp1)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(packets) < 1 {
		t.Fatal("packets should not empty")
	}
	if !reflect.DeepEqual(p1, packets[0]) {
		t.Fatalf("expect: %v, got: %v", p1, packets[0])
	}
}
