package config

import "os"

type Config struct {
	AppEnv     string
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
		HTTPPort:   os.Getenv("HTTP_PORT"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBDatabase: os.Getenv("DB_DATABASE"),
		DBUsername: os.Getenv("DB_USERNAME"),
		DBPassword: os.Getenv("DB_PASSWORD"),
	}
}
