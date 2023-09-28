package controller

import (
	"transferSrv/controller/middleware"

	"github.com/gin-gonic/gin"
)

func Route() *gin.Engine {
	router := gin.Default()

	registerRoute(router)

	return router
}

func registerRoute(router *gin.Engine) {

	transGroup := router.Group("/transfer").Use(middleware.Auth())
	{
		transGroup.POST("/sol", handleTransferSOL)
		transGroup.POST("/token", handleTransferToken)
	}

}
