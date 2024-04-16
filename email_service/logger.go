package main

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

type logger struct {
	next PasswordService
}

func NewLogger(next PasswordService) PasswordService {
	return &logger{
		next: next,
	}
}

func (s *logger) SendResetPasswordEmail(ctx context.Context, email string, token string) (e string, err error) {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:  true,
		PadLevelText: true,
	})
	defer func(begin time.Time) {
		logrus.WithFields(
			logrus.Fields{
				"took":  time.Since(begin),
				"error": err,
				"email": email,
			}).Info("auth service - validate")
	}(time.Now())

	return s.next.SendResetPasswordEmail(ctx, email, token)
}

func (s *logger) ValidateResetCode(ctx context.Context, reqBody ValidateTokenBody, db PasswordResetDB) (err error) {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:  true,
		PadLevelText: true,
	})
	defer func(begin time.Time) {
		logrus.WithFields(
			logrus.Fields{
				"took":        time.Since(begin),
				"error":       err,
				"requestBody": reqBody,
			}).Info("auth service - validate")
	}(time.Now())
	return s.next.ValidateResetCode(ctx, reqBody, db)
}
