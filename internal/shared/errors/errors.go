package errors

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

var ErrCipherText = errors.New("cipher text too short")

type fieldError struct {
	err validator.FieldError
}

type FieldMessage struct {
	FieldName string `json:"fieldName,omitempty"`
	Message   string `json:"message,omitempty"`
}

type StandardError struct {
	Timestamp  string         `json:"timestamp,omitempty"`
	StatusCode int            `json:"statusCode,omitempty"`
	Message    string         `json:"message,omitempty"`
	Errors     []FieldMessage `json:"errors,omitempty"`
}

func (s *StandardError) Error() string {
	return s.Message
}

func NewStandardError(msg string, code int) StandardError {
	return StandardError{
		Timestamp:  time.Now().Format(time.RFC3339),
		StatusCode: code,
		Message:    msg,
		Errors:     []FieldMessage{},
	}
}

func (s *StandardError) AddError(fieldName, message string) {
	fieldMessage := FieldMessage{fieldName, message}
	s.Errors = append(s.Errors, fieldMessage)
}

func BuildError(err error) StandardError {
	var ve validator.ValidationErrors
	var errResp StandardError

	if errors.As(err, &ve) {
		errResp = NewStandardError("Validation Error", http.StatusUnprocessableEntity)
		for _, e := range ve {
			errResp.AddError(e.StructField(), fieldError{err: e}.String())
		}
	} else {
		errResp = NewStandardError(err.Error(), http.StatusInternalServerError)
	}

	return errResp
}

func (q fieldError) String() string {
	var sb strings.Builder

	sb.WriteString("validation failed on field '" + q.err.Field() + "'")
	sb.WriteString(", condition: " + q.err.ActualTag())

	// Print condition parameters, e.g. oneof=red blue -> { red blue }
	if q.err.Param() != "" {
		sb.WriteString(" { " + q.err.Param() + " }")
	}

	if q.err.Value() != nil && q.err.Value() != "" {
		sb.WriteString(fmt.Sprintf(", actual: %v", q.err.Value()))
	}

	return sb.String()
}
