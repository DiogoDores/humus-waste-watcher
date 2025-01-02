package utils

import (
	dotenv "github.com/joho/godotenv"
)

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func LoadEnv() {
	err := dotenv.Load()
	Check(err)
}