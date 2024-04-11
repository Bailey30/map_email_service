package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	s := &http.Server{
		Addr:    "3002",
		Handler: mux,
	}

	mux.Handle("/validate", middleware(http.HandlerFunc(ValidateJTW)))

	err := s.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := next.ServeHTTP(w, r); err != nil {
			// will this work to centralise error handling?
		}
	})

}

func ValidateJTW(w http.ResponseWriter, r *http.Request) error {
	tokenString := r.Header.Get("Authorization")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// validate the token signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing methds: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})

	// WRITE RESPONSE HEADERS ETC
	// HANDERS DONT RETURN ANYTHING, JUST WRITE HEADERS

	if err != nil {
		return err, http.StatusUnauthorized
	}

	if !token.Valid {
		return nil, http.StatusUnauthorized
	}

	return

}
