package types

import "go.mongodb.org/mongo-driver/v2/bson"

type Review struct {
	Rating  float32       `json:"rating" bson:"rating"`
	Comment string        `json:"comment" bson:"comment"`
	UserId  bson.ObjectID `json:"userid" bson:"userid"`
}
