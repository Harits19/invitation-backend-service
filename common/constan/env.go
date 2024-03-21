package constan

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var OBJECT_STORAGE_KEY string

func InitEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	OBJECT_STORAGE_KEY = os.Getenv("OBJECT_STORAGE_KEY")

}
