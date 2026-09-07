package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found. Falling back to environment variables.")
	}

	app, err := InitializeApp()
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
