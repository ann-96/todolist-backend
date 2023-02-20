package main

import (
	"github.com/ann-96/todo-go-backend/app"
	"github.com/ann-96/todo-go-backend/app/controllers"
	"github.com/spf13/viper"
)

func main() {
	app := app.App{}
	commonSettings := controllers.Settings{}
	viper.SetDefault("TODO_HTTP_HOST", "localhost")
	viper.SetDefault("TODO_HTTP_PORT", "8080")
	viper.SetDefault("USER_HTTP_HOST", "localhost")
	viper.SetDefault("USER_HTTP_PORT", "9080")
	viper.SetDefault("SQL_HOST", "localhost")
	viper.SetDefault("SQL_PORT", "5432")
	viper.SetDefault("SQL_USER", "postgres")
	viper.SetDefault("SQL_PASS", "postgres")
	viper.SetDefault("SQL_DBNAME", "postgres")
	viper.SetDefault("JWT_KEY", "my-secret-key-my-secret-key-my-secret-key")

	viper.BindEnv("SQL_HOST")
	viper.BindEnv("SQL_PORT")
	viper.BindEnv("SQL_USER")
	viper.BindEnv("SQL_PASS")
	viper.BindEnv("SQL_DBNAME")
	viper.BindEnv("JWT_KEY")
	commonSettings.SqlHost = viper.GetString("SQL_HOST")
	commonSettings.SqlPort = viper.GetString("SQL_PORT")
	commonSettings.SqlUser = viper.GetString("SQL_USER")
	commonSettings.SqlPass = viper.GetString("SQL_PASS")
	commonSettings.SqlName = viper.GetString("SQL_DBNAME")
	commonSettings.JwtKey = viper.GetString("JWT_KEY")

	app.UserController = commonSettings
	app.TodoController = commonSettings

	viper.BindEnv("USER_HTTP_HOST")
	viper.BindEnv("USER_HTTP_PORT")
	app.UserController.Host = viper.GetString("USER_HTTP_HOST")
	app.UserController.Port = viper.GetString("USER_HTTP_PORT")

	viper.BindEnv("TODO_HTTP_HOST")
	viper.BindEnv("TODO_HTTP_PORT")
	app.TodoController.Host = viper.GetString("TODO_HTTP_HOST")
	app.TodoController.Port = viper.GetString("TODO_HTTP_PORT")

	app.Run()
}
