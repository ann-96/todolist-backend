package db

import "github.com/ann-96/todo-go-backend/app/models"

type TodoSqlDB interface {
	Update(input *models.Todo) (*models.Todo, error)
	List(start int, count int) (*models.TodoList, error)
	Add(input *models.AddTodoRequest) (*models.Todo, error)
	Delete(id int) error
	Migrate() error
}
