package utils

import (
	"encoding/binary"
	"llog"
)

// AppendHead 给bytes加上长度
func AppendHead(buf []byte) ([]byte, int32) {
	bufLen := len(buf)
	retval := make([]byte, bufLen+4)
	binary.LittleEndian.PutUint32(retval[0:], uint32(bufLen))
	copy(retval[4:], buf)
	return retval, int32(bufLen + 4)
}

// MergeBytes .
func MergeBytes(bufs ...interface{}) []byte {
	//4+1+4+buf
	var bufLen uint32
	var bufCount = len(bufs)
	//for _, buf := range bufs {
	for i := 0; i < bufCount; i++ {
		buf := bufs[i]
		switch buf.(type) {
		case uint32:
			bufLen += 4
		case int32:
			bufLen += 4
		case int8:
			bufLen++
		case byte:
			bufLen++
		case string:
			b := []byte(buf.(string))
			//to C# binary
			l := len(b)
			if l < 128 {
				bufLen++
			} else if l < 65536 {
				bufLen += 2
			} else {
				llog.Error("WriteMsg:未处理的字符串长度")
			}
			bufLen += uint32(l)
		case []byte:
			bufLen += uint32(len(buf.([]byte)))
		case int64:
			bufLen += 8
		default:
			llog.Error("WriteMsg:未处理的类型")
		}
	}

	arg := make([]byte, bufLen)
	// binary.LittleEndian.PutUint32(arg[0:], uint32(msgcode))
	// binary.LittleEndian.PutUint32(arg[5:], uint32(bufLen))

	l := 0
	//for _, buf := range bufs {
	for i := 0; i < bufCount; i++ {
		buf := bufs[i]
		switch buf.(type) {
		case uint32:
			binary.LittleEndian.PutUint32(arg[l:], buf.(uint32))
			l += 4
		case int32:
			binary.LittleEndian.PutUint32(arg[l:], uint32(buf.(int32)))
			l += 4
		case int8:
			arg[l] = byte(buf.(int8))
			l++
		case byte:
			arg[l] = buf.(byte)
			l++
		case string:
			b := []byte(buf.(string))
			//to C# binary
			le := len(b)
			if le < 128 {
				arg[l] = byte(le)
				l++
				copy(arg[l:], b)
				l += le
			} else if le < 65536 {
				arg[l] = byte(le%128 + 128)
				arg[l+1] = byte(le / 128)
				l += 2
				copy(arg[l:], b)
				l += le
			} else {
				llog.Error("WriteMsg:未处理的字符串长度")
			}
		case []byte:
			b := buf.([]byte)
			copy(arg[l:], b)
			l += len(b)
		case int64:
			binary.LittleEndian.PutUint64(arg[l:], uint64(buf.(int64)))
			l += 8
		default:
		}
	}

	return arg

}
