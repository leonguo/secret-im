package lib

import (
	"github.com/labstack/echo"
	"../models"
	"encoding/json"
	"encoding/base64"
	"strings"
	"crypto/sha1"
	"bytes"
)

// 获取账号信息
func GetAccountInfo(c echo.Context) (account models.Account, err error) {
	number := c.Get("PhoneNumber").(string)
	accounts := new(models.Accounts)
	accounts.GetAccountByNumber(number)
	data := accounts.Data
	err = json.Unmarshal([]byte(data), &account)
	return
}

// 获取账号信息
func GetOtherAccountInfo(c echo.Context, phoneNumber string) (account models.Account, err error) {
	accounts := new(models.Accounts)
	accounts.GetAccountByNumber(phoneNumber)
	data := accounts.Data
	err = json.Unmarshal([]byte(data), &account)
	return
}

// 联系人decode

func DecodeToken(encoded string) (b []byte, err error) {
	encoded = strings.Replace(encoded, "-", "+", -1)
	encoded = strings.Replace(encoded, "_", "/", -1)
	return base64.RawStdEncoding.DecodeString(encoded)
}

// 更新通讯录
func ContactToken(number string) (hs []byte) {
	h := sha1.New()
	h.Write([]byte(number))
	hs = h.Sum(nil)
	var b bytes.Buffer
	b.Write(hs)
	b.Truncate(10)
	hs = b.Bytes()
	return
}

// 获取账号信息
func GetAccountInfoByNumber(number string) (account models.Account, err error) {
	accounts := new(models.Accounts)
	accounts.GetAccountByNumber(number)
	data := accounts.Data
	err = json.Unmarshal([]byte(data), &account)
	return
}
