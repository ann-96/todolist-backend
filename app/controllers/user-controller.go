package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ann-96/todo-go-backend/app/db"
	"github.com/ann-96/todo-go-backend/app/models"
	"github.com/ann-96/todo-go-backend/app/tools"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

const passwordComplexity = 50

type userController struct {
	e  *echo.Echo
	db db.TodoSqlDB
	Settings
}

func NewUserController(settings Settings, db db.TodoSqlDB) (*userController, error) {

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
		Settings: settings,

		db: db,
		e:  e,
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
	if *req.Password != *req.Password2 {
		return c.JSON(http.StatusBadRequest, &models.ErrResponse{Msg: "Entered passwords didn't match"})
	}
	if err := passwordvalidator.Validate(*req.Password, passwordComplexity); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrResponse{Msg: err.Error()})
	}

	*req.Login = strings.ToLower(*req.Login)

	err := controller.db.Register(req)
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "users_login_key"` {
			return c.JSON(http.StatusConflict, &models.ErrResponse{Msg: "the user already exists"})
		}
		return c.JSON(http.StatusInternalServerError, &models.ErrResponse{Msg: err.Error()})
	}

	return c.NoContent(http.StatusCreated)
}

func (controller *userController) Login(c echo.Context) error {
	req := &models.LoginRequest{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrResponse{Msg: err.Error()})
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrResponse{Msg: err.Error()})
	}

	*req.Login = strings.ToLower(*req.Login)

	userId, err := controller.db.Login(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.ErrResponse{Msg: err.Error()})
	}

	expirationTime := time.Now().Add(30 * 24 * time.Hour)
	claims := &models.Claims{
		UserID: *userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Id:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(controller.Settings.JwtKey))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, tokenString)
}

func (controller *userController) Logout(c echo.Context) error {
	return c.JSON(http.StatusNoContent, nil)
}

func (controller *userController) NewContext(r *http.Request, w http.ResponseWriter) echo.Context {
	return controller.e.NewContext(r, w)
}

func TokenToUserID(input, key string) (int, error) {
	token, err := jwt.ParseWithClaims(input, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(key), nil
	})

	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*models.Claims)
	if ok && token.Valid {
		if claims == nil {
			return 0, fmt.Errorf("invalid token")
		}
		return claims.UserID, nil
	}
	return 0, fmt.Errorf("invalid token")
}
