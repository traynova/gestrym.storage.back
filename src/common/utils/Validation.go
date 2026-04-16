package utils

import (
	"regexp"
	"sync"

	"github.com/go-playground/validator/v10"
)

type ValidatorUtil interface {
	New() *validator.Validate
}

type validatorUtil struct {
	validator *validator.Validate
}

var (
	validatorInstance *validatorUtil
	validatorOnce     sync.Once
)

func GetValidator() ValidatorUtil {
	validatorOnce.Do(func() {
		v := validator.New()
		registerCustomValidations(v)
		validatorInstance = &validatorUtil{
			validator: v,
		}
	})
	return validatorInstance
}

func (v *validatorUtil) New() *validator.Validate {
	return v.validator
}

func registerCustomValidations(validate *validator.Validate) {
	_ = validate.RegisterValidation("regexp", func(fl validator.FieldLevel) bool {
		pattern := fl.Param()
		if pattern == "" {
			return false
		}
		value := fl.Field().String()
		matched, err := regexp.MatchString(pattern, value)
		return err == nil && matched
	})
}
