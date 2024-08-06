package validator

import (
	ioc "github.com/Ignaciojeria/einar-ioc"
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validator *validator.Validate
}

func (cv *Validator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func init() {
	ioc.Registry(NewValidator)
}
func NewValidator() *Validator {
	return &Validator{validator: validator.New()}
}
