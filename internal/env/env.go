package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	SERVER_PORT string

	MONGO_DB_NAME             string
	MONGO_COLLECTION_CONTESTS string
	MONGO_COLLECTION_GAMERS   string
	MONGO_ROOT_USER           string
	MONGO_ROOT_PASSWORD       string
	MONGO_PORT                string

	URL string
}

func InitENV() *Env {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	e := Env{
		SERVER_PORT:               os.Getenv("SERVER_PORT"),
		MONGO_DB_NAME:             os.Getenv("MONGO_DB_NAME"),
		MONGO_COLLECTION_CONTESTS: os.Getenv("MONGO_COLLECTION_CONTESTS"),
		MONGO_COLLECTION_GAMERS:   os.Getenv("MONGO_COLLECTION_GAMERS"),
		MONGO_ROOT_USER:           os.Getenv("MONGO_ROOT_USER"),
		MONGO_ROOT_PASSWORD:       os.Getenv("MONGO_ROOT_PASSWORD"),
		MONGO_PORT:                os.Getenv("MONGO_PORT"),
		URL:                       os.Getenv("URL"),
	}

	return &e
}
