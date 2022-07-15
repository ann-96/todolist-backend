package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/ann-96/todo-go-backend/app/controllers"
)

type App struct {
	controllers.Settings
}

func (app *App) Run() {
	const serviceNum = 1
	var wg sync.WaitGroup

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	wg.Add(serviceNum)

	notifyQuit := make([]chan struct{}, serviceNum)
	for i := range notifyQuit {
		notifyQuit[i] = make(chan struct{}, 1)
	}
	go app.runTodoRest(&wg, notifyQuit[0])

	<-quit
	for i := range notifyQuit {
		notifyQuit[i] <- struct{}{}
	}

	wg.Wait()
}

func (app *App) runTodoRest(wg *sync.WaitGroup, quit chan struct{}) {
	c, err := controllers.NewTodoController(app.Settings)
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
