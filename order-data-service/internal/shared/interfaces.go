package shared

type Validator interface {
	ValidateStruct(s interface{}) error
}
