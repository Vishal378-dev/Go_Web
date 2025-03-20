package types

type Room struct {
	Class      string  `json:"room" bson:"room"`
	RoomNumber uint16  `json:"roomnumber" bson:"roomnumber"`
	IsBooked   bool    `json:"isbooked" bson:"isbooked"`
	Price      float32 `json:"price" bson:"price"`
	Features   any     `json:"feature" bson:"feature"`
}
