package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Evilcmd/Invoice-llm-/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type apiConfigDefn struct {
	DB            *database.Queries
	MongoDBCLient *mongo.Collection
	jwtSecret     string
	userID        uuid.UUID
}

func main() {

	godotenv.Load()
	dbURL := os.Getenv("DBURL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

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

	dbQueries := database.New(db)
	apiConfig := apiConfigDefn{dbQueries, MongoClient, JWT_SECRET, uuid.UUID{}}

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
