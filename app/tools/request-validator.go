package tools

import (
	"errors"

	"github.com/go-playground/validator"
)

type Validator struct {
	Validator *validator.Validate
}

func (cv *Validator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		switch err.Error() {
		case "Key: 'RegisterRequest.LoginRequest.Login' Error:Field validation for 'Login' failed on the 'lte' tag":
			return errors.New("login is too long")
		case "Key: 'RegisterRequest.LoginRequest.Login' Error:Field validation for 'Login' failed on the 'gte' tag":
			return errors.New("login is too short")
		case "Key: 'RegisterRequest.LoginRequest.Login' Error:Field validation for 'Login' failed on the 'alphanum' tag":
			return errors.New("login should only contain letters and numbers")
		case "Key: 'RegisterRequest.LoginRequest.Password' Error:Field validation for 'Password' failed on the 'lte' tag":
			return errors.New("password is too long")
		default:
			return err
		}
	}
	return nil
}

func NewValidator() *Validator {
	return &Validator{Validator: validator.New()}
}
