package main

import "net/http"

func rootFunc(res http.ResponseWriter, req *http.Request) {
	respondWithJson(res, 200, "Hello User")
}
