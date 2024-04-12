package main

import (
	"context"
)

type EmailService interface {
	SendResetPasswordEmail(context.Context, string) (string, error)
}

type emailService struct{}

func NewEmailService() *emailService {
	return &emailService{}
}

func (s *emailService) SendResetPasswordEmail(ctx context.Context, email string) (string, error) {
	return email, nil
}
