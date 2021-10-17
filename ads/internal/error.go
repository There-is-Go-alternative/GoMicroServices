package internal

import "fmt"

const (
	InvalidAdID  string = "Ad ID is invalid"
	AdNotFound   string = "Ad not found"
	MissingParam string = "Missing param"
	AuthorizationError string = "Authorization token is not in the headers"
	AdIsClose string = "Ad is closed"
	Unauthorized string = "You don't have the authorization"
	InternalServerError string = "Intertal server error"
)

const (
	NotFound = iota
	BadRequest
	DatabaseError
)

type CustomError struct {
	Code int
	Err  error
}

func NewInternalError(code int, message string) (*CustomError) {
	var err = new(CustomError);

	err.Code = code
	err.Err =  fmt.Errorf("%s", message)

	return err
}

func (err CustomError) GetCode() int {
	return (err.Code)
}

func (err CustomError) Error() string {
	return fmt.Sprintf("[%d] %s\n", err.Code, err.Err)
}
