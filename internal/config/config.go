package config

import "os"

type Config struct {
	AppEnv     string
	AppName    string
	LogLevel   string
	LogFile    string
	HTTPPort   string
	DBHost     string
	DBPort     string
	DBDatabase string
	DBUsername string
	DBPassword string
}

func Load() Config {
	return Config{
		AppEnv:     os.Getenv("APP_ENV"),
		AppName:    os.Getenv("APP_NAME"),
		LogLevel:   os.Getenv("LOG_LEVEL"),
		LogFile:    os.Getenv("LOG_FILE"),
		HTTPPort:   os.Getenv("HTTP_PORT"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBDatabase: os.Getenv("DB_DATABASE"),
		DBUsername: os.Getenv("DB_USERNAME"),
		DBPassword: os.Getenv("DB_PASSWORD"),
	}
}
