package db

import "github.com/ann-96/todo-go-backend/app/models"

type TodoSqlDB interface {
	Update(input *models.Todo, userId int) (*models.Todo, error)
	List(start int, count int, userId int) (*models.TodoList, error)
	Add(input *models.AddTodoRequest, userId int) (*models.Todo, error)
	Delete(id int, userId int) error

	Register(input *models.RegisterRequest) error
	Login(input *models.LoginRequest) (*int, error)

	Migrate() error
}
