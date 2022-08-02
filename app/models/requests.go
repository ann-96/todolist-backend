package models

type ListRequest struct {
	Start string `json:"start" form:"start" query:"start"`
	Count string `json:"count" form:"count" query:"count"`
}
