package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type apiConfigDefn struct {
	MongoDBUserCLient *mongo.Collection
	MongoDBCLient     *mongo.Collection
	jwtSecret         string
	userID            primitive.ObjectID
}

func main() {

	godotenv.Load()

	JWT_SECRET := os.Getenv("JWT_SECRET")

	MongoUrl := os.Getenv("MONGOURL")
	MongoClientDriver, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(MongoUrl))
	if err != nil {
		log.Fatal("Mongo Connectin Failed", err.Error())
	}
	defer func() {
		if err := MongoClientDriver.Disconnect(context.TODO()); err != nil {
			log.Fatal("error disconnecting mongodb: ", err.Error())
		}
	}()

	MongoClient := MongoClientDriver.Database("InvoiceLLM").Collection("InvoiceLLM")
	MongoUserClient := MongoClientDriver.Database("InvoiceLLM").Collection("Users")

	apiConfig := apiConfigDefn{MongoUserClient, MongoClient, JWT_SECRET, primitive.NilObjectID}

	router := http.NewServeMux()

	router.HandleFunc("GET /", rootFunc)

	router.HandleFunc("GET /health", checkHealth)
	router.HandleFunc("GET /err", errCheck)

	router.HandleFunc("POST /upload", uploadFile)

	router.HandleFunc("POST /signup", apiConfig.signup)
	router.HandleFunc("POST /login", apiConfig.login)

	router.HandleFunc("POST /invoices", apiConfig.authenticate(apiConfig.saveInvoice))
	router.HandleFunc("GET /invoices", apiConfig.authenticate(apiConfig.getAllInvoiceForUser))

	server := http.Server{
		Addr:    ":8080",
		Handler: corsMiddleware(router),
	}
	fmt.Println("Starting server on port 8080")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Error starting server: ", err.Error())
	}
}
