package Protobuf

import (
	"errors"

	"github.com/golang/protobuf/proto"
)

// ErrWrongValueType 用protobuf序列化产生问题的错误
var ErrWrongValueType = errors.New("protobuf: convert on wrong type value")

// Serializer implements the serialize.Serializer interface
type Serializer struct{}

// NewSerializer returns a new Serializer.
func NewSerializer() *Serializer {
	return &Serializer{}
}

// Marshal returns the protobuf encoding of v.
func (s *Serializer) Marshal(v interface{}) ([]byte, error) {
	pb, ok := v.(proto.Message)
	if !ok {
		return nil, ErrWrongValueType
	}
	return proto.Marshal(pb)
}

// Unmarshal parses the protobuf-encoded data and stores the result
// in the value pointed to by v.
func (s *Serializer) Unmarshal(data []byte, v interface{}) error {
	pb, ok := v.(proto.Message)
	if !ok {
		return ErrWrongValueType
	}
	return proto.Unmarshal(data, pb)
}
