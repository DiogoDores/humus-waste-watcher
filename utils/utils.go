package utils

import (
	"log"
	dotenv "github.com/joho/godotenv"
)

func CheckError(msg string, e error) {
	if e != nil {
		log.Fatalf("%s: %v", msg, e)
	}
}

func LoadEnv() {
	err := dotenv.Load("../.env")
	CheckError("Failed to load environment variables", err)
}