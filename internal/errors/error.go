package errors

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
