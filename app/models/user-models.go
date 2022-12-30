package models

type LoginRequest struct {
	Login    *string `json:"login" validate:"required,alphanum,gte=3,lte=50"`
	Password *string `json:"password" validate:"required,lte=50"`
}

type RegisterRequest struct {
	LoginRequest
	Password2 *string `json:"password2" validate:"required"`
}
