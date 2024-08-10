package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type InvoiceWithUserID struct {
	UserId  string
	Invoice details
}

func (apiConfig *apiConfigDefn) saveInvoice(res http.ResponseWriter, req *http.Request) {
	invoice := details{}
	invoiceReceived, err := io.ReadAll(req.Body)
	if err != nil {
		respondWithError(res, 406, fmt.Sprintf("error reading request body: %v", err.Error()))
		return
	}
	err = json.Unmarshal(invoiceReceived, &invoice)
	if err != nil {
		respondWithError(res, 406, fmt.Sprintf("error unmarshalling: %v", err.Error()))
		return
	}

	mongoObject := InvoiceWithUserID{apiConfig.userID.String(), invoice}

	result, err := apiConfig.MongoDBCLient.InsertOne(context.TODO(), mongoObject)
	if err != nil {
		respondWithError(res, 406, fmt.Sprintf("error inserting into database: %v", err.Error()))
		return
	}

	respondWithJson(res, 200, result)
}
