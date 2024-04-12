package main

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type logger struct {
	next AuthService
}

func NewLogger(next AuthService) AuthService {
	return &logger{
		next: next,
	}
}

func (s *logger) validate(r *http.Request) (valid bool, err error) {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:  true,
		PadLevelText: true,
	})
	defer func(begin time.Time) {
		logrus.WithFields(
			logrus.Fields{
				"took":  time.Since(begin),
				"error": err,
				"valid": valid,
			}).Info("auth service - validate")
	}(time.Now())

	return s.next.validate(r)
}
