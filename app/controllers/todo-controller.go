package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ann-96/todo-go-backend/app/db"
	"github.com/ann-96/todo-go-backend/app/models"
	"github.com/ann-96/todo-go-backend/app/redis"
	"github.com/ann-96/todo-go-backend/app/tools"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type todoController struct {
	e  *echo.Echo
	db db.TodoSqlDB
	Settings
	userCache redis.SessionCache
}

func NewTodoController(settings Settings, db db.TodoSqlDB, userCache redis.SessionCache) (*todoController, error) {
	e := echo.New()
	e.Validator = tools.NewValidator()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "Authorization"},
	}))
	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:authorization",
		Validator: func(key string, c echo.Context) (bool, error) {
			res := userCache.GetSession(key)
			if res == nil {
				return false, nil
			}
			c.Set("userId", *res)
			return true, nil
		},
	}))

	controller := &todoController{
		Settings:  settings,
		userCache: userCache,
		db:        db,
		e:         e,
	}

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
	userId := c.Get("userId").(int)

	req := &models.AddTodoRequest{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrResponse{Msg: err.Error()})
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrResponse{Msg: err.Error()})
	}

	res, err := controller.db.Add(req, userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.ErrResponse{Msg: err.Error()})
	}

	return c.JSON(http.StatusOK, *res)
}

func (controller *todoController) Update(c echo.Context) error {
	userId := c.Get("userId").(int)

	req := &models.Todo{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrResponse{Msg: err.Error()})
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrResponse{Msg: err.Error()})
	}

	res, err := controller.db.Update(req, userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrResponse{Msg: err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

func (controller *todoController) List(c echo.Context) error {
	userId := c.Get("userId").(int)

	req := &models.ListRequest{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrResponse{Msg: err.Error()})
	}
	startVal, err := strconv.Atoi(req.Start)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrResponse{Msg: err.Error()})
	}
	countVal, err := strconv.Atoi(req.Count)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrResponse{Msg: err.Error()})
	}

	res, err := controller.db.List(startVal, countVal, userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.ErrResponse{Msg: err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

func (controller *todoController) Delete(c echo.Context) error {
	userId := c.Get("userId").(int)

	req := &struct {
		Id int `json:"id" binding:"validate"`
	}{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrResponse{Msg: err.Error()})
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrResponse{Msg: err.Error()})
	}

	err := controller.db.Delete(req.Id, userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, &models.ErrResponse{Msg: err.Error()})
	}

	return c.JSON(http.StatusOK, nil)
}

func (controller *todoController) NewContext(r *http.Request, w http.ResponseWriter) echo.Context {
	return controller.e.NewContext(r, w)
}
