package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func partToString(part genai.Part) (string, error) {
	switch v := part.(type) {
	case genai.Text:
		return string(v), nil
	default:
		return "", fmt.Errorf("error in type")
	}
}

func main() {

	//Pdf
	args := []string{
		"-layout",
		"-nopgbrk",
		"Sample Invoice2.pdf",
		"-",
	}
	cmd := exec.Command("pdftotext", args...)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	invoice := string(output)

	// Gemini Config
	ctx := context.Background()
	godotenv.Load()
	GeminiApiKey := os.Getenv("GEMINI_API_KEY")
	client, err := genai.NewClient(ctx, option.WithAPIKey(GeminiApiKey))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	model := client.GenerativeModel("gemini-1.5-flash")

	// Prompt
	resp, err := model.GenerateContent(ctx, genai.Text(invoice+"\nFrom the above invoice, return the customer_details(name, address, phone_number, email), array of product_details(name, rate, quantity, total_amount), total_amount and amount_payable in Json format if any of the fields mentioned doesnt exist then mention 0 in its place and give the numbers without commas and no need to mention it is json"))
	if err != nil {
		log.Fatal(err)
	}

	response := resp.Candidates[0].Content.Parts[0]
	invoiceData, err := partToString(response)
	if err != nil {
		log.Fatal("Error in converting part to string", err.Error())
	}

	fmt.Println(invoiceData)

	type details struct {
		CustomerDetails struct {
			Name        string `json:"name"`
			Address     string `json:"address"`
			PhoneNumber string `json:"phone_number"`
			Email       string `json:"email"`
		} `json:"customer_details"`
		ProductDetails []struct {
			Name        string `json:"name"`
			Rate        string `json:"rate"`
			Quantity    string `json:"quantity"`
			TotalAmount string `json:"total_amount"`
		} `json:"product_details"`
		TotalAmount   string `json:"total_amount"`
		AmountPayable string `json:"amount_payable"`
	}
	detailsobj := &details{}

	err = json.Unmarshal([]byte(invoiceData), detailsobj)
	if err != nil {
		log.Fatal("Error unmarshalling: ", err.Error())
	}

	fmt.Println(*detailsobj)

}
