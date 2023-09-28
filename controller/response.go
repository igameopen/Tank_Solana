package controller

import (
	"net/http"

	"transferSrv/infra/common"

	"github.com/gin-gonic/gin"
)

type RespObj struct {
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
	Code   int32       `json:"code"`
	Status string      `json:"status"`
}

func Resp(c *gin.Context, data interface{}, code int32) {
	if code != common.ERR_CODE_OK {
		RespFailed(c, data, code)
		return
	}
	RespOK(c, data)
}

func RespOK(c *gin.Context, data interface{}) {
	d := RespObj{
		Data:   data,
		Status: "ok",
	}

	c.JSON(http.StatusOK, d)
}

func RespFailed(c *gin.Context, data interface{}, code int32) {
	d := RespObj{
		Data:   data,
		Code:   code,
		Status: "failed",
		Msg:    common.ErrMsg[code],
	}

	c.JSON(http.StatusOK, d)
}

func RespDownload(c *gin.Context, fileName string) {
	c.File(fileName)
}
