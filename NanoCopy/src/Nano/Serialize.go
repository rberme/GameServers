package Nano

import (
	"Nano/Serialize"
	"Nano/Serialize/Protobuf"
)

// Default serializer
var serializer Serialize.Serializer = Protobuf.NewSerializer()

// SetSerializer customize application serializer, which automatically Marshal
// and UnMarshal handler payload
func SetSerializer(seri Serialize.Serializer) {
	serializer = seri
}
