package models

type LoginRequest struct {
	Login    *string `json:"login" validate:"required"`
	Password *string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	LoginRequest
	Password2 *string `json:"password2" validate:"required"`
}
