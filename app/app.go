package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/ann-96/todo-go-backend/app/controllers"
	"github.com/ann-96/todo-go-backend/app/db"
)

type App struct {
	TodoController controllers.Settings
	UserController controllers.Settings
}

func (app *App) Run() {

	const serviceNum = 2
	var wg sync.WaitGroup

	quit := make(chan os.Signal, serviceNum)
	signal.Notify(quit, os.Interrupt)

	wg.Add(serviceNum)

	notifyQuit := make([]chan struct{}, serviceNum)
	for i := range notifyQuit {
		notifyQuit[i] = make(chan struct{}, serviceNum)
	}
	go app.runTodoRest(&wg, notifyQuit[0])
	go app.runUserRest(&wg, notifyQuit[1])

	<-quit
	for i := range notifyQuit {
		notifyQuit[i] <- struct{}{}
	}

	wg.Wait()
}

func (app *App) runTodoRest(wg *sync.WaitGroup, quit chan struct{}) {
	db, err := getDBConnectionsFromSettings(&app.TodoController)
	if err != nil {
		panic(err)
	}

	c, err := controllers.NewTodoController(app.TodoController, db)
	if err != nil {
		panic(err)
	}

	errChan := make(chan error, 1)
	go func() {
		if err := c.Run(); err != nil {
			errChan <- err
		}
	}()

	select { // stopping on both an error or a SIGINT
	case err := <-errChan:
		c.Logger().Printf("Fatal error: %v\n", err)
	case <-quit:
		c.Logger().Printf("Shutting down gracefully")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := c.Shutdown(ctx); err != nil {
		c.Logger().Fatal(err)
	}

	wg.Done()
}

func (app *App) runUserRest(wg *sync.WaitGroup, quit chan struct{}) {
	db, err := getDBConnectionsFromSettings(&app.UserController)
	if err != nil {
		panic(err)
	}

	c, err := controllers.NewUserController(app.UserController, db)
	if err != nil {
		panic(err)
	}

	errChan := make(chan error, 1)
	go func() {
		if err := c.Run(); err != nil {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		c.Logger().Printf("Fatal error: %v\n", err)
	case <-quit:
		c.Logger().Printf("Shutting down gracefully")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := c.Shutdown(ctx); err != nil {
		c.Logger().Fatal(err)
	}

	wg.Done()
}

func getDBConnectionsFromSettings(settings *controllers.Settings) (sql db.TodoSqlDB, err error) {
	sql, err = db.CreatePostgresDB(
		db.Settings{
			IP:       settings.SqlHost,
			Port:     settings.SqlPort,
			User:     settings.SqlUser,
			Password: settings.SqlPass,
			Name:     settings.SqlName,
		},
	)
	if err != nil {
		return
	}

	err = sql.Migrate()
	return
}
