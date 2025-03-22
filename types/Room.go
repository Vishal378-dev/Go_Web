package types

import (
	"fmt"
	"slices"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Room struct {
	Class        string        `json:"class" bson:"class"`
	RoomNumber   uint16        `json:"roomnumber" bson:"roomnumber"`
	IsBooked     *bool         `json:"isbooked" bson:"isbooked"`
	Price        float32       `json:"price" bson:"price"`
	Features     any           `json:"feature" bson:"feature"`
	HotelID      bson.ObjectID `json:"hotelid" bson:"hotelid"`
	RoomCategory string        `json:"roomcategory" bson:"roomcategory"`
}

var (
	class    = []string{"deluxe", "suite", "standard", "budget"}
	category = []string{"single", "double", "triple"}
)

func (r *Room) RequestValidation() error {
	fmt.Println(r.Class)
	if !slices.Contains(class, strings.ToLower(r.Class)) {
		return fmt.Errorf("invalid room class. only availabe room classes are - %v", class)
	}
	if !slices.Contains(category, strings.ToLower(r.RoomCategory)) {
		return fmt.Errorf("invalid roomcategory. only availabe room categories are - %v", category)
	}
	if r.IsBooked == nil {
		return fmt.Errorf("please provide the booking status")
	}

	return nil
}
