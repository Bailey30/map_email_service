package main

import ()

func main() {
	auth := &authService{}

	// logging middleware
	service := NewLogger(auth)

	server := NewJSONServer(service, ":3002")
	server.Run()
}

// func middleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		next.ServeHTTP(w, r)
// 	})
// }
//
// func ValidateJTW(w http.ResponseWriter, r *http.Request) {
// 	authHeader := r.Header.Get("Authorization")
// 	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
// 	fmt.Println(tokenString)
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		// validate the token signing method
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			fmt.Println("not ok parsing token")
// 			return nil, fmt.Errorf("Unexpected signing methds: %v", token.Header["alg"])
// 		}
// 		// my auth secret key, same as front end, put in env file
// 		return []byte("secret"), nil
// 	})
//
// 	fmt.Println("jwt ", token)
//
// 	// WRITE RESPONSE HEADERS ETC
// 	// HANDERS DONT RETURN ANYTHING, JUST WRITE HEADERS
//
// 	if err != nil {
// 		fmt.Println(err)
// 		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
// 	}
//
// 	if !token.Valid {
// 		fmt.Println("token not valid")
// 		writeJSON(w, http.StatusUnauthorized, map[string]any{"": ""})
// 	}
//
// 	writeJSON(w, http.StatusOK, map[string]string{"message": "valid token"})
// 	return
// }
//
// func writeJSON(w http.ResponseWriter, statusCode int, value any) error {
// 	w.WriteHeader(statusCode)
// 	return json.NewEncoder(w).Encode(value)
// }
