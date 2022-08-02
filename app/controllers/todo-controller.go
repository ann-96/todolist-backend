package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"

	"github.com/ann-96/todo-go-backend/app/db"
	"github.com/ann-96/todo-go-backend/app/models"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type todoController struct {
	e  *echo.Echo
	db db.TodoSqlDB
	Settings
}

type Validator struct {
	validator *validator.Validate
}

type errResponse struct {
	Msg string `json:"msg"`
}

func (cv *Validator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func NewTodoController(settings Settings) (*todoController, error) {
	e := echo.New()

	e.Validator = &Validator{validator: validator.New()}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	controller := &todoController{Settings: settings}
	db, err := db.CreatePostgresDB(
		db.Settings{
			IP:       settings.SqlHost,
			Port:     settings.SqlPort,
			User:     settings.SqlUser,
			Password: settings.SqlPass,
			Name:     settings.SqlName,
		},
	)
	if err != nil {
		return nil, err
	}

	err = db.Migrate()
	if err != nil {
		return nil, err
	}

	controller.db = db
	controller.e = e

	controller.e.POST("/todo/add", controller.Add)
	controller.e.POST("/todo/update", controller.Update)
	controller.e.GET("/todo/list", controller.List)
	controller.e.POST("/todo/delete", controller.Delete)

	return controller, nil
}

func (controller *todoController) Run() error {
	connectionString := fmt.Sprintf("%v:%v", controller.Host, controller.Port)
	return controller.e.Start(connectionString)
}

func (controller *todoController) Shutdown(ctx context.Context) error {
	return controller.e.Shutdown(ctx)
}

func (controller *todoController) Logger() echo.Logger {
	return controller.e.Logger
}

func (controller *todoController) Add(c echo.Context) error {
	req := &models.AddTodoRequest{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, &errResponse{Msg: err.Error()})
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, &errResponse{Msg: err.Error()})
	}

	res, err := controller.db.Add(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &errResponse{Msg: err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

func (controller *todoController) Update(c echo.Context) error {
	req := &models.Todo{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, &errResponse{Msg: err.Error()})
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, &errResponse{Msg: err.Error()})
	}

	res, err := controller.db.Update(req)
	if err != nil {
		return c.JSON(http.StatusNotFound, &errResponse{Msg: err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

func (controller *todoController) List(c echo.Context) error {

	req := &models.ListRequest{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, &errResponse{Msg: err.Error()})
	}
	startVal, err := strconv.Atoi(req.Start)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &errResponse{Msg: err.Error()})
	}
	countVal, err := strconv.Atoi(req.Count)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &errResponse{Msg: err.Error()})
	}

	res, err := controller.db.List(startVal, countVal)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &errResponse{Msg: err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

func (controller *todoController) Delete(c echo.Context) error {
	req := &struct {
		Id int `json:"id" binding:"validate"`
	}{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, &errResponse{Msg: err.Error()})
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, &errResponse{Msg: err.Error()})
	}

	err := controller.db.Delete(req.Id)
	if err != nil {
		return c.JSON(http.StatusNotFound, &errResponse{Msg: err.Error()})
	}

	return c.JSON(http.StatusOK, nil)
}
