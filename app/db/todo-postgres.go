package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/ann-96/todo-go-backend/app/models"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Settings struct {
	IP       string
	Port     string
	User     string
	Password string
	Name     string
}

type postgresDB struct {
	sql *sql.DB
}

func CreatePostgresDB(settings Settings) (*postgresDB, error) {
	res := &postgresDB{}
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", settings.IP, settings.Port, settings.User, settings.Password, settings.Name)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return nil, err
	}

	res.sql = db
	return res, nil
}

func (db *postgresDB) Migrate() error {
	driver, err := postgres.WithInstance(db.sql, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://./migrations", "postgres", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func (db *postgresDB) Update(input *models.Todo) (*models.Todo, error) {
	updateStmt := "UPDATE todos SET task=$1, completed=$2 WHERE id=$3 RETURNING id;"
	row := db.sql.QueryRow(updateStmt, input.Text, input.Completed, input.Id)
	var id int
	if err := row.Scan(&id); err != nil {
		return nil, err
	}

	return input, nil
}

func (db *postgresDB) List() (*models.TodoList, error) { // TODO: pagination(both on front and back end)
	rows, err := db.sql.Query("SELECT id, task, completed FROM todos;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := models.TodoList{}
	for rows.Next() {
		var id int
		var text string
		var completed bool
		if err := rows.Scan(&id, &text, &completed); err != nil {
			return nil, err
		}

		todo := models.Todo{
			Id: id,
			AddTodoRequest: models.AddTodoRequest{
				Text:      &text,
				Completed: &completed,
			},
		}

		list = append(list, todo)
	}

	return &list, nil
}

func (db *postgresDB) Add(input *models.AddTodoRequest) (*models.Todo, error) {
	stmnt, err := db.sql.Prepare("INSERT INTO todos(task, completed) VALUES ($1, $2) RETURNING id;")
	if err != nil {
		return nil, err
	}

	res := &models.Todo{
		AddTodoRequest: *input,
	}

	row := stmnt.QueryRow(input.Text, input.Completed)
	if err := row.Scan(&res.Id); err != nil {
		return nil, err
	}

	return res, nil
}

func (db *postgresDB) Delete(id int) error {
	updateStmt := "DELETE FROM todos WHERE id=$1;"
	_, err := db.sql.Exec(updateStmt, id)
	if err != nil {
		return err
	}

	return nil
}
