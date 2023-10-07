package val

import (
	"regexp"
	"strconv"

	"github.com/go-playground/validator/v10"
)

var cpfRegexp = regexp.MustCompile(`^\d{3}\.?\d{3}\.?\d{3}-?\d{2}$`)

// IsCPF verifies if the given string is a valid CPF document.
func IsCPF(doc string) bool {
	return validate(doc, cpfRegexp, calculateCPFVerifierDigits)
}

func TagIsCPF(fl validator.FieldLevel) bool {
	return IsCPF(fl.Field().String())
}

func validate(doc string, regexp *regexp.Regexp, calculateVerifierDigits func(string) string) bool {
	if !regexp.MatchString(doc) {
		return false
	}

	doc = sanitize(doc)

	// Invalidates documents with all digits equal.
	if allEq(doc) {
		return false
	}

	verifierDigits := calculateVerifierDigits(doc[:len(doc)-2])
	return doc == doc[:len(doc)-2]+verifierDigits
}

func sanitize(data string) string {
	data = regexp.MustCompile(`[^\d]`).ReplaceAllString(data, "")
	return data
}

// allEq checks if every rune in a given string is equal.
func allEq(doc string) bool {
	base := doc[0]
	for i := 1; i < len(doc); i++ {
		if base != doc[i] {
			return false
		}
	}

	return true
}

// calculateCPFVerifierDigits calculates the verifier digits for a given CPF document.
func calculateCPFVerifierDigits(doc string) string {
	weights := []int{10, 9, 8, 7, 6, 5, 4, 3, 2}
	firstDigit := calculateVerifierDigit(doc, weights)
	weights = append([]int{11}, weights...)
	secondDigit := calculateVerifierDigit(doc+firstDigit, weights)
	return firstDigit + secondDigit
}

func calculateVerifierDigit(doc string, weights []int) string {
	var sum int
	modulo11 := 11
	for i, r := range doc {
		sum += int(r-'0') * weights[i]
	}
	remainder := sum % modulo11
	if remainder < 2 {
		return "0"
	}
	return strconv.Itoa(modulo11 - remainder)
}
