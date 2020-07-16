package service

import (
	"errors"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func SetupValidator() error {
	validate := binding.Validator.Engine()
	v, ok := validate.(*validator.Validate)
	if !ok {
		return errors.New("setup validator error")
	}

	// string length >= 10
	return v.RegisterValidation("my-val", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) >= 10
	})
}
