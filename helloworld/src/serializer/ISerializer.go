package serializer

//ISerializer 序列化接口
type ISerializer interface {
	Encode(interface{}) ([]byte, error)
	Decode([]byte, interface{}) error
}
