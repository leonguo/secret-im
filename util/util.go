package util

import (
	"crypto/rand"
	"bytes"
	"encoding/binary"
)

// 生成附件ID
func GenerateAttachmentId() (int64, error) {
	bytesArray := make([]byte, 8)
	_, err := rand.Read(bytesArray)
	if err != nil {
		return 0, err
	}
	bytesArray[0] = bytesArray[0] & 0x7F
	return BytesToInt(bytesArray), nil
}

//整形转换成字节
func IntToBytes(n int64) []byte {
	tmp := int64(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, tmp)
	return bytesBuffer.Bytes()
}

//字节转换成整形
func BytesToInt(b []byte) int64 {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp int64
	binary.Read(bytesBuffer, binary.BigEndian, &tmp)
	return int64(tmp)
}
