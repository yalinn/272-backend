package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	MONGO_URI      string
	MONGO_DBNAME   string
	REDIS_URI      string
	REDIRECT_URI   string
	JWT_SECRET_KEY string
	PORT           string = ":3000"
	BSL_URI        string
)

var (
	IMAP_S_HOST string
	IMAP_T_HOST string
	IMAP_PORT   string = "993"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file")
	} else {
		log.Println("Successfully loaded .env file")
	}
	MONGO_URI = os.Getenv("MONGO_URI")
	MONGO_DBNAME = os.Getenv("MONGO_DBNAME")
	REDIS_URI = os.Getenv("REDIS_URI")
	REDIRECT_URI = os.Getenv("REDIRECT_URI")
	JWT_SECRET_KEY = os.Getenv("JWT_SECRET_KEY")
	PORT = os.Getenv("PORT")
	BSL_URI = os.Getenv("BSL_URI")
	IMAP_S_HOST = os.Getenv("IMAP_S_HOST")
	IMAP_T_HOST = os.Getenv("IMAP_T_HOST")
	IMAP_PORT = os.Getenv("IMAP_PORT")

}

func Getenv(key string) string {
	if os.Getenv(key) == "" {
		log.Fatalf("Error getting %s from .env file", key)
	}
	return os.Getenv(key)
}
