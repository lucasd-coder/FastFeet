package val_test

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/lucasd-coder/fast-feet/business-service/pkg/val"
)

var valInst *validator.Validate

func TestPattern(t *testing.T) {
	type testStruct struct {
		Field string `validate:"pattern"`
	}

	validCases := []testStruct{
		{Field: "abcd1234"},
		{Field: "abc_def"},
		{Field: "áàâãéèêíïóôõöúçñÁÀÂÃÉÈÍÏÓÔÕÖÚÇÑ:\\/@#,.+-"},
	}

	invalidCases := []testStruct{
		{Field: "abcd!@#$"},
		{Field: "abc@&%%%$&&def"},
		{Field: "A B C"},
	}

	valInst = validator.New()

	if err := valInst.RegisterValidation("pattern", val.Pattern); err != nil {
		t.Errorf("err register validation pattern error: %v", err)
	}

	for _, c := range validCases {
		err := valInst.Struct(c)
		if err != nil {
			t.Errorf("expected %v to be valid, but got error: %v", c.Field, err)
		}
	}

	for _, c := range invalidCases {
		err := valInst.Struct(c)
		if err == nil {
			t.Errorf("expected %v to be invalid, but got no error", c.Field)
		}
	}
}
