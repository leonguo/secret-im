package controllers

import (
	"github.com/labstack/echo"
	"../lib"
	"net/http"
	"../models"
	"encoding/json"
)

func GetKeys(c echo.Context) (err error) {
	// 获取设备信息
	account, err := lib.GetAccountInfo(c)
	if err != nil {
		return lib.Resp(c, http.StatusBadRequest, "get account info error", err)
	}
	device := &account.Devices[0]
	var key models.Keys
	count := key.GetKeysCount(account.Number, device.Id)
	data := make(map[string]int)
	data["count"] = count
	c.Logger().Debug(data)
	return c.JSON(http.StatusOK, data)
}

func GetKeysDevice(c echo.Context) (err error) {
	number := c.Param("number")
	// 获取设备信息
	account, err := lib.GetOtherAccountInfo(c, number)
	if err != nil {
		return lib.Resp(c, http.StatusBadRequest, "get account info error", err)
	}
	device := &account.Devices[0]
	data := make(map[string]interface{})
	data["identityKey"] = account.IdentityKey
	key := new(models.Keys)
	key.GetKeysFirst(account.Number, device.Id)
	preKey := make(map[string]interface{})
	preKey["keyId"] = key.KeyId
	preKey["publicKey"] = key.PublicKey
	var result = [1]map[string]interface{}{}
	result[0] = echo.Map{
		"deviceId":       device.Id,
		"registrationId": device.RegistrationId,
		"signedPreKey":   device.SignedPreKey,
		"preKey":         preKey,
	}
	data["devices"] = result
	return c.JSON(http.StatusOK, data)
}

func PutKeys(c echo.Context) (err error) {
	preKeys := new(models.PreKeyState)
	if err = c.Bind(preKeys); err != nil {
		return lib.Resp(c, http.StatusBadRequest, "param bind error", err)
	}
	// 获取设备信息
	account, err := lib.GetAccountInfo(c)
	if err != nil {
		return lib.Resp(c, http.StatusBadRequest, "get account info error", err)
	}
	device := &account.Devices[0]
	updateAccount := false
	if device.SignedPreKey != preKeys.SignedPreKey {
		device.SignedPreKey = preKeys.SignedPreKey
		updateAccount = true
	}
	if account.IdentityKey != preKeys.IdentityKey {
		updateAccount = true
		account.IdentityKey = preKeys.IdentityKey
	}

	// save account
	if updateAccount {
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
		hs := lib.ContactToken(account.Number)
		err = models.UpdateDirectory(hs, device.Voice, device.Video)
		if err != nil {
			return lib.Resp(c, http.StatusInternalServerError, "", err)
		}
	}
	// save keys
	var keys models.Keys
	keys.UpdateKeys(account.Number, device.Id, preKeys.PreKeys)
	return lib.RespOk(c, http.StatusOK, "")
}

func GetKeysSigned(c echo.Context) error {
	// 获取设备信息
	account, err := lib.GetAccountInfo(c)
	if err != nil {
		return lib.Resp(c, http.StatusBadRequest, "get account info error", err)
	}
	device := &account.Devices[0]
	data := device.SignedPreKey
	return c.JSON(http.StatusOK, data)
}

func PutKeysSigned(c echo.Context) error {
	signedPreKeys := new(models.SignedPreKey)
	if err := c.Bind(signedPreKeys); err != nil {
		return lib.Resp(c, http.StatusBadRequest, "param bind error", err)
	}
	// 获取设备信息
	account, err := lib.GetAccountInfo(c)
	if err != nil {
		return lib.Resp(c, http.StatusBadRequest, "get account info error", err)
	}
	device := &account.Devices[0]
	device.SignedPreKey.KeyId = signedPreKeys.KeyId
	device.SignedPreKey.PublicKey = signedPreKeys.PublicKey
	device.SignedPreKey.Signature = signedPreKeys.Signature

	// save account
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
	hs := lib.ContactToken(account.Number)
	err = models.UpdateDirectory(hs, device.Voice, device.Video)
	if err != nil {
		return lib.Resp(c, http.StatusInternalServerError, "", err)
	}
	return c.JSON(http.StatusOK, "")
}
