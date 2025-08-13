package config

import (
	"os"

	"github.com/joho/godotenv"
)

type (
	Container struct {
		App   *App
		Token *Token
		DB    *DB
		HTTP  *HTTP
	}
	App struct {
		Name string
		Env  string
	}
	Token struct {
		Duration        string
		Secret          string
		DurationRefresh string
		SecretRefresh   string
	}
	DB struct {
		Connection string
		Host       string
		Port       string
		User       string
		Password   string
		Name       string
		SSL        string
	}
	HTTP struct {
		Env            string
		URL            string
		Port           string
		AllowedOrigins string
	}
)

func New() (*Container, error) {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			return nil, err
		}
	}

	app := &App{
		Name: os.Getenv("APP_NAME"),
		Env:  os.Getenv("APP_ENV"),
	}

	token := &Token{
		Duration:        os.Getenv("TOKEN_DURATION"),
		Secret:          os.Getenv("TOKEN_SECRET"),
		DurationRefresh: os.Getenv("REFRESH_TOKEN_DURATION"),
		SecretRefresh:   os.Getenv("REFRESH_TOKEN_SECRET"),
	}

	db := &DB{
		Connection: os.Getenv("DB_CONNECTION"),
		Host:       os.Getenv("DB_HOST"),
		Port:       os.Getenv("DB_PORT"),
		User:       os.Getenv("DB_USER"),
		Password:   os.Getenv("DB_PASSWORD"),
		Name:       os.Getenv("DB_NAME"),
		SSL:        os.Getenv("DB_SSL"),
	}

	http := &HTTP{
		Env:            os.Getenv("APP_ENV"),
		URL:            os.Getenv("HTTP_URL"),
		Port:           os.Getenv("HTTP_PORT"),
		AllowedOrigins: os.Getenv("HTTP_ALLOWED_ORIGINS"),
	}

	return &Container{
		app,
		token,
		db,
		http,
	}, nil
}
