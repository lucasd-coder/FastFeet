package val

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func ObjectID(fl validator.FieldLevel) bool {
	if fl.Field().String() == "" {
		return true
	}

	objectIDRegex := regexp.MustCompile(`^[a-f\d]{24}$`)
	value := fl.Field().String()
	return objectIDRegex.MatchString(value)
}
