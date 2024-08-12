package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type userDefn struct {
	Uname  string `json:"uname"`
	Passwd string `json:"passwd"`
}

type jwtSend struct {
	Token string `json:"token"`
}

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"uname"`
	Passwd   string             `bson:"passwd"`
}

func (apiConfig *apiConfigDefn) signup(res http.ResponseWriter, req *http.Request) {

	user := userDefn{}
	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		respondWithError(res, 406, fmt.Sprintf("error io reading request body: %v", err.Error()))
		return
	}

	err = json.Unmarshal(reqBody, &user)
	if err != nil {
		respondWithError(res, 406, fmt.Sprintf("error unmarshalling user login data: %v", err.Error()))
		return
	}

	if user.Uname == "" || user.Passwd == "" {
		respondWithError(res, 406, "did not get username or password")
		return
	}

	filter := bson.D{{Key: "uname", Value: user.Uname}}
	userRet := apiConfig.MongoDBUserCLient.FindOne(context.TODO(), filter)
	err = userRet.Err()
	if err == nil || err != mongo.ErrNoDocuments {
		respondWithError(res, 406, "Duplicate user or some other error:")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Passwd), 12)
	if err != nil {
		respondWithError(res, 406, fmt.Sprintf("error generating hash: %v", err.Error()))
		return
	}

	params := User{
		ID:       primitive.NewObjectID(),
		Username: user.Uname,
		Passwd:   string(hash),
	}

	_, err = apiConfig.MongoDBUserCLient.InsertOne(context.TODO(), params)
	if err != nil {
		respondWithError(res, 406, fmt.Sprintf("error creating user: %v", err.Error()))
		return
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "InvoiceLLM",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		Subject:   params.ID.Hex(),
	})

	signedJwtTokenString, err := jwtToken.SignedString([]byte(apiConfig.jwtSecret))
	if err != nil {
		respondWithError(res, 406, fmt.Sprintf("error signing the token: %v", err.Error()))
		return
	}

	respondWithJson(res, 200, jwtSend{signedJwtTokenString})

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

	filter := bson.D{{Key: "uname", Value: user.Uname}}
	userRet := apiConfig.MongoDBUserCLient.FindOne(context.TODO(), filter)
	err = userRet.Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			respondWithError(res, 404, "user not found")
			return
		}
		respondWithError(res, 500, fmt.Sprintf("error finding user: %v", err.Error()))
		return
	}

	userRetreived := User{}
	userRet.Decode(&userRetreived)

	err = bcrypt.CompareHashAndPassword([]byte(userRetreived.Passwd), []byte(user.Passwd))
	if err != nil {
		respondWithError(res, 406, "wrong password")
		return
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "InvoiceLLM",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		Subject:   userRetreived.ID.Hex(),
	})

	signedJwtTokenString, err := jwtToken.SignedString([]byte(apiConfig.jwtSecret))
	if err != nil {
		respondWithError(res, 406, fmt.Sprintf("error signing the token: %v", err.Error()))
		return
	}

	respondWithJson(res, 200, jwtSend{signedJwtTokenString})
}
