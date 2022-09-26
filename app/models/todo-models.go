package models

type AddTodoRequest struct {
	Text      *string `json:"text" validate:"required"`
	Completed *bool   `json:"completed" validate:"required"`
}

type Todo struct {
	AddTodoRequest
	Id int `json:"id" validate:"required"`
}

type TodoList struct {
	List           []Todo `json:"list" validate:"required"`
	Count          int    `json:"count" validate:"required"`
	CompletedCount int    `json:"completedCount" validate:"required"`
}

type ListRequest struct {
	Start string `json:"start" form:"start" query:"start"`
	Count string `json:"count" form:"count" query:"count"`
}
