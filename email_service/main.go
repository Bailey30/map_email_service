package main

import (
	"flag"
	"log"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// "context"
// "fmt"
// "log"

func main() {
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal(envErr)
	}

	passwordResetDb := NewPasswordResetDB("postgresql://postgres:postgres@localhost:5432/passwordreset?sslmode=disable")
	// passwordResetDb.connect()
	// passwordResetDb.CreateTable()
	// err := passwordResetDb.Create()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	svc := NewPasswordService()

	listenAddr := flag.String("listenaddr", ":3001", "listen address the service is running")
	flag.Parse()

	server := NewJSONAPIServer(*listenAddr, svc, *passwordResetDb)
	server.Run()

}
