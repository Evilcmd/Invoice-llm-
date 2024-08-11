package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (apiConfig *apiConfigDefn) authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		signedJwtTokenString := req.Header.Get("Authorization")
		if len(signedJwtTokenString) == 0 {
			respondWithError(res, 406, "expected authoriztion")
			return
		}
		signedJwtTokenString = strings.Split(signedJwtTokenString, " ")[1]

		jwtToken, err := jwt.ParseWithClaims(signedJwtTokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(apiConfig.jwtSecret), nil
		})
		if err != nil {
			respondWithError(res, http.StatusUnauthorized, fmt.Sprintf("Unauthorized: %v", err.Error()))
			return
		}

		userid, err := jwtToken.Claims.GetSubject()
		if err != nil {
			respondWithError(res, http.StatusUnauthorized, fmt.Sprintf("Error getting id from jwtToken: %v", err.Error()))
			return
		}

		userUUID, err := uuid.Parse(userid)
		if err != nil {
			respondWithError(res, http.StatusUnauthorized, fmt.Sprintf("Error parsing id from jwtToken: %v", err.Error()))
			return
		}

		apiConfig.userID = userUUID

		// skipping refresh tokens for now

		next.ServeHTTP(res, req)
	}
}

// CORS Middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace '*' with specific origins if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

