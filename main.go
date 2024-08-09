package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()

	godotenv.Load()

	GeminiApiKey := os.Getenv("GEMINI_API_KEY")

	client, err := genai.NewClient(ctx, option.WithAPIKey(GeminiApiKey))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	resp, err := model.GenerateContent(ctx, genai.Text("Write a story about a AI and magic"))
	if err != nil {
		log.Fatal(err)
	}

	for res := range resp.Candidates {
		fmt.Println(res)
	}

	fmt.Println(resp.Candidates[0].Content.Parts[0])

}
