package BufferUtils

import (
	"bytes"
	"encoding/binary"
)

// GetBytes 获得对象字节数组
func GetBytes(data interface{}) []byte {
	result := bytes.NewBuffer([]byte{})
	binary.Write(result, binary.LittleEndian, data)
	return result.Bytes()
}

// //MergeBytes 合并所有的对象为[]byte
// func MergeBytes(datas ...interface{}) []byte {
// 	result := bytes.NewBuffer([]byte{})
// 	temp := make([]byte, 4)
// 	for _, data := range datas {
// 		switch v:=data.(type) {
// 		case int8:
// 			binary.LittleEndian.PutUint32(lenBytes, uint32(temp))
// 		case int16:
// 			binary.LittleEndian.PutUint32(lenBytes, uint32(temp))
// 		case int32:
// 		}
// 		binary.Write(result, binary.LittleEndian, data)
// 	}
// 	return result.Bytes()
// }

// AppendNumBytes 在头部添加数据
func AppendNumBytes(num int, buff []byte) []byte {
	result := bytes.NewBuffer([]byte{})
	lenBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(lenBytes, uint32(num))
	result.Write(lenBytes)
	result.Write(buff)
	return result.Bytes()
}

// AppendHeadBytes 给buff头上加上数据长度
func AppendHeadBytes(buff []byte) []byte {
	result := bytes.NewBuffer([]byte{})

	lenBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(lenBytes, uint32(len(buff)+4))
	result.Write(lenBytes)
	result.Write(buff)
	// binary.Write(result, binary.LittleEndian, len(buff))
	// binary.Write(result, binary.LittleEndian, buff)
	return result.Bytes()
}
