package constants

var (
	DateFormat        string = "dd-mm-yyyy"
	InvalidRoomId     string = "invalid room id"
	RoomAlreadyBooked string = "room is already booked"
	InsuffientBalance string = "insuffient balance to book the room"
	StartDateError    string = "start date should not be less end date"
	EndDateError      string = "end date cannot be less than start date"
)

func InsertErrorMessage(entity string) string {
	return "error while inserting - " + entity
}

func UpdateErrorMessage(entity string) string {
	return "error while updating - " + entity
}
