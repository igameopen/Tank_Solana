package controller

import (
	"transferSrv/controller/dto"
	"transferSrv/infra/common"
	"transferSrv/infra/library/log"
	"transferSrv/service"

	"github.com/gin-gonic/gin"
)

// 转换SOL
func handleTransferSOL(c *gin.Context) {

	var (
		req  dto.TransferReq
		code int32
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("handleTransferSOL param err : ", err)
		RespFailed(c, nil, common.ERR_CODE_PARAM_ERR)
		return
	}

	sig, err := service.TransferSOL(req.Wallet, req.Amount)
	if err != nil {
		code = common.ERR_CODE_SYS_ERR
	} else {
		code = common.ERR_CODE_OK
	}

	Resp(c, sig, code)
}

// 转换token
func handleTransferToken(c *gin.Context) {

	var (
		req  dto.TransferReq
		code int32
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("handleTransferToken param err : ", err)
		RespFailed(c, nil, common.ERR_CODE_PARAM_ERR)
		return
	}

	sig, err := service.TransferToken(req.Wallet, req.Amount)
	if err != nil {
		code = common.ERR_CODE_SYS_ERR
	} else {
		code = common.ERR_CODE_OK
	}

	Resp(c, sig, code)
}
