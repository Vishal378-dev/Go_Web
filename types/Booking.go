package types

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Booking struct {
	UserId     bson.ObjectID `json:"userid" bson:"userid"`
	RoomId     bson.ObjectID `json:"roomid" bson:"roomid"`
	StartDate  time.Time     `json:"startdate" bson:"startdate"`
	EndDate    time.Time     `json:"enddate" bson:"enddate"`
	AmountPaid float32       `json:"amountpaid" bson:"amountpaid"`
}

type BookingRequest struct {
	RoomId    bson.ObjectID `json:"roomid" bson:"roomid"`
	StartDate string        `json:"startdate" bson:"startdate"`
	EndDate   string        `json:"enddate" bson:"enddate"`
}
