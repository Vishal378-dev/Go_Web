package types

type Coordinates struct {
	Latitude  float32 `json:"latitude,omitempty" bson:"latitude,omitempty"`
	Longitude float32 `json:"longitude,omitempty" bson:"longitude,omitempty"`
}

type Address struct {
	LandMark    string      `json:"landmark,omitempty" bson:"landmark,omitempty"`
	City        string      `json:"city,omitempty" bson:"city,omitempty"`
	State       string      `json:"state,omitempty" bson:"state,omitempty"`
	Street      string      `json:"street,omitempty" bson:"street,omitempty"`
	Pincode     int32       `json:"pincode,omitempty" bson:"pincode,omitempty"`
	Coordinates Coordinates `json:"coordinates,omitempty" bson:"coordinates,omitempty"`
}
