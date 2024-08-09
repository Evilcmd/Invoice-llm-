package main

import (
	"encoding/json"
	"net/http"
)

func respondWithJson(res http.ResponseWriter, statusCode int, payload interface{}) {
	dat, _ := json.Marshal(payload)
	res.WriteHeader(statusCode)
	res.Write(dat)
}

func respondWithError(res http.ResponseWriter, statusCode int, message string) {
	payload := struct {
		Error string `json:"error"`
	}{
		message,
	}
	respondWithJson(res, statusCode, payload)
}

func checkHealth(res http.ResponseWriter, req *http.Request) {
	payload := struct {
		Message string `json:"message"`
	}{
		"OK",
	}
	respondWithJson(res, 200, payload)
}

func errCheck(res http.ResponseWriter, req *http.Request) {
	respondWithError(res, 400, "checking error response function")
}
