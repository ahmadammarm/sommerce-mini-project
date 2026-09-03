package main

import (
	"fmt"

	"github.com/ahmadammarm/sommerce-mini-project/config"
)

func main() {
	fmt.Println("Loading environment variables...")
	env := config.LoadEnv()

	fmt.Println("Connecting to the database...")
	db := config.InitDatabase(env)

	if db != nil {
		fmt.Println("SUCCESS: Database connection and pool established successfully!")
	}
}