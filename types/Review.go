package types

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Review struct {
	Rating  float32       `json:"rating" bson:"rating"`
	Comment string        `json:"comment" bson:"comment"`
	UserId  bson.ObjectID `json:"userid" bson:"userid"`
	Created time.Time     `json:"created_at" bson:"created_at"`
	Updated time.Time     `json:"updated_at" bson:"updated_at"`
}
