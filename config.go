package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	OwnerID       int64
	TelegramToken string
	GeminiKey     string
)

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, falling back to environment variables:", err)
	}

	TelegramToken = os.Getenv("TELEGRAM_TOKEN")
	if TelegramToken == "" {
		log.Fatal("TELEGRAM_TOKEN is not set")
	}

	GeminiKey = os.Getenv("GEMINI_API_KEY")
	if GeminiKey == "" {
		log.Fatal("GEMINI_API_KEY is not set")
	}

	ownerIDStr := os.Getenv("OWNER_ID")
	if ownerIDStr == "" {
		log.Fatal("OWNER_ID is not set")
	}

	var parseErr error
	OwnerID, parseErr = strconv.ParseInt(ownerIDStr, 10, 64)
	if parseErr != nil {
		log.Fatal("Invalid OWNER_ID:", parseErr)
	}
}
