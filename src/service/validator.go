package service

import (
	"errors"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func SetupValidator() error {
	// 获取 validator 句柄
	validate := binding.Validator.Engine()
	v, ok := validate.(*validator.Validate)
	if !ok {
		return errors.New("setup validator error")
	}

	// 注册一个 string length >= 10 的规则
	err := v.RegisterValidation("my-val", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) >= 10
	})
	if err != nil {
		return err
	}

	// 注册 tag 为 form，
	// Password int `form:"pwd"`
	// 这里意思类似json 是指将表单中从传入的 pwd 值赋值到 Password 属性中。
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		tag := strings.SplitN(field.Tag.Get("form"), ",", 2)[0]
		if tag == "-" {
			return ""
		}
		return tag
	})

	return nil
}
