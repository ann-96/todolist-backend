package models

const (
	ResponseBadRequest = "bad request"
)

type ErrResponse struct {
	Msg string `json:"msg"`
}
