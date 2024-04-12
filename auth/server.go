package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type APIFunc func(*http.Request) (bool, error)

type JSONServer struct {
	svc  AuthService
	port string
}

func NewJSONServer(svc AuthService, port string) *JSONServer {
	return &JSONServer{
		svc:  svc,
		port: port,
	}
}

func (s *JSONServer) Run() {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":3002",
		Handler: mux,
	}

	mux.Handle("/validate", middleware(http.HandlerFunc(s.ValidateJWT)))

	fmt.Println("auth server running")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

// handler function
func (s *JSONServer) ValidateJWT(w http.ResponseWriter, r *http.Request) {
	// businness logic from service
	tokenValid, err := s.svc.validate(r)

	// WRITE RESPONSE HEADERS ETC
	// HANDERS DONT RETURN ANYTHING, JUST WRITE HEADERS

	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
	}

	if !tokenValid {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "Token not valid."})
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "valid token"})
	return
}

func writeJSON(w http.ResponseWriter, statusCode int, value any) error {
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(value)
}
