package val

import (
	"time"

	"github.com/go-playground/validator/v10"
)

func DateTime(fl validator.FieldLevel) bool {
	if fl.Field().String() == "" {
		return true
	}
	_, err := time.Parse(time.RFC3339, fl.Field().String())
	return err == nil
}
