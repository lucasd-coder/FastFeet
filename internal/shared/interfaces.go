package shared

import "context"

type Validator interface {
	ValidateStruct(s interface{}) error
}

type Publish interface {
	Send(ctx context.Context, msg *Message) error
}
