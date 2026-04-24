package validator

import (
	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator *validator.Validate
}

func New() *CustomValidator {
	return &CustomValidator{validator: validator.New()}
}

func (cv *CustomValidator) Validate(i any) error {
	return cv.validator.Struct(i)
}
