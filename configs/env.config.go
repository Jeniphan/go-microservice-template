package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func InitConfigEnv() {
	envFile := ".env"

	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		log.Println(".env file not found, using environment variables instead")
		return
	}

	err := godotenv.Load(envFile)
	if err != nil {
		log.Printf("Warning: error loading %s file: %v", envFile, err)
		return
	}

	log.Println("Successfully loaded .env file")
}
