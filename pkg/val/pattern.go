package val

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func Pattern(fl validator.FieldLevel) bool {
	pattern := regexp.MustCompile(`^[a-zA-Z0-9_áàâãéèêíïóôõöúçñÁÀÂÃÉÈÍÏÓÔÕÖÚÇÑ:\\/@#*,.?!+\-\s]*$`)
	return pattern.MatchString(fl.Field().String())
}
