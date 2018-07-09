package controllers

import (
	"github.com/labstack/echo"
	"net/http"
	"../lib"
	"encoding/json"
	"../models"
	"../../db/redis"
)

func GetDirectoryToken(c echo.Context) error {
	token := c.Param("token")
	bt, err := lib.DecodeToken(token)
	if err != nil {
		return lib.Resp(c, http.StatusBadRequest, "", err)
	}
	value, err := redis.RedisDirectoryManager().HGet("directory", string(bt[:])).Result()
	if err != nil {
		return lib.Resp(c, http.StatusBadRequest, "", err)
	}

	var ds models.DirectoryStat
	err = json.Unmarshal([]byte(value), &ds)
	if err != nil {
		return lib.Resp(c, http.StatusBadRequest, "", err)
	}
	ds.Token = token
	var dr models.DirectoryResp
	dr = models.DirectoryResp(ds)
	return c.JSON(http.StatusOK, dr)
}

func PutDirectoryToken(c echo.Context) error {
	dt := new(models.DirectoryTokens)
	if err := c.Bind(dt); err != nil {
		return lib.Resp(c, http.StatusBadRequest, "param bind error", err)
	}

	var dss models.DirectoryStats
	for _, token := range dt.Contacts {
		c.Logger().Debug(token)
		bt, _ := lib.DecodeToken(token)
		value, err := redis.RedisDirectoryManager().HGet("directory", string(bt[:])).Result()
		if err != nil {
			c.Logger().Debug(" redis >>>>", err)
			continue
		}
		var ds models.DirectoryStat
		err = json.Unmarshal([]byte(value), &ds)
		if err != nil {
			c.Logger().Debug(" json >>>>", err)
			continue
		}
		ds.Token = token
		var dr models.DirectoryResp
		dr = models.DirectoryResp(ds)
		dss.Contacts = append(dss.Contacts, dr)
	}
	return c.JSON(http.StatusOK, dss)
}
