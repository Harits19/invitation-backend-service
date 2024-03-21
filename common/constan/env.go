package constan

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type IEnv struct {
	OBJECT_STORAGE_ACCESS_KEY string
	OBJECT_STORAGE_SECRET_KEY string
	OBJECT_STORAGE_REGION     string
	OBJECT_STORAGE_ENDPOINT   string
}

var ENV IEnv

func InitEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	ENV = IEnv{
		OBJECT_STORAGE_ACCESS_KEY: os.Getenv("OBJECT_STORAGE_ACCESS_KEY"),
		OBJECT_STORAGE_SECRET_KEY: os.Getenv("OBJECT_STORAGE_SECRET_KEY"),
		OBJECT_STORAGE_REGION:     os.Getenv("OBJECT_STORAGE_REGION"),
		OBJECT_STORAGE_ENDPOINT:   os.Getenv("OBJECT_STORAGE_ENDPOINT"),
	}

}
