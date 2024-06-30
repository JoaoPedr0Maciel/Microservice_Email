package handler

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

// Estrutura para armazenar temporariamente os códigos de verificação
type VerificationStore struct {
	sync.Mutex
	codes map[string]int // Mapa para armazenar códigos por email
}

var verificationStore = VerificationStore{
	codes: make(map[string]int),
}

func SendVerifyCode(ctx *gin.Context) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Erro ao carregar variavel de ambiente")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":  "failed to load environment variables",
			"status": http.StatusInternalServerError,
		})
		return
	}

	appPassword := os.Getenv("APP_PASSWORD")
	email := os.Getenv("EMAIL")

	request := VerifyEmailWithCodeRequest{}

	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":  "invalid data",
			"status": http.StatusInternalServerError,
		})
		return
	}

	randomNumber := rand.Intn(900) + 100

	message := fmt.Sprintf("Seu código de verificação é %d", randomNumber)

	mail := gomail.NewMessage()
	mail.SetHeader("From", request.SenderEmail)
	mail.SetHeader("To", request.ReceiverEmail)
	mail.SetHeader("Subject", "Código de verificação")
	mail.SetBody("text/plain", message)

	config := gomail.NewDialer("smtp.gmail.com", 587, email, appPassword)

	if err := config.DialAndSend(mail); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":  "error to send the verification code",
			"status": http.StatusInternalServerError,
		})
		return
	}

	// Armazena o código de verificação com o email do destinatário no VerificationStore
	verificationStore.Lock()
	verificationStore.codes[request.ReceiverEmail] = randomNumber
	verificationStore.Unlock()

	ctx.JSON(http.StatusOK, gin.H{
		"success": "email was sent successfully",
		"status":  http.StatusOK,
	})
}

func ConfirmVerificationCode(ctx *gin.Context) {
	request := struct {
		ReceiverEmail string `json:"receiverEmail"`
		VerifyCode    int    `json:"verifyCode"`
	}{}

	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":  "invalid data",
			"status": http.StatusInternalServerError,
		})
		return
	}

	// Busca o código de verificação armazenado para o email do destinatário
	verificationStore.Lock()
	storedCode, ok := verificationStore.codes[request.ReceiverEmail]
	verificationStore.Unlock()

	if !ok || storedCode != request.VerifyCode {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":  "invalid verification code",
			"status": http.StatusBadRequest,
		})
		return
	}

	// Remove o código de verificação após a confirmação bem-sucedida
	verificationStore.Lock()
	delete(verificationStore.codes, request.ReceiverEmail)
	verificationStore.Unlock()

	ctx.JSON(http.StatusOK, gin.H{
		"success": "verification successful",
		"status":  http.StatusOK,
	})
}
