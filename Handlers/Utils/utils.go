package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vishal/reservation_system/types"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func ResponseWriter(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func CommonError(err error, status int) types.ErrorResponse {
	return types.ErrorResponse{
		Error:  err.Error(),
		Status: status,
	}
}

func CreateEmailIndex(collection *mongo.Collection) {

	indexes, err := collection.Indexes().List(context.TODO())
	if err != nil {
		log.Fatal("Error while fetching the Index : - ", err.Error())
	}
	var indexExist bool
	for indexes.Next(context.TODO()) {
		var index bson.M
		err := indexes.Decode(&index)
		if err != nil {
			log.Fatal("Error while decoding the Index : - ", err.Error())
		}
		if index["name"] == "email_1" {
			indexExist = true
			break
		}
	}
	if !indexExist {
		fmt.Println("index running for first time")
		indexModel := mongo.IndexModel{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		}

		_, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
		if err != nil {
			log.Fatal("Error while creating index:-", err)
		} else {
			fmt.Println("Unique index on email created successfully.")
		}
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func ComparePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func NewAccessToken(claims types.UserClaims) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return accessToken.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
}

func ParseToken(tokenString string, secretKey string) (*types.UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &types.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*types.UserClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
