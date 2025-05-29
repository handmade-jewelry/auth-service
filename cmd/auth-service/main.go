package main

import (
	"context"
	"github.com/handmade-jewelry/auth-service/internal/app"
	"github.com/handmade-jewelry/auth-service/logger"
	"log"
)

func main() {
	err := logger.Init()
	if err != nil {
		log.Fatalf("cannot init logger: %v", err)
	}
	defer logger.Sync()

	ctx := context.Background()
	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	err = a.Run(ctx)
	if err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
