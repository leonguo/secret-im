package controllers

import (
	"github.com/labstack/echo"
	"../lib"
	"net/http"
)

type profile struct {
	Name        string `json:"name"`
	Avatar      string `json:"avatar"`
	IdentityKey string `json:"identityKey"`
}

func GetProfile(c echo.Context) error {
	number := c.Param("number")
	account, err := lib.GetOtherAccountInfo(c, number)
	if err != nil {
		return lib.Resp(c, http.StatusNotFound, "get account info not found", err)
	}
	pro := profile{}
	pro.Name = account.Name
	pro.Avatar = account.Avatar
	pro.IdentityKey = account.IdentityKey

	return c.JSON(http.StatusOK, pro)
}
