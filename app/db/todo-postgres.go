package db

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
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

func (db *postgresDB) Update(input *models.Todo, userId int) (*models.Todo, error) {
	updateStmt := "UPDATE todos SET task=$1, completed=$2 WHERE id=$3 AND userid=$4 RETURNING id;"
	row := db.sql.QueryRow(updateStmt, input.Text, input.Completed, input.Id, userId)
	var id int
	if err := row.Scan(&id); err != nil {
		return nil, err
	}

	return input, nil
}

func (db *postgresDB) List(start int, count int, userId int) (*models.TodoList, error) {
	stm := `
		SELECT id, task, completed FROM todos 
		WHERE userid=$3
		ORDER BY id ASC 
		LIMIT $1 
		OFFSET $2;
		`
	stmt, err := db.sql.Prepare(stm)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(count, start, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []models.Todo{}
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

	row := db.sql.QueryRow("SELECT COUNT(*) FROM todos WHERE userid=$1;", userId)
	var totalCount int
	if err := row.Scan(&totalCount); err != nil {
		return nil, err
	}

	row = db.sql.QueryRow("SELECT COUNT(*) FROM todos WHERE completed = true AND userid=$1;", userId)
	var completedCount int
	if err := row.Scan(&completedCount); err != nil {
		return nil, err
	}

	return &models.TodoList{
		List:           list,
		Count:          totalCount,
		CompletedCount: completedCount,
	}, nil
}

func (db *postgresDB) Add(input *models.AddTodoRequest, userId int) (*models.Todo, error) {
	stmnt, err := db.sql.Prepare("INSERT INTO todos(task, completed, userid) VALUES ($1, $2, $3) RETURNING id;")
	if err != nil {
		return nil, err
	}

	res := &models.Todo{
		AddTodoRequest: *input,
	}

	row := stmnt.QueryRow(input.Text, input.Completed, userId)
	if err := row.Scan(&res.Id); err != nil {
		return nil, err
	}

	return res, nil
}

func (db *postgresDB) Delete(id int, userId int) error {
	updateStmt := "DELETE FROM todos WHERE id=$1 AND userid=$2;"
	_, err := db.sql.Exec(updateStmt, id, userId)
	if err != nil {
		return err
	}

	return nil
}

func (db *postgresDB) Register(input *models.RegisterRequest) error {
	query := "INSERT INTO users(login, passwordhash) values($1, $2) RETURNING id;"
	row := db.sql.QueryRow(query, input.Login, db.hash(*input.Password))
	var id int
	if err := row.Scan(&id); err != nil {
		return err
	}

	return nil
}

func (db *postgresDB) Login(input *models.LoginRequest) (*int, error) {
	query := "SELECT id FROM users WHERE login=$1 AND passwordhash=$2;"
	row := db.sql.QueryRow(query, input.Login, db.hash(*input.Password))
	var id int
	if err := row.Scan(&id); err != nil {
		return nil, err
	}

	return &id, nil
}

func (db *postgresDB) hash(in string) string {
	input := fmt.Sprintf("salt_%s_salt", in)
	res := md5.Sum([]byte(input))
	return hex.EncodeToString(res[:])
}
