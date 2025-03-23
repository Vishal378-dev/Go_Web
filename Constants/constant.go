package constants

var (
	DateFormat    string = "dd-mm-yyyy"
	InvalidRoomId string = "invalid room id"
)

func InsertErrorMessage(entity string) string {
	return "error while inserting - " + entity
}

func UpdateErrorMessage(entity string) string {
	return "error while updating - " + entity
}
