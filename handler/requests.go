package handler

type SendEmailRequest struct {
	SenderEmail   string
	ReceiverEmail string
	Subject       string
	Text          string
}

type VerifyEmailWithCodeRequest struct {
	SenderEmail   string
	ReceiverEmail string
}
