package serializer

import (
	"errors"

	"github.com/golang/protobuf/proto"
)

var errWrongValueType = errors.New("protobuf: convert on wrong type value")

//ProtobufSerializer Protobuf序列化
type ProtobufSerializer struct {
}

//NewProtobuf 新建序列化对象
func NewProtobuf() *ProtobufSerializer {
	return &ProtobufSerializer{}
}

//Encode Protobuf编码
func (s *ProtobufSerializer) Encode(v interface{}) ([]byte, error) {
	pb, ok := v.(proto.Message)
	if !ok {
		return nil, errWrongValueType
	}
	return proto.Marshal(pb)
}

//Decode Protobuf解码
func (s *ProtobufSerializer) Decode(data []byte, v interface{}) error {
	pb, ok := v.(proto.Message)
	if !ok {
		return errWrongValueType
	}
	return proto.Unmarshal(data, pb)
}
