package validator

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
}

func NewValidator() echo.Validator {
	return &CustomValidator{validator: validator.New()}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	cv.validator.RegisterValidation("is_int", IsInt)
	return cv.validator.Struct(i)
}

func IsInt(fl validator.FieldLevel) bool {
	_, err := strconv.Atoi(fl.Field().String())
	if err != nil {
		return false
	}
	return true
}
