package utils

import (
	"log"
	"strconv"
	dotenv "github.com/joho/godotenv"
)

func CheckError(msg string, e error) {
	if e != nil {
		log.Fatalf("%s: %v", msg, e)
	}
}

func Stoi(s string) int {
	i, err := strconv.Atoi(s)
	CheckError("Failed to convert string to int", err)
	return i
}

func LoadEnv() {
	err := dotenv.Load(".env")
	CheckError("Failed to load environment variables", err)
}