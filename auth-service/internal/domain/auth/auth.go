package auth

import "github.com/lucasd-coder/fast-feet/auth-service/internal/shared"

type RolesEnum string

const (
	Unknown RolesEnum = "UNKNOWN"
	Admin   RolesEnum = "ADMIN"
	User    RolesEnum = "USER"
)

type Register struct {
	FirstName string `json:"firstName,omitempty" validate:"pattern"`
	LastName  string `json:"lastName,omitempty" validate:"pattern"`
	Username  string `json:"username,omitempty" validate:"required,email"`
	Password  string `json:"password,omitempty" validate:"min=8,containsany=!@#?*"`
	Roles     string `json:"roles,omitempty" validate:"required,oneof=ADMIN USER"`
}

type Response struct {
	ID string `json:"id,omitempty"`
}

func (r *Register) Validate(val shared.Validator) error {
	return val.ValidateStruct(r)
}

var RolesEnumString = map[RolesEnum]string{
	Admin:   "admin",
	User:    "user",
	Unknown: "unknown",
}

func (r RolesEnum) String() string {
	return RolesEnumString[r]
}

func GetRolesFromString(s string) RolesEnum {
	switch s {
	case "ADMIN":
		return Admin
	case "USER":
		return User
	default:
		return Unknown
	}
}

type FindUserByEmail struct {
	Email string `json:"email" validate:"required,email"`
}

func (f *FindUserByEmail) Validate(val shared.Validator) error {
	return val.ValidateStruct(f)
}

type UserRepresentation struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Enabled  bool   `json:"enabled,omitempty"`
	Email    string `json:"email,omitempty"`
}

type GetUserID struct {
	ID string `json:"id" validate:"required,uuid"`
}

func (g *GetUserID) Validate(val shared.Validator) error {
	return val.ValidateStruct(g)
}
