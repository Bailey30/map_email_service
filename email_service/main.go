package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	// "os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// "context"
// "fmt"
// "log"

func main() {
	// allow reading of env variables from env file
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal(envErr)
	}

	database_url := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable", os.Getenv("INSTANCE_CONNECTION_NAME"), os.Getenv("DB_USER"), os.Getenv("DB_NAME"))
	// create new database instance
	passwordResetDb := NewPasswordResetDB(database_url)
	// passwordResetDb := NewPasswordResetDB("postgresql://postgres:postgres@127.0.0.1:5432/passwordreset?sslmode=disable")

	err := passwordResetDb.CreateTable()
	if err != nil {
		log.Fatal(err)
	}

	// create new service instance
	svc := NewPasswordService()

	// apply logger middleware
	serviceWithLogger := NewLogger(svc)

	// create flag and parse any in command line
	listenAddr := flag.String("listenaddr", ":3001", "listen address the service is running")
	flag.Parse()

	// create new server with service and database
	server := NewJSONAPIServer(*listenAddr, serviceWithLogger, *passwordResetDb)
	server.Run()

}
