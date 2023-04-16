package shared

import (
	"errors"
	"fmt"
)

var ErrExtractResponse = errors.New("error while extracting response from API")
var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyExist = errors.New("user already exist")

type HTTPError struct {
	StatusCode int
	Message    string
}

func (r *HTTPError) Error() string {
	return fmt.Sprintf("[%d] %s", r.StatusCode, r.Message)
}
