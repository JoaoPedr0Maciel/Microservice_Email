package router

import (
	"send-email-microservice/handler"

	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine) {
	v1 := router.Group("/v1")
	{
		v1.POST("/email/send", handler.SendEmail)
		v1.POST("/email/code", handler.SendVerifyCode)
		v1.POST("/email/confirm", handler.ConfirmVerificationCode)
	}
}
