package util

import (
	"encoding/base64"
	"crypto/aes"
	"crypto/cipher"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"log"
	"io"
	"crypto/rand"
)

func AESCBCEncrypt(plainText []byte, signalingKey string) (encryptPlainText []byte, err error) {
	//log.Printf("******** signalingKey >>>> %v", signalingKey)
	signalingKeyBytes, err := base64.StdEncoding.DecodeString(signalingKey)
	if err != nil {
		return nil, err
	}
	// 取32位 加密key
	cipherKey := make([]byte, 32)
	if len(signalingKeyBytes) < 32 {
		return nil, err
	}
	copy(cipherKey, signalingKeyBytes)
	//log.Printf(" ********  cipherkey >>>> %v", cipherKey)
	// 取HMAC key
	macKey := make([]byte, 20)
	log.Printf("signalingKeyBytes >>>> %v", signalingKeyBytes)
	copy(macKey, signalingKeyBytes[32:])
	//log.Printf(" ********  macKey >>>> %v", macKey)

	block, err := aes.NewCipher(cipherKey) //生成加密用的block
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	plainText = PKCS5Padding(plainText, blockSize)

	// 对IV有随机性要求，但没有保密性要求，所以常见的做法是将IV包含在加密文本当中
	ciphertext := make([]byte, aes.BlockSize+len(plainText))
	//随机一个block大小作为IV
	//采用不同的IV时相同的秘钥将会产生不同的密文，可以理解为一次加密的session
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	//log.Printf("******** iv >>>> %v", iv)
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plainText)

	// 版本号
	version := []byte{0x01}
	hmacNew := GenerateHMAC(BytesCombine(version, ciphertext), macKey)
	truncatedMac := make([]byte, 10)
	copy(truncatedMac, hmacNew)
	//log.Printf("******** truncatedMac >>>> %v", truncatedMac)
	encryptPlainText = BytesCombine(version, ciphertext, truncatedMac)
	//log.Printf("******** encryptPlainText >>>> %v", encryptPlainText)

	return

}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	//填充
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func GenerateHMAC(message, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	return mac.Sum(nil)
}

func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}
