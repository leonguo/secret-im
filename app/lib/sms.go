package lib

import (
	"../../config"
	"math/rand"
	"time"
	"strconv"
	"github.com/mitchellh/mapstructure"
	"encoding/hex"
	"crypto/sha1"
)

type device struct {
	Number string
	Code string
}

func containTestNumber(number string) (bool bool, code string) {
	bool = false
	code = ""
	configDevices := config.AppConfig.Get("testDevices")
	for _, j := range configDevices.([]interface{}) {
		var device device
		 mapstructure.Decode(j,&device)
		 if device.Number == number {
		 	bool = true
		 	code = device.Code
		 	return
		 }
	}
	return
}

func VerificationCode(number string) (code string) {
	ok, code := containTestNumber(number)
	if !ok {
		code = ""
		rand.Seed(time.Now().UnixNano())
		for i := 1; i < 7; i++ {
			code = code + strconv.Itoa(rand.Intn(9))
		}
	}
	return
}

func GenSalt() (salt string) {
	rand.Seed(time.Now().UnixNano())
	salt = ""
	for i := 1; i < 10; i++ {
		salt = salt + strconv.Itoa(rand.Intn(9))
	}
	return
}

func GenAuthKey(salt string, password string) (authKey string) {
	h := sha1.New()
	h.Write([]byte(salt + password))
	authKey = hex.EncodeToString(h.Sum(nil))
	return
}