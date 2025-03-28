package types

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Hotel struct {
	ID              interface{}     `bson:"_id,omitempty"`
	Name            string          `json:"name" bson:"name"`
	Description     string          `json:"description" bson:"description"`
	Star            int8            `json:"star" bson:"star"`
	Review          []Review        `json:"review" bson:"review"`
	Address         Address         `json:"address" bson:"address"`
	Amenities       any             `json:"amenities" bson:"amenities"`
	AdditionalInfo1 string          `json:"additionalinfo1,omitempty" bson:"additionalinfo1,omitempty"`
	AdditionalInfo2 string          `json:"additionalinfo2,omitempty" bson:"additionalinfo2,omitempty"`
	AdditionalInfo3 any             `json:"additionalinfo3,omitempty" bson:"additionalinfo3,omitempty"`
	TypesOfRooms    []string        `json:"typesofrooms,omitempty" bson:"typesofrooms,omitempty"`
	Rooms           []bson.ObjectID `json:"rooms,omitempty" bson:"rooms,omitempty"`
	Created         time.Time       `json:"created_at" bson:"created_at"`
	Updated         time.Time       `json:"updated_at" bson:"updated_at"`
}
