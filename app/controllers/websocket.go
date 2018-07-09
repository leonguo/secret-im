package controllers

import (
	"github.com/labstack/echo"
	"github.com/gorilla/websocket"
	"../../util"
	"net/http"
	"strings"
)

var upgrader = websocket.Upgrader{}

func GetWebSocket(c echo.Context) error {
	number := c.QueryParam("login")
	password := c.QueryParam("password") // 注册conn
	if number == "" || password == "" {
		return c.JSON(http.StatusBadRequest, "")
	}
	if !strings.ContainsAny(number, "+") {
		number = "+" + strings.TrimSpace(number)
	}
	c.Logger().Debug(">> new conn >>>>>>", number)
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	defer conn.Close()
	hub := c.Get("Hub").(*util.Hub)
	client := &util.Client{Hub: hub, Conn: conn, Id: number, Send: make(chan []byte, util.MaxMessageSize)}
	client.Hub.Register <- client

	c.Logger().Debug("register client >>>> ", client)
	go client.WritePump()
	go client.ReadPump()
	select {}
}

func GetWebSocketPrevisioning(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}
