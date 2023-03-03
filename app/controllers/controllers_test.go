package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ann-96/todo-go-backend/app/controllers"
	"github.com/ann-96/todo-go-backend/app/models"
	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

const secretKey = "my-secret-key-my-secret-key-my-secret-key"

func TestRegisterAndLogin(t *testing.T) {
	h, err := controllers.NewUserController(controllers.Settings{JwtKey: secretKey}, getMockSQL())
	require.NoError(t, err)

	login, passw := "useruser", "Somepassw@1"
	registerReq := models.RegisterRequest{
		LoginRequest: models.LoginRequest{
			Login:    &login,
			Password: &passw,
		},
		Password2: &passw,
	}

	c, rec := getRequestContext(t, http.MethodPost, registerReq, h.NewContext)
	require.NoError(t, h.Register(c))
	require.Equal(t, http.StatusCreated, rec.Code, rec.Body.String())

	c, rec = getRequestContext(t, http.MethodPost, registerReq.LoginRequest, h.NewContext)
	require.NoError(t, h.Login(c))
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())

	err = json.Unmarshal(rec.Body.Bytes(), &jwtToken)
	require.NoError(t, err)
}

func TestLoginLength(t *testing.T) {

	h, err := controllers.NewUserController(controllers.Settings{JwtKey: secretKey}, getMockSQL())
	require.NoError(t, err)

	tooLongLogin := getRandomString(120)
	tooShortLogin := getRandomString(2)
	passw := "Passwd@1"

	registerReq := models.RegisterRequest{
		LoginRequest: models.LoginRequest{
			Login:    &tooLongLogin,
			Password: &passw,
		},
		Password2: &passw,
	}

	c, rec := getRequestContext(t, http.MethodPost, registerReq, h.NewContext)
	require.NoError(t, h.Register(c))
	require.Equal(t, http.StatusBadRequest, rec.Code, rec.Body.String())

	errResp := models.ErrResponse{}
	body := rec.Body.Bytes()
	err = json.Unmarshal(body, &errResp)
	require.NoError(t, err, "cannot unmarshal the response into error response")
	require.Equal(t, "login is too long", errResp.Msg)

	registerReq.Login = &tooShortLogin
	c, rec = getRequestContext(t, http.MethodPost, registerReq, h.NewContext)
	require.NoError(t, h.Register(c))
	require.Equal(t, http.StatusBadRequest, rec.Code, rec.Body.String())

	errResp = models.ErrResponse{}
	body = rec.Body.Bytes()
	err = json.Unmarshal(body, &errResp)
	require.NoError(t, err, "cannot unmarshal the response into error response")
	require.Equal(t, "login is too short", errResp.Msg)
}

func TestRegisterPasswordLength(t *testing.T) {
	h, err := controllers.NewUserController(controllers.Settings{JwtKey: secretKey}, getMockSQL())
	require.NoError(t, err)

	login := "validuser"
	tooShortPassword := "abc"
	tooLongPassword := getRandomString(121)

	// Test password that is too short
	registerReq := models.RegisterRequest{
		LoginRequest: models.LoginRequest{
			Login:    &login,
			Password: &tooShortPassword,
		},
		Password2: &tooShortPassword,
	}
	c, rec := getRequestContext(t, http.MethodPost, registerReq, h.NewContext)
	require.NoError(t, h.Register(c))
	require.Equal(t, http.StatusBadRequest, rec.Code, rec.Body.String())

	errResp := models.ErrResponse{}
	body := rec.Body.Bytes()
	err = json.Unmarshal(body, &errResp)
	require.NoError(t, err, "cannot unmarshal the response into error response")
	require.Equal(t, "insecure password, try including more special characters, using uppercase letters, using numbers or using a longer password", errResp.Msg)

	// Test password that is too long
	registerReq.Password = &tooLongPassword
	registerReq.Password2 = &tooLongPassword
	c, rec = getRequestContext(t, http.MethodPost, registerReq, h.NewContext)
	require.NoError(t, h.Register(c))
	require.Equal(t, http.StatusBadRequest, rec.Code, rec.Body.String())

	errResp = models.ErrResponse{}
	body = rec.Body.Bytes()
	err = json.Unmarshal(body, &errResp)
	require.NoError(t, err, "cannot unmarshal the response into error response")
	require.Equal(t, "password is too long", errResp.Msg)
}

func TestRegisterLoginAlphanum(t *testing.T) {

	h, err := controllers.NewUserController(controllers.Settings{JwtKey: secretKey}, getMockSQL())
	require.NoError(t, err)

	login := "asdfasdf1@"
	passw := "Passwd1@"

	registerReq := models.RegisterRequest{
		LoginRequest: models.LoginRequest{
			Login:    &login,
			Password: &passw,
		},
		Password2: &passw,
	}

	c, rec := getRequestContext(t, http.MethodPost, registerReq, h.NewContext)
	require.NoError(t, h.Register(c))
	require.Equal(t, http.StatusBadRequest, rec.Code, rec.Body.String())

	errResp := models.ErrResponse{}
	body := rec.Body.Bytes()
	err = json.Unmarshal(body, &errResp)
	require.NoError(t, err, "cannot unmarshal the response into error response")
	require.Equal(t, "login should only contain letters and numbers", errResp.Msg)

}

func TestUpdateTodos(t *testing.T) {
	mockSQL := getMockSQL()

	userController, err := controllers.NewUserController(controllers.Settings{JwtKey: secretKey}, mockSQL)
	require.NoError(t, err)
	todoController, err := controllers.NewTodoController(controllers.Settings{JwtKey: secretKey}, mockSQL)
	require.NoError(t, err)

	login := "loginlogin"
	passw := "Passwd@jwklfnjknfkj1"

	registerReq := models.RegisterRequest{
		LoginRequest: models.LoginRequest{
			Login:    &login,
			Password: &passw,
		},
		Password2: &passw,
	}

	todoText := getRandomString(200)
	todoCompleted := true

	addTodo := models.AddTodoRequest{
		Text:      &todoText,
		Completed: &todoCompleted,
	}

	newTodoText := getRandomString(300)
	newTodoCompleted := true

	c, rec := getRequestContext(t, http.MethodPost, registerReq, userController.NewContext)
	require.NoError(t, userController.Register(c))
	require.Equal(t, http.StatusCreated, rec.Code, rec.Body.String())

	c, rec = getRequestContext(t, http.MethodPost, registerReq.LoginRequest, userController.NewContext)
	require.NoError(t, userController.Login(c))
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())

	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &jwtToken))

	c, rec = getRequestContext(t, http.MethodPost, addTodo, userController.NewContext)
	userId, err := controllers.TokenToUserID(jwtToken, secretKey)
	require.NoError(t, err)
	c.Set("userId", userId)

	require.NoError(t, todoController.Add(c))
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())

	resultTodo := models.Todo{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resultTodo))
	require.Equal(t, addTodo, resultTodo.AddTodoRequest)

	updateTodo := models.Todo{
		Id: resultTodo.Id,
		AddTodoRequest: models.AddTodoRequest{
			Text:      &newTodoText,
			Completed: &newTodoCompleted,
		},
	}
	c, rec = getRequestContext(t, http.MethodPost, updateTodo, userController.NewContext)
	userId, err = controllers.TokenToUserID(jwtToken, secretKey)
	require.NoError(t, err)
	c.Set("userId", userId)

	require.NoError(t, todoController.Update(c))
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())

	resTodo := models.Todo{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resTodo))
	require.Equal(t, updateTodo, resTodo)
}

func TestTodos(t *testing.T) {
	// controllers
	mockSQL := getMockSQL()

	userController, err := controllers.NewUserController(controllers.Settings{JwtKey: secretKey}, mockSQL)
	require.NoError(t, err)

	todoController, err := controllers.NewTodoController(controllers.Settings{JwtKey: secretKey}, mockSQL)
	require.NoError(t, err)

	// requests
	login := "loginlogin"
	passw := "Passwd@jwklfnjknfkj1"

	registerReq := models.RegisterRequest{
		LoginRequest: models.LoginRequest{
			Login:    &login,
			Password: &passw,
		},
		Password2: &passw,
	}

	todoText := getRandomString(200)
	todoCompleted := true

	addTodo := models.AddTodoRequest{
		Text:      &todoText,
		Completed: &todoCompleted,
	}

	// register a user
	c, rec := getRequestContext(t, http.MethodPost, registerReq, userController.NewContext)
	require.NoError(t, userController.Register(c))
	require.Equal(t, http.StatusCreated, rec.Code, rec.Body.String())

	// login the user
	c, rec = getRequestContext(t, http.MethodPost, registerReq.LoginRequest, userController.NewContext)
	require.NoError(t, userController.Login(c))
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())

	// save session id for future use
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &jwtToken))

	// create a todo entry
	c, rec = getRequestContext(t, http.MethodPost, addTodo, userController.NewContext)
	userId, err := controllers.TokenToUserID(jwtToken, secretKey)
	require.NoError(t, err)
	c.Set("userId", userId)

	require.NoError(t, todoController.Add(c))
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())

	resultTodo := models.Todo{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resultTodo))
	require.Equal(t, addTodo, resultTodo.AddTodoRequest)
}

func TestRasswordComplexity(t *testing.T) {
	h, err := controllers.NewUserController(controllers.Settings{JwtKey: secretKey}, getMockSQL())
	require.NoError(t, err)

	login := "loginlogin"
	passw := "password"

	registerReq := models.RegisterRequest{
		LoginRequest: models.LoginRequest{
			Login:    &login,
			Password: &passw,
		},
		Password2: &passw,
	}

	c, rec := getRequestContext(t, http.MethodPost, registerReq, h.NewContext)
	require.NoError(t, h.Register(c))
	require.Equal(t, http.StatusBadRequest, rec.Code, rec.Body.String())

	errResp := models.ErrResponse{}
	body := rec.Body.Bytes()
	err = json.Unmarshal(body, &errResp)
	require.NoError(t, err, "cannot unmarshal the response into error response")
	require.Equal(t, "insecure password, try including more special characters, using uppercase letters, using numbers or using a longer password", errResp.Msg)

}

func TestRasswordMatch(t *testing.T) {
	h, err := controllers.NewUserController(controllers.Settings{JwtKey: secretKey}, getMockSQL())
	require.NoError(t, err)

	login := "loginloginlogin"
	passw := "password333222@@@password"
	passwtwo := "ppassword333222@@@password"

	registerReq := models.RegisterRequest{
		LoginRequest: models.LoginRequest{
			Login:    &login,
			Password: &passw,
		},
		Password2: &passwtwo,
	}
	c, rec := getRequestContext(t, http.MethodPost, registerReq, h.NewContext)
	require.NoError(t, h.Register(c))
	require.Equal(t, http.StatusBadRequest, rec.Code, rec.Body.String())

	errResp := models.ErrResponse{}
	body := rec.Body.Bytes()
	err = json.Unmarshal(body, &errResp)
	require.NoError(t, err, "cannot unmarshal the response into error response")
	require.Equal(t, "Entered passwords didn't match", errResp.Msg)

}

func TestLoginLetterCase(t *testing.T) {
	h, err := controllers.NewUserController(controllers.Settings{JwtKey: secretKey}, getMockSQL())
	require.NoError(t, err)

	login := "loginlogin"
	passw := "Passwd@jwklfnjknfkj1"

	registerReq := models.RegisterRequest{
		LoginRequest: models.LoginRequest{
			Login:    &login,
			Password: &passw,
		},
		Password2: &passw,
	}

	c, rec := getRequestContext(t, http.MethodPost, registerReq, h.NewContext)
	require.NoError(t, h.Register(c))
	require.Equal(t, http.StatusCreated, rec.Code, rec.Body.String())

	login = "lOginlogin"
	c, rec = getRequestContext(t, http.MethodPost, registerReq, h.NewContext)
	require.NoError(t, h.Register(c))
	require.Equal(t, http.StatusConflict, rec.Code, rec.Body.String())

	errResp := models.ErrResponse{}
	body := rec.Body.Bytes()
	err = json.Unmarshal(body, &errResp)
	require.NoError(t, err, "cannot unmarshal the response into error response")
	require.Equal(t, "the user already exists", errResp.Msg)

	c, rec = getRequestContext(t, http.MethodPost, registerReq.LoginRequest, h.NewContext)
	require.NoError(t, h.Login(c))
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
}

type userIdType int
type todoIdType int
type mockSqlDB struct {
	todos      map[userIdType]map[todoIdType]*models.Todo
	nextTodoID todoIdType
	users      map[userIdType]*models.LoginRequest
	nextUserID userIdType
}

var (
	jwtToken string
)

func getMockSQL() *mockSqlDB {
	return &mockSqlDB{
		todos: make(map[userIdType]map[todoIdType]*models.Todo),
		users: make(map[userIdType]*models.LoginRequest),
	}
}

func (db *mockSqlDB) Update(input *models.Todo, userId int) (*models.Todo, error) {
	_, ok := db.users[userIdType(userId)]
	if !ok {
		return nil, fmt.Errorf("the user id %v is not found", userId)
	}

	todo, ok := db.todos[userIdType(userId)][todoIdType(*input.Id)]
	if !ok {
		return nil, errors.New("entry not found for the user")
	}

	if input.Completed != nil {
		todo.Completed = input.Completed
	}
	if input.Text != nil {
		todo.Text = input.Text
	}

	return todo, nil
}

func (db *mockSqlDB) List(start int, count int, userId int) (*models.TodoList, error) {
	todosCount := len(db.todos[userIdType(userId)])
	if todosCount == 0 || todosCount < start+count {
		return nil, nil
	}
	completedCount := 0
	for i := range db.todos[userIdType(userId)] {
		if *db.todos[userIdType(userId)][i].Completed {
			completedCount++
		}
	}

	res := models.TodoList{
		Count:          todosCount,
		CompletedCount: completedCount,
		List:           make([]models.Todo, 0, count),
	}

	counter := 0
	for i := range db.todos[userIdType(userId)] {
		if counter < start {
			counter++
			continue
		}
		if counter >= start+count {
			counter++
			break
		}
		res.List = append(res.List, *db.todos[userIdType(userId)][i])
		counter++
	}

	return &res, nil
}

func (db *mockSqlDB) Add(input *models.AddTodoRequest, userId int) (*models.Todo, error) {
	if db.todos[userIdType(userId)] == nil {
		db.todos[userIdType(userId)] = make(map[todoIdType]*models.Todo)
	}
	id := int(db.nextTodoID)
	db.todos[userIdType(userId)][db.nextTodoID] = &models.Todo{
		Id:             &id,
		AddTodoRequest: *input,
	}
	defer func() { db.nextTodoID++ }()

	return db.todos[userIdType(userId)][db.nextTodoID], nil
}

func (db *mockSqlDB) Delete(id int, userID int) error {
	delete(db.todos[userIdType(userID)], todoIdType(id))

	return nil
}

func (db *mockSqlDB) Register(input *models.RegisterRequest) error {
	exists := false
	for i := range db.users {
		if *db.users[i].Login == *input.Login && *db.users[i].Password == *input.Password {
			exists = true
		}
	}
	if exists {
		return fmt.Errorf(`pq: duplicate key value violates unique constraint "users_login_key"`)
	}
	db.users[db.nextUserID] = &input.LoginRequest
	db.nextUserID++
	return nil
}

func (db *mockSqlDB) Login(input *models.LoginRequest) (userIdType *int, err error) {
	for i := range db.users {
		if *db.users[i].Login == *input.Login && *db.users[i].Password == *input.Password {
			userIdType = (*int)(&i)
		}
	}
	if userIdType == nil {
		err = fmt.Errorf("user or password is invalid")
	}
	return
}

func (db *mockSqlDB) Migrate() error {
	return nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func getRandomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

type contextGen func(r *http.Request, w http.ResponseWriter) echo.Context

func getRequestContext(t require.TestingT, method string, in interface{}, newContext contextGen) (echo.Context, *httptest.ResponseRecorder) {
	reqJson, err := json.Marshal(in)
	require.NoError(t, err)

	req := httptest.NewRequest(method, "/", bytes.NewReader(reqJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, jwtToken)

	rec := httptest.NewRecorder()
	require.NoError(t, err)

	return newContext(req, rec), rec
}
