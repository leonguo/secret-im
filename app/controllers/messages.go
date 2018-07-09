package controllers

import (
	"github.com/labstack/echo"
	"net/http"
	"../lib"
	"../models"
	"../../util"
	"time"
	"github.com/golang/protobuf/proto"
	"log"
	"encoding/base64"
)

// 发送消息
func GetMessageSend(c echo.Context) error {
	destination := c.Param("destination")
	c.Logger().Debug("get param destination", destination)
	inMessages := new(models.IncomingMessageList)
	if err := c.Bind(inMessages); err != nil {
		return lib.Resp(c, http.StatusBadRequest, "", err)
	}
	c.Logger().Debug("http get message >>>>>>>>>> ", inMessages)
	// 获取设备信息
	//source, err := lib.GetAccountInfo(c, accountDb)
	//if err != nil {
	//	return lib.Resp(c, http.StatusBadRequest, "get account info error", err)
	//}
	//var isSyncMessage bool
	//isSyncMessage = false
	//if destination == source.Number {
	//	isSyncMessage = true
	//}
	//if inMessages.Relay == "" {
	//	destinationAccount, _ := lib.GetOtherAccountInfo(c, accountDb, destination)
	//	destinationDevice := destinationAccount.Devices[0]
	//
	//	send, messageContent := sendLocalMessage(c, source, destination, destinationDevice, inMessages, isSyncMessage)
	//	c.Set("Message", messageContent)
	//	if !send {
	//		// 发送离线消息
	//		messageDb := db.ConnectMessagePG()
	//		defer messageDb.Close()
	//		message := new(models.OutGoingMessage)
	//		message.Content = messageContent.Message
	//		message.Timestamp = time.Now().UnixNano() / 1e6
	//		message.Type = inMessages.Messages[0].Type
	//		message.SourceDevice = source.Devices[0].Id
	//		message.Source = source.Number
	//		message.Destination = destination
	//		message.SaveMessage(messageDb)
	//	}
	//}
	return c.JSON(http.StatusOK, "")

}

// 获取离线消息内容
func GetMessagePending(c echo.Context) error {
	// 获取设备信息
	account, err := lib.GetAccountInfo(c)
	if err != nil {
		return lib.Resp(c, http.StatusBadRequest, "get account info error", err)
	}
	device := &account.Devices[0]
	message := new(models.OutGoingMessage)
	if device.Id == 0 {
		return c.JSON(http.StatusOK, "")
	}
	messages := message.GetMessageForDevice(account.Number, device.Id)
	return c.JSON(http.StatusOK, messages)
}

func sendLocalMessage(c echo.Context, source models.Account, destination string, destinationDevice models.Device, inMessages *models.IncomingMessageList, isSyncMessage bool) (send bool, singleMsg util.SingleMessage) {
	// 校验设备列表
	// 校验注册ID
	if inMessages.Timestamp == 0 {
		inMessages.Timestamp = time.Now().UnixNano() / 1e6
	}
	for _, message := range inMessages.Messages {
		//body, _ := base64.StdEncoding.DecodeString(message.Body)
		content, _ := base64.StdEncoding.DecodeString(message.Content)
		outMessage := &models.OutMessage{}
		//outMessage.Message = body
		outMessage.Timestamp = proto.Int64(inMessages.Timestamp)
		outMessage.SourceDevice = proto.Int64(source.Devices[0].Id)
		outMessage.Source = proto.String(source.Number)
		outMessage.Content = content
		outMessage.Type = models.OutMessage_Type(int32(message.Type)).Enum()
		out, err := proto.Marshal(outMessage)
		if err != nil {
			log.Println("Failed to encode message:", err)
		}
		c.Logger().Debug(">>>>> message content ... ", outMessage)
		singleMsg := new(util.SingleMessage)
		singleMsg.Id = destination
		singleMsg.Message, _ = util.AESCBCEncrypt(out, destinationDevice.SignalingKey)
		// send message
		log.Printf("v Message v>>>>  %v", singleMsg.Message)
		//hub := c.Get("Hub").(*util.Hub)
		//if len(hub.Clients) > 0 && len(singleMsg.Id) > 0 {
		//	for c, _ := range hub.Clients {
		//		if c.Id == singleMsg.Id {
		//			c.Hub.SendSingle <- singleMsg
		//			send = true
		//			break
		//		}
		//	}
		//}
	}
	return
}
