package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/smtp"
	"os"
)

type PasswordService interface {
	SendResetPasswordEmail(context.Context, string, string) (string, error)
}

type passwordService struct{}

func NewPasswordService() *passwordService {
	return &passwordService{}
}

func (s *passwordService) SendResetPasswordEmail(ctx context.Context, email string, token string) (string, error) {
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

func GenerateRandomToken() (string, error) {
	// Generate a byte slice of length 16 to hold random bytes
	tokenBytes := make([]byte, 16)

	// Fill the tokenBytes slice with random bytes
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}

	// Encode the random bytes to hexadecimal string representation
	return hex.EncodeToString(tokenBytes), nil
}

// Function to hash the token using SHA256
func HashToken(token string) string {
	// Create a new SHA256 hasher
	hasher := sha256.New()

	// Write the token bytes to the hasher
	hasher.Write([]byte(token))

	// Calculate the SHA256 hash of the token
	hashedToken := hasher.Sum(nil)

	// Encode the hashed token to hexadecimal string representation
	return hex.EncodeToString(hashedToken)
}
