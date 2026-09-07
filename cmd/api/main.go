package main

import "log"

func main() {
	app, err := InitializeApp()
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
