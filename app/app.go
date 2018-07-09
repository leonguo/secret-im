package app

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"./controllers"
	"../util"
)

var Server *echo.Echo

func Init() {
	Server = echo.New()
	Server.Debug = true
	// load default config
	Server.Use(middleware.Logger())
	Server.Use(middleware.Recover())
	Server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	Server.Use(middleware.RequestID())
	// 注册链接信息
	hub := util.NewHub()
	go hub.Run()
	Server.Use(
		func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				c.Set("Hub", hub)
				return next(c)
			}
		})
	Server.Static("/ws", "./public")
	Server.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "test!")
	})

	// api接口组
	accountApi := Server.Group("/v1/accounts")
	accountApi.GET("/sms/code/:number", controllers.GetSmsCode)

	//验证header的参数
	authApi := Server.Group("/v1")
	authApi.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		AuthScheme: "Basic",
		Validator: func(key string, c echo.Context) (bool, error) {
			checkOk, err := util.AuthorizationHeader(key, c)
			return checkOk, err
		},
	}))
	// accounts
	authApi.PUT("/accounts/code/:verification_code", controllers.PutVerifyAccount)
	authApi.PUT("/accounts/attributes", controllers.PutAttribute)

	// profile
	authApi.GET("/profile/:number", controllers.GetProfile)

	// directory
	authApi.GET("/directory/:token", controllers.GetDirectoryToken)
	authApi.PUT("/directory/tokens", controllers.PutDirectoryToken)

	// message
	authApi.PUT("/messages/:destination", controllers.GetMessageSend)
	authApi.GET("/messages", controllers.GetMessagePending)

	// devices
	authApi.GET("/devices", controllers.GetDevices)
	authApi.DELETE("/devices/:device_id", controllers.DeleteDevices)
	authApi.PUT("/devices/:verification_code", controllers.PutDevices)

	// attachments
	//authApi.GET("/attachments", controllers.GetAttachment)
	//authApi.GET("/attachments/:attachmentId", controllers.GetAttachmentId)
	// ws
	wsAuthApi := Server.Group("/v1")
	wsAuthApi.GET("/websocket/", controllers.GetWebSocket)
	wsAuthApi.GET("/websocket/provisioning/", controllers.GetWebSocketPrevisioning)

	api2 := Server.Group("/v2")
	api2.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		AuthScheme: "Basic",
		Validator: func(key string, c echo.Context) (bool, error) {
			checkOk, err := util.AuthorizationHeader(key, c)
			return checkOk, err
		},

	}))
	api2.GET("/keys", controllers.GetKeys)
	api2.GET("/keys/:number/:device_id", controllers.GetKeysDevice)
	api2.PUT("/keys", controllers.PutKeys)
	api2.GET("/keys/signed", controllers.GetKeysSigned)
	api2.PUT("/keys/signed", controllers.PutKeysSigned)
}
