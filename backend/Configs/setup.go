package Configs

import (
	"github.com/joho/godotenv"
	"log"
)

func IntializeDotEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Unable to initialize DotEnv (" + err.Error() + "). Fallback to native.")
	}
}
