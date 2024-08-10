package main

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func (apiConfig *apiConfigDefn) getAllInvoiceForUser(res http.ResponseWriter, req *http.Request) {
	filter := bson.D{{Key: "userid", Value: apiConfig.userID.String()}}
	invoices, err := apiConfig.MongoDBCLient.Find(context.TODO(), filter)
	if err != nil {
		respondWithError(res, 406, fmt.Sprintf("error fetching invoices: %v", err.Error()))
		return
	}
	defer invoices.Close(context.TODO())

	invoiceSlice := make([]InvoiceWithUserID, 0)

	err = invoices.All(context.TODO(), &invoiceSlice)
	if err != nil {
		respondWithError(res, 406, fmt.Sprintf("error converting cursor to slice of structs: %v", err.Error()))
		return
	}

	respondWithJson(res, 200, invoiceSlice)
}
