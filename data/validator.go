package data

import (
	"github.com/go-playground/validator/v10"
	"strings"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func Validate(s interface{}) error {
	return validate.Struct(s)
}

func ValidateUser(user string) error {
	user = strings.ToLower(user)

	return validate.Var(user, "oneof=carlitos chio tecnologer")
}
