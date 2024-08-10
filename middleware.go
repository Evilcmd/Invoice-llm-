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
