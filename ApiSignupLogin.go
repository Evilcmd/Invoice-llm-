package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Evilcmd/Invoice-llm-/internal/database"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userDefn struct {
	Uname  string `json:"uname"`
	Passwd string `json:"passwd"`
}

func (apiConfig *apiConfigDefn) signup(res http.ResponseWriter, req *http.Request) {

	user := userDefn{}
	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		respondWithError(res, 406, fmt.Sprintf("error io reading request bodu: %v", err.Error()))
		return
	}

	err = json.Unmarshal(reqBody, &user)
	if err != nil {
		respondWithError(res, 406, fmt.Sprintf("error unmarshalling user login data: %v", err.Error()))
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Passwd), 12)
	if err != nil {
		respondWithError(res, 406, fmt.Sprintf("error generating hash: %v", err.Error()))
		return
	}

	params := database.CreateUserParams{
		ID:       uuid.New(),
		UserName: user.Uname,
		Passwd:   string(hash),
	}

	err = apiConfig.DB.CreateUser(context.Background(), params)
	if err != nil {
		respondWithError(res, 406, fmt.Sprintf("error creating user: %v", err.Error()))
		return
	}

	respondWithJson(res, 200, "success")

}

func (apiConfig *apiConfigDefn) login(res http.ResponseWriter, req *http.Request) {

	user := userDefn{}
	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		respondWithError(res, 406, fmt.Sprintf("error io reading request bodu: %v", err.Error()))
		return
	}

	err = json.Unmarshal(reqBody, &user)
	if err != nil {
		respondWithError(res, 406, fmt.Sprintf("error unmarshalling user login data: %v", err.Error()))
		return
	}

	userRetreived, err := apiConfig.DB.GetUserByUsername(context.Background(), user.Uname)
	if err != nil {
		respondWithError(res, 406, fmt.Sprintf("error fetching admin: %v", err.Error()))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(userRetreived.Passwd), []byte(user.Passwd))
	if err != nil {
		respondWithError(res, 406, "wrong password")
		return
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "InvoiceLLM",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		Subject:   userRetreived.ID.String(),
	})

	signedJwtTokenString, err := jwtToken.SignedString([]byte(apiConfig.jwtSecret))
	if err != nil {
		respondWithError(res, 406, fmt.Sprintf("error signing the token: %v", err.Error()))
		return
	}

	respondWithJson(res, 200, signedJwtTokenString)
}
