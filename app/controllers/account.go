package controllers

import (
	"github.com/labstack/echo"
	"net/http"
	"../../util"
	"../lib"
	"../models"
	"time"
	"encoding/json"
	"../../db/redis"
)

// 获取短信验证码
func GetSmsCode(c echo.Context) error {
	// get params
	number := c.Param("number")
	if isVal := util.IsValidNumber(number); !isVal {
		return lib.Resp(c, http.StatusBadRequest, "valid number false", number)
	}
	// TODO 限流控制

	// generate verification code
	code := lib.VerificationCode(number)

	// store code
	key := "pending_account2::" + number
	redis.RedisCacheManager().Set(key, code, 70*time.Second)

	var pendingAccounts = new(models.PendingAccounts)
	pendingAccounts.Number = number
	pendingAccounts.VerificationCode = code
	pendingAccounts.Timestamp = time.Now().Unix()
	pendingAccounts.SaveOrUpdatePendingAccounts()
	// send sms ignore
	return lib.RespOk(c, http.StatusCreated, "")
}

func PutVerifyAccount(c echo.Context) (err error) {
	verificationCode := c.Param("verification_code")
	number := c.Get("PhoneNumber")
	password := c.Get("Password")
	attribute := new(models.AccountAttributes)
	if err = c.Bind(attribute); err != nil {
		return lib.Resp(c, http.StatusBadRequest, "param bind error", err)
	}
	// check code
	key := "pending_account2::" + number.(string)
	code, err := redis.RedisDirectoryManager().Get(key).Result()
	if err != nil {
		return lib.Resp(c, http.StatusInternalServerError, "redis data not found", err)
	}
	if verificationCode != code {
		return lib.Resp(c, http.StatusBadRequest, "redis get data err", "")
	}
	device := models.Device{}
	device.Id = 1
	device.Salt = lib.GenSalt()
	device.RegistrationId = attribute.RegistrationId
	device.FetchesMessages = attribute.FetchesMessages
	device.SignalingKey = attribute.SignalingKey
	device.Video = attribute.Video
	device.Voice = attribute.Voice
	device.Name = attribute.Name
	device.AuthToken = lib.GenAuthKey(device.Salt, password.(string))
	device.Name = ""
	devices := models.Devices{}
	devices.Devices = append(devices.Devices, device)
	account := models.Account{}
	account.Number = number.(string)
	account.Devices = devices.Devices
	accountString, _ := json.Marshal(account)
	// create account
	clearAccount := models.Accounts{}
	clearAccount.Number = number.(string)
	clearAccount.DeleteAccount()

	accounts := models.Accounts{}
	accounts.Number = number.(string)
	accounts.Data = string(accountString)
	accounts.CreateAccount()
	pendingAccounts := models.PendingAccounts{}
	pendingAccounts.Number = number.(string)
	pendingAccounts.RemovePendingAccount()

	return lib.RespOk(c, http.StatusCreated, "")
}

func PutAttribute(c echo.Context) (err error) {
	attribute := new(models.AccountAttributes)
	if err = c.Bind(attribute); err != nil {
		return lib.Resp(c, http.StatusBadRequest, "param bind error", err)
	}
	account, err := lib.GetAccountInfo(c)
	if err != nil {
		return lib.Resp(c, http.StatusBadRequest, "get account info error", err)
	}
	devices := models.Devices{}
	devices.Devices = account.Devices
	device := &devices.Devices[0]
	device.RegistrationId = attribute.RegistrationId
	device.FetchesMessages = attribute.FetchesMessages
	device.SignalingKey = attribute.SignalingKey
	device.Video = attribute.Video
	device.Voice = attribute.Voice
	device.Name = attribute.Name
	accountString, err := json.Marshal(account)
	if err != nil {
		return lib.Resp(c, http.StatusBadRequest, "param Marshal account error", err)
	}
	accountM := models.Accounts{}
	accountM.GetAccountByNumber(account.Number)
	accountM.Data = string(accountString)
	err = accountM.UpdateAccount()
	if err != nil {
		return lib.Resp(c, http.StatusInternalServerError, "", err)
	}
	c.Logger().Debug(">>>> success update data")
	hs := lib.ContactToken(account.Number)
	c.Logger().Debug(">>>> success update data", string(hs))
	err = models.UpdateDirectory(hs, device.Voice, device.Video)
	if err != nil {
		return lib.Resp(c, http.StatusInternalServerError, "", err)
	}
	return lib.RespOk(c, http.StatusOK, "")
}
