package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func SendEmail(ctx *gin.Context) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Erro ao carregar variavel de ambiente")
		return
	}

	app_password := os.Getenv("APP_PASSWORD")
	email := os.Getenv("EMAIL")
	request := SendEmailRequest{}

	ctx.Bind(&request)

	mail := gomail.NewMessage()
	mail.SetHeader("From", request.SenderEmail)
	mail.SetHeader("To", request.ReceiverEmail)
	mail.SetHeader("Subject", request.Subject)
	mail.SetBody("text/plain", request.Text)

	config := gomail.NewDialer("smtp.gmail.com", 587, email, app_password)

	if err := config.DialAndSend(mail); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to send email, try again later",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": "your email was send",
	})
}
