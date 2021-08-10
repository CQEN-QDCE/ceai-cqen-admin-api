package main

import (
	"log"

	cmd "github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/commands"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cmd.Execute()
}
