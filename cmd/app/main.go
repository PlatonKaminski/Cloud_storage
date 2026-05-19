package main

import (
	"cloud_storage/internal/app"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
