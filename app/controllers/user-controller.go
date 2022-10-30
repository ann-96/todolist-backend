package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ann-96/todo-go-backend/app/db"
	"github.com/ann-96/todo-go-backend/app/models"
	"github.com/ann-96/todo-go-backend/app/redis"
	"github.com/ann-96/todo-go-backend/app/tools"
	"github.com/google/uuid"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type userController struct {
	e  *echo.Echo
	db db.TodoSqlDB
	Settings
	userCache redis.SessionCache
}

func NewUserController(settings Settings, db db.TodoSqlDB, userCache redis.SessionCache) (*userController, error) {

	e := echo.New()
	e.Validator = tools.NewValidator()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	// e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
	// 	fmt.Printf("%s", string(reqBody))
	// }))

	usercontroller := &userController{
		Settings:  settings,
		userCache: userCache,
		db:        db,
		e:         e,
	}

	usercontroller.e.POST("/users/register", usercontroller.Register)
	usercontroller.e.POST("/users/login", usercontroller.Login)
	usercontroller.e.POST("/users/logout", usercontroller.Logout)

	return usercontroller, nil
}

func (controller *userController) Run() error {
	connectionString := fmt.Sprintf("%v:%v", controller.Host, controller.Port)
	return controller.e.Start(connectionString)
}

func (controller *userController) Shutdown(ctx context.Context) error {
	return controller.e.Shutdown(ctx)
}

func (controller *userController) Logger() echo.Logger {
	return controller.e.Logger
}

func (controller *userController) Register(c echo.Context) error {
	req := &models.RegisterRequest{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrResponse{Msg: err.Error()})
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrResponse{Msg: err.Error()})
	}

	err := controller.db.Register(req)
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "users_login_key"` {
			return c.JSON(http.StatusConflict, &models.ErrResponse{Msg: "the user already exists"})
		}
		return c.JSON(http.StatusInternalServerError, &models.ErrResponse{Msg: err.Error()})
	}

	return c.JSON(http.StatusCreated, nil)
}

func (controller *userController) Login(c echo.Context) error {
	req := &models.LoginRequest{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrResponse{Msg: err.Error()})
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrResponse{Msg: err.Error()})
	}

	userId, err := controller.db.Login(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.ErrResponse{Msg: err.Error()})
	}

	id := uuid.New()
	uuid := id.String()
	controller.userCache.CreateSession(uuid, *userId)

	return c.JSON(http.StatusOK, uuid)
}

func (controller *userController) Logout(c echo.Context) error {
	if auth, ok := c.Request().Header["Authorization"]; ok && auth[0] != "" {
		controller.userCache.DeleteSession(auth[0])
	}
	return c.JSON(http.StatusNoContent, nil)
}

func (controller *userController) NewContext(r *http.Request, w http.ResponseWriter) echo.Context {
	return controller.e.NewContext(r, w)
}
