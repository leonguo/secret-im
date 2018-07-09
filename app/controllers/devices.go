package controllers

import (
	"github.com/labstack/echo"
	"../lib"
	"net/http"
	"strconv"
	"encoding/json"
	"../models"
)

type devices struct {
	Name        string `json:"name"`
	Avatar      string `json:"avatar"`
	IdentityKey string `json:"identityKey"`
}

func GetDevices(c echo.Context) error {
	account, err := lib.GetAccountInfo(c)
	if err != nil {
		return lib.Resp(c, http.StatusNotFound, "get account info not found", err)
	}
	devices := account.Devices
	c.Logger().Debug(len(devices))
	deviceInfoList := make([]map[string]interface{}, len(devices)-1)
	for _, device := range devices {
		deviceInfo := map[string]interface{}{}
		deviceInfo["id"] = device.Id
		deviceInfo["name"] = device.Name
		deviceInfo["lastSeen"] = device.LastSeen
		deviceInfo["Created"] = device.Created
		deviceInfoList = append(deviceInfoList, deviceInfo)
	}
	return c.JSON(http.StatusOK, deviceInfoList)
}

func DeleteDevices(c echo.Context) error {
	deviceId := c.Param("device_id")
	if deviceId == "" || deviceId == "1" {
		return c.JSON(http.StatusUnauthorized, "")
	}
	id, err := strconv.ParseInt(deviceId, 10, 64)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "")
	}
	// 获取设备信息
	account, err := lib.GetAccountInfo(c)
	if err != nil {
		return lib.Resp(c, http.StatusBadRequest, "get account info error", err)
	}
	for index, device := range account.Devices {
		if device.Id == id {
			account.Devices = append(account.Devices[:index], account.Devices[index+1:]...)
			break
		}
	}
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
	return c.JSON(http.StatusOK, "")
}

func PutDevices(c echo.Context) error {
	//number := c.Get("PhoneNumber")
	password := c.Get("Password")
	attribute := new(models.AccountAttributes)
	if err := c.Bind(attribute); err != nil {
		return lib.Resp(c, http.StatusBadRequest, "param bind error", err)
	}
	// 获取设备信息
	account, err := lib.GetAccountInfo(c)
	if err != nil {
		return lib.Resp(c, http.StatusBadRequest, "get account info error", err)
	}
	nextId := len(account.Devices)
	device := models.Device{}
	device.Id = int64(nextId)
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
	account.Devices = devices.Devices

	accountString, _ := json.Marshal(account)

	accountM := models.Accounts{}
	accountM.Number = account.Number
	accountM.Data = string(accountString)
	err = accountM.UpdateAccount()
	if err != nil {
		return lib.Resp(c, http.StatusInternalServerError, "", err)
	}
	return c.JSON(http.StatusOK, "")
}
