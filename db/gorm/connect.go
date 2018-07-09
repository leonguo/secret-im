package gorm

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"../../config"
	"fmt"
)

var db *gorm.DB
var msgDb *gorm.DB

func init() {
	var err error
	connectString := "host=" + config.AppConfig.GetString("accountdb.pg_host") + " port=" + config.AppConfig.GetString("accountdb.pg_port") + " user=" + config.AppConfig.GetString("accountdb.pg_user") + " dbname=" + config.AppConfig.GetString("accountdb.pg_dbname") + " password=" + config.AppConfig.GetString("accountdb.pg_password") + " sslmode=disable"
	db, err = gorm.Open("postgres", connectString)
	if err != nil {
		panic(fmt.Errorf("Fatal err when db connect: %s \n", err))
	}
	db.DB().SetMaxIdleConns(5)
	db.DB().SetMaxOpenConns(100)
	db.LogMode(true)

	msgConnectString := "host=" + config.AppConfig.GetString("messagedb.pg_host") + " port=" + config.AppConfig.GetString("messagedb.pg_port") + " user=" + config.AppConfig.GetString("messagedb.pg_user") + " dbname=" + config.AppConfig.GetString("messagedb.pg_dbname") + " password=" + config.AppConfig.GetString("messagedb.pg_password") + " sslmode=disable"
	msgDb, err = gorm.Open("postgres", msgConnectString)
	if err != nil {
		panic(fmt.Errorf("Fatal err when db connect: %s \n", err))
	}
	msgDb.DB().SetMaxIdleConns(5)
	msgDb.DB().SetMaxOpenConns(100)
	msgDb.LogMode(true)
	return
}

func ClosePg() {
	err := db.Close()
	if err != nil {
		panic(err)
	}
	return
}

func AccountManager() *gorm.DB {
	return db
}

func MsgManager() *gorm.DB {
	return msgDb
}
