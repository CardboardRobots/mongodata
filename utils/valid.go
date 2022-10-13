package utils

import "github.com/go-playground/validator"

var validate = validator.New()

type Valid interface {
	Valid() error
}

func StructValidate(data any) error {
	return validate.Struct(data)
}
