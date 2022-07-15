package models

type AddTodoRequest struct {
	Text      *string `json:"text" validate:"required"`
	Completed *bool   `json:"completed" validate:"required"`
}

type Todo struct {
	AddTodoRequest
	Id int `json:"id" validate:"required"`
}

type TodoList []Todo
