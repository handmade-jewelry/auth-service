package main

import (
	app "github.com/handmade-jewellery/auth-service/internal/app"
	"log"
)

func main() {
	//ctx := context.Background()
	a, err := app.NewApp()
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	err = a.Run()
	if err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
