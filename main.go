package main

import (
	"github.com/ann-96/todo-go-backend/app"
	"github.com/spf13/viper"
)

func main() {
	app := app.App{}

	viper.SetDefault("HTTP_HOST", "localhost")
	viper.SetDefault("HTTP_PORT", "8080")
	viper.SetDefault("SQL_HOST", "localhost")
	viper.SetDefault("SQL_PORT", "5432")
	viper.SetDefault("SQL_USER", "postgres")
	viper.SetDefault("SQL_PASS", "postgres")
	viper.SetDefault("SQL_DBNAME", "postgres")

	viper.BindEnv("HTTP_HOST")
	viper.BindEnv("HTTP_PORT")
	app.Host = viper.GetString("HTTP_HOST")
	app.Port = viper.GetString("HTTP_PORT")

	viper.BindEnv("SQL_HOST")
	viper.BindEnv("SQL_PORT")
	viper.BindEnv("SQL_USER")
	viper.BindEnv("SQL_PASS")
	viper.BindEnv("SQL_DBNAME")
	app.SqlHost = viper.GetString("SQL_HOST")
	app.SqlPort = viper.GetString("SQL_PORT")
	app.SqlUser = viper.GetString("SQL_USER")
	app.SqlPass = viper.GetString("SQL_PASS")
	app.SqlName = viper.GetString("SQL_DBNAME")

	app.Run()
}
