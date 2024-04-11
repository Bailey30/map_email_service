package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
)

type ResetEmailBody struct {
	Email string `json:"email"`
}

type ResetEmailResponse struct {
	Email string `json:"email"`
}

type APIFunc func(context.Context, http.ResponseWriter, *http.Request) error

type JSONAPIServer struct {
	listenAddr string
	svc        EmailService
}

func NewJSONAPIServer(listenAddr string, svc EmailService) *JSONAPIServer {
	return &JSONAPIServer{
		listenAddr: listenAddr,
		svc:        svc,
	}
}

func (s *JSONAPIServer) Run() {
	http.HandleFunc("/", makeHTTPHandlerFunc(s.handleSendPasswordResetEmail))

	fmt.Printf("listening on port %s", s.listenAddr)
	http.ListenAndServe(s.listenAddr, nil)

}

// middleware
func makeHTTPHandlerFunc(apiFunc APIFunc) http.HandlerFunc {
	ctx := context.Background()

	ctx = context.WithValue(ctx, "requestID", rand.Intn(10000000))

	return func(w http.ResponseWriter, r *http.Request) {
		// set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")                            // Allow requests from any origin
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")          // Allow specified methods
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization") // Allow specified headers

		if err := apiFunc(ctx, w, r); err != nil {
			// centralised error handling
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		}
	}
}

func (s *JSONAPIServer) handleSendPasswordResetEmail(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// check for jwt
	jwt, err := validateJWT(w, r)
	if err != nil {
		return err
	}
	fmt.Println("jwt:", jwt)

	// read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// parse the json
	reqBody := new(ResetEmailBody)
	err = json.Unmarshal(body, &reqBody)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// business logic
	email, err := s.svc.SendResetPasswordEmail(ctx, reqBody.Email)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// create the reponse value
	resetEmailResponse := ResetEmailResponse{
		Email: email,
	}

	// encode and return json
	return writeJSON(w, http.StatusOK, &resetEmailResponse)
}

func writeJSON(w http.ResponseWriter, s int, v any) error {
	w.WriteHeader(s)
	return json.NewEncoder(w).Encode(v)
}

func validateJWT(w http.ResponseWriter, r *http.Request) (string, error) {
	// get the Authorization header from the request
	authHeader := r.Header.Get("Authorization")

	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		// exract token from the Authorization header
		token := strings.TrimPrefix(authHeader, "Bearer ")

		return token, nil
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}

	return "", errors.New("JWT token missing")
}
