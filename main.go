package main

import "flag"

// "context"
// "fmt"
// "log"

func main() {
	svc := NewEmailService()
	//
	// email, err := svc.SendResetPasswordEmail(context.Background(), "test@email.com")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(email)

	listenAddr := flag.String("listenaddr", ":3001", "listen address the service is running")
	flag.Parse()

	server := NewJSONAPIServer(*listenAddr, svc)
	server.Run()

}
