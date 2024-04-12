package main

import (
	"context"
	"fmt"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

type EmailService interface {
	SendResetPasswordEmail(context.Context, string) (string, error)
}

type emailService struct{}

func NewEmailService() *emailService {
	return &emailService{}
}

func (s *emailService) SendResetPasswordEmail(ctx context.Context, email string) (string, error) {
	envErr := godotenv.Load()
	if envErr != nil {
		return "", envErr
	}
	// auth credentials
	password := os.Getenv("APP_PASSWORD")
	senderEmail := "guinnessmapservices@gmail.com"
	host := "smtp.gmail.com"
	port := "587"
	fullServerAddress := host + ":" + port

	// message contents
	subject := "Reset password"
	mailBody := "Here will be the password reset link"

	recipient := []string{email}

	// compose message
	message := []byte("To: " + recipient[0] + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + mailBody + "\r\n")

	// set up authentication
	auth := smtp.PlainAuth("", senderEmail, password, host)

	// connect to the SMTP server
	err := smtp.SendMail(fullServerAddress, auth, senderEmail, recipient, message)
	if err != nil {
		return "", err
	}

	fmt.Println("email sent to", email)

	return email, nil
}
