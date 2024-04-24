package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/smtp"
	"os"
	"time"
)

type PasswordService interface {
	SendResetPasswordEmail(context.Context, string, string) (string, error)
	ValidateResetCode(context.Context, ValidateTokenBody, PasswordResetDB) (ResetCode, error)
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
	mailBody := os.Getenv("GUINNESS_MAP_URL") + "/resetpassword/" + token

	fmt.Println("guinness map url", os.Getenv("GUINNESS_MAP_URL"))

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

func (s *passwordService) ValidateResetCode(ctx context.Context, reqBody ValidateTokenBody, db PasswordResetDB) (ResetCode, error) {
	// hash the token from the request to compare with what is in the database
	hashedToken := HashToken(reqBody.Token)

	// get token and user id from db using hashed token from token in reqBody - cannot send user id in request as they wont be logged in
	resetToken, err := db.getByToken(hashedToken)
	if err != nil {
		return ResetCode{}, err
	}

	// get any tokens from the database using the user id
	resetTokens, err := db.getAllById(resetToken.UserId)
	if err != nil {
		return ResetCode{}, err
	}

	// loop over the tokens in the database to find a match
	var matchedToken ResetCode
	for _, tokenInDb := range resetTokens {
		if tokenInDb.HashedCode == hashedToken {
			matchedToken = tokenInDb
		}
	}

	// if not token has been matched return error - may have already been redeemed
	if matchedToken == (ResetCode{}) {
		fmt.Println("Invalid password reset token")
		return ResetCode{}, errors.New("Invalid password reset token")
	}

	// check if the token has not expired
	if matchedToken.Expiry.Before(time.Now()) {
		fmt.Println("Token has expired")
		return ResetCode{}, errors.New("Token has expired")
	}

	// delete anything in the database related to user
	deleteErr := db.deleteAllById(resetToken.UserId)
	if deleteErr != nil {
		return ResetCode{}, err
	}

	return resetToken, nil
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
