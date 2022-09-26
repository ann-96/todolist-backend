package tools

import (
	"net/http"

	"github.com/go-playground/validator"
	echo "github.com/labstack/echo/v4"
)

type Validator struct {
	Validator *validator.Validate
}

func (cv *Validator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func NewValidator() *Validator {
	return &Validator{Validator: validator.New()}
}
