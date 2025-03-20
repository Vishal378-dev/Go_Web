package types

import "go.mongodb.org/mongo-driver/v2/bson"

type Room struct {
	Class      string        `json:"class" bson:"class"`
	RoomNumber uint16        `json:"roomnumber" bson:"roomnumber"`
	IsBooked   bool          `json:"isbooked" bson:"isbooked"`
	Price      float32       `json:"price" bson:"price"`
	Features   any           `json:"feature" bson:"feature"`
	HotelID    bson.ObjectID `json:"hotelid" bson:"hotelid"`
}
