package db

import (
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func DB_Connection(MONGODB_URI string) *mongo.Client {
	if MONGODB_URI == "" {
		panic("Missing/Invalid Mongo Uri")
	}
	client, err := mongo.Connect(options.Client().ApplyURI(MONGODB_URI))
	if err != nil {
		connectionError := fmt.Errorf("Missing/Invalid Mongo Uri - %s", err.Error())
		panic(connectionError)
	}
	return client
}
