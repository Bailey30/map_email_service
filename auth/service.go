package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type AuthService interface {
	validate(r *http.Request) (bool, error)
}

type authService struct{}

func (svc *authService) validate(r *http.Request) (bool, error) {
	authHeader := r.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// validate the token signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing methds: %v", token.Header["alg"])
		}
		// my auth secret key, same as front end, put in env file
		return []byte("secret"), nil
	})

	return token.Valid, err
}
