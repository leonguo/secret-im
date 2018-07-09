package main

import (
	"./app"
	"./config"
)

func main() {
	app.Init()
	app.Server.Logger.Fatal(app.Server.Start(config.AppConfig.GetString("system.port")))
}
