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
	UserId int    `json:"userId"`
	Email  string `json:"email"`
}

type ValidateTokenBody struct {
	UserId int    `json:"userId"`
	Token  string `json:"token"`
}

type ResetEmailResponse struct {
	Email string `json:"email"`
}

type ValidateTokenResponse struct {
	Success bool `json:"success"`
}

type APIFunc func(context.Context, http.ResponseWriter, *http.Request) error

type JSONAPIServer struct {
	listenAddr string
	svc        PasswordService
	db         PasswordResetDB
}

func NewJSONAPIServer(listenAddr string, svc PasswordService, db PasswordResetDB) *JSONAPIServer {
	return &JSONAPIServer{
		listenAddr: listenAddr,
		svc:        svc,
		db:         db,
	}
}

func (s *JSONAPIServer) Run() {
	http.HandleFunc("/sendemail", makeHTTPHandlerFunc(s.handleSendPasswordResetEmail))
	http.HandleFunc("/validatetoken", makeHTTPHandlerFunc(s.handleValidateToken))

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
			fmt.Println("api err", err)
			// centralised error handling - not sure i like this. Cant specify status code easily
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		}
	}
}

func (s *JSONAPIServer) handleSendPasswordResetEmail(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// check for jwt in the gateway before this service is called

	// check if user email exists in the database - in the nextjs server

	// create new token and store hashed (SHA256) version in database
	// create link that is sent to user with the token
	// when the user sends their new password, the token is sent to this service to be validated, then the password is changed in the other database

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

	// create token
	token, err := GenerateRandomToken()
	if err != nil {
		return err
	}

	// created hashed version of the token for the database
	hashedToken := HashToken(token)

	// add hashed token to database
	createErr := s.db.create(CreateResetCodeParams{reqBody.UserId, hashedToken})
	if createErr != nil {
		return createErr
	}

	// business logic
	email, err := s.svc.SendResetPasswordEmail(ctx, reqBody.Email, token)
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

func (s *JSONAPIServer) handleValidateToken(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// parse the json
	reqBody := new(ValidateTokenBody)
	err = json.Unmarshal(body, &reqBody)
	if err != nil {
		fmt.Println(err)
		return err
	}

	validateErr := s.svc.ValidateResetCode(ctx, *reqBody, s.db)
	if validateErr != nil {
		return validateErr
	}

	// return that the token is valid
	// create the reponse value
	validateTokenResponse := ValidateTokenResponse{
		Success: true,
	}

	// encode and return json
	return writeJSON(w, http.StatusOK, &validateTokenResponse)
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
