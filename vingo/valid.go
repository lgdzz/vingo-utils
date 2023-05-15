package vingo

import (
	"github.com/go-playground/validator/v10"
)

var Valid *validator.Validate

func InitValidateService() {
	Valid = validator.New()
}
