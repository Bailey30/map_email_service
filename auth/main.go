package main

import ()

func main() {
	auth := &authService{}

	// logging middleware
	service := NewLogger(auth)

	server := NewJSONServer(service, ":3002")
	server.Run()
}
