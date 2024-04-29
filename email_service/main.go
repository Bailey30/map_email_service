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

	fmt.Println("db url", os.Getenv("DATABASE_URL"))
	// postgresql//postgres:postgres@127.0.0.1:5432/passwordreset?sslmode=disable
	// create new database instance
	passwordResetDb := NewPasswordResetDB(os.Getenv("DATABASE_URL"))
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
