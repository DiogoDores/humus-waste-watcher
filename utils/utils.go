package utils

import (
	"fmt"
	"strconv"

	dotenv "github.com/joho/godotenv"
)

func Stoi(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("failed to convert string to int: %w", err)
	}
	return i, nil
}

func LoadEnv() error {
	err := dotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("failed to load environment variables: %w", err)
	}
	return nil
}
