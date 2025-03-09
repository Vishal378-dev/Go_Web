package types

type ErrorResponse struct {
	Error  string
	Status int
}

type SuccessResponse struct {
	Data   interface{}
	Status int
}
