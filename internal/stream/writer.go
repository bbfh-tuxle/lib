package stream

import (
	"bytes"
	"encoding/binary"
	"io"
)

func WriteUint32(buffer io.Writer, number uint32) {
	byteSlice := make([]byte, 4)
	binary.BigEndian.PutUint32(byteSlice, number)
	buffer.Write(byteSlice)
}

func WriteUint64(buffer io.Writer, number uint64) {
	byteSlice := make([]byte, 8)
	binary.BigEndian.PutUint64(byteSlice, number)
	buffer.Write(byteSlice)
}

func WriteInt64(buffer io.Writer, number int64) {
	byteSlice := make([]byte, 8)
	binary.BigEndian.PutUint64(byteSlice, uint64(number))
	buffer.Write(byteSlice)
}

func WriteChars47(buffer io.Writer, str string) {
	var byteUser [47]byte
	copy(byteUser[:], []byte(str))
	buffer.Write(byteUser[:])
}

func WriteChars16(buffer io.Writer, str string) {
	var byteUser [16]byte
	copy(byteUser[:], []byte(str))
	buffer.Write(byteUser[:])
}

func WriteChars256(buffer io.Writer, str string) {
	var byteUser [256]byte
	copy(byteUser[:], []byte(str))
	buffer.Write(byteUser[:])
}

func WriteChars1024(buffer io.Writer, str string) {
	var byteUser [1024]byte
	copy(byteUser[:], []byte(str))
	buffer.Write(byteUser[:])
}

func WriteChars8192(buffer io.Writer, str string) {
	var byteUser [8192]byte
	copy(byteUser[:], []byte(str))
	buffer.Write(byteUser[:])
}

type writableAsString interface {
	Write(io.Writer)
}

func WriteToString(writer writableAsString) string {
	var buffer bytes.Buffer
	writer.Write(&buffer)
	return buffer.String()
}
