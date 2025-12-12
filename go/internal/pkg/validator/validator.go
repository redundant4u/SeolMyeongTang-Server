package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

type EchoValidator struct {
	v *validator.Validate
}

func New() *EchoValidator {
	v := validator.New()
	v.RegisterValidation("k8slabel", validateK8sLabel)

	return &EchoValidator{v: v}
}

func (e *EchoValidator) Validate(i interface{}) error {
	return e.v.Struct(i)
}

func validateK8sLabel(fl validator.FieldLevel) bool {
	label := fl.Field().String()
	pattern := `^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$`

	reg := regexp.MustCompile(pattern)
	result := reg.MatchString(label)

	return result
}
