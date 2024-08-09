package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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

func readPdfPromptLLM(fname string) (details, error) {

	//Pdf
	args := []string{
		"-layout",
		"-nopgbrk",
		fname,
		"-",
	}
	cmd := exec.Command("pdftotext", args...)
	output, err := cmd.Output()
	if err != nil {
		return details{}, err
	}

	invoice := string(output)

	// Gemini Config
	ctx := context.Background()
	godotenv.Load()
	GeminiApiKey := os.Getenv("GEMINI_API_KEY")
	client, err := genai.NewClient(ctx, option.WithAPIKey(GeminiApiKey))
	if err != nil {
		return details{}, err
	}
	defer client.Close()
	model := client.GenerativeModel("gemini-1.5-flash")

	// Prompt
	resp, err := model.GenerateContent(ctx, genai.Text(invoice+"\nFrom the above invoice, return the customer_details(name, address, phone_number, email), array of product_details(name, rate, quantity, total_amount), total_amount and amount_payable in Json format if any of the fields mentioned doesnt exist then mention 0 in its place and give the numbers without commas and no need to mention it is json and dont add any backticks at the beginning and ending"))
	if err != nil {
		return details{}, err
	}

	response := resp.Candidates[0].Content.Parts[0]
	invoiceData, err := partToString(response)
	if err != nil {
		return details{}, err
	}

	detailsobj := &details{}

	err = json.Unmarshal([]byte(invoiceData), detailsobj)
	if err != nil {
		return details{}, err
	}

	return *detailsobj, nil
}

func uploadFile(res http.ResponseWriter, req *http.Request) {
	req.ParseMultipartForm(10 << 20)
	file, _, err := req.FormFile("myFile")
	if err != nil {
		respondWithError(res, http.StatusUnprocessableEntity, fmt.Sprintf("Error Retrieving the File: %v", err.Error()))
		return
	}
	defer file.Close()

	tempFile, err := os.CreateTemp("temp", "upload-*.pdf")

	if err != nil {
		respondWithError(res, http.StatusUnprocessableEntity, fmt.Sprintf("Error Creating temp file: %v", err.Error()))
		return
	}
	defer tempFile.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		respondWithError(res, http.StatusUnprocessableEntity, fmt.Sprintf("Error reading the File: %v", err.Error()))
		return
	}

	tempFile.Write(fileBytes)

	invoice, err := readPdfPromptLLM(tempFile.Name())
	if err != nil {
		respondWithError(res, 406, fmt.Sprintf("error calling the function: %v", err.Error()))
		return
	}

	respondWithJson(res, 200, invoice)

	err = os.Remove(tempFile.Name())
	if err != nil {
		log.Println(err)
	}
}

func main() {

	router := http.NewServeMux()

	router.HandleFunc("POST /upload", uploadFile)

	server := http.Server{
		Addr:    ":8080",
		Handler: router}

	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Error starting server: ", err.Error())
	}
}
