package shared

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Nerzal/gocloak/v13"
	"github.com/go-playground/validator/v10"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrExtractResponse = errors.New("error while extracting response from API")
var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyExist = errors.New("user already exist")
var ErrUserUnauthorized = errors.New("error mission not permission")

type HTTPError struct {
	StatusCode int
	Message    string
}

func (r *HTTPError) Error() string {
	return fmt.Sprintf("[%d] %s", r.StatusCode, r.Message)
}

func fieldViolation(field string, err error) *errdetails.BadRequest_FieldViolation {
	return &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: err.Error(),
	}
}

func InvalidArgumentError(violations []*errdetails.BadRequest_FieldViolation) error {
	badRequest := &errdetails.BadRequest{FieldViolations: violations}

	return status.Errorf(codes.InvalidArgument, "invalid parameters: %v", badRequest)
}

func UnauthenticatedError(err error) error {
	return status.Errorf(codes.Unauthenticated, "unauthorized: %s", err)
}

func NotFoundError(err error) error {
	return status.Errorf(codes.NotFound, "not found: %s", err)
}

func ValidationErrors(err error) error {
	var valErrs validator.ValidationErrors

	var edetails []*errdetails.BadRequest_FieldViolation

	if errors.As(err, &valErrs) {
		for _, e := range valErrs {
			edetails = append(edetails, fieldViolation(e.StructField(), e))
		}
	}

	return InvalidArgumentError(edetails)
}

func CheckError(err error) error {
	var apiError *gocloak.APIError

	switch {
	case errors.As(err, &apiError):
		switch apiError.Code {
		case http.StatusNotFound:
			return NotFoundError(err)
		case http.StatusConflict:
			return status.Error(codes.AlreadyExists, err.Error())
		case http.StatusUnauthorized:
			return status.Error(codes.Unauthenticated, err.Error())
		default:
			return err
		}
	case errors.Is(err, ErrUserNotFound):
		return NotFoundError(err)
	case errors.Is(err, ErrUserUnauthorized):
		return UnauthenticatedError(err)
	default:
		return err
	}
}
