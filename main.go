package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	GeminiApiKey := os.Getenv("GEMINI_API_KEY")
	fmt.Println(GeminiApiKey)
}
