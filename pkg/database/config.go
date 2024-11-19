package database

import "github.com/hilmiikhsan/simple-messaging-app/pkg/env"

type Config struct {
	DBUser        string
	DBPassword    string
	DBHost        string
	DBPort        string
	DBName        string
	AppHost       string
	AppPort       string
	AppName       string
	AppSecret     string
	AppPortSocket string
	MongoDBUri    string
}

func LoadConfig() *Config {
	return &Config{
		DBUser:        env.GetEnv("DB_USER", ""),
		DBPassword:    env.GetEnv("DB_PASSWORD", ""),
		DBHost:        env.GetEnv("DB_HOST", "127.0.0.1"),
		DBPort:        env.GetEnv("DB_PORT", "3306"),
		DBName:        env.GetEnv("DB_NAME", ""),
		AppHost:       env.GetEnv("APP_HOST", "0.0.0.0"),
		AppPort:       env.GetEnv("APP_PORT", "4000"),
		AppName:       env.GetEnv("APP_NAME", ""),
		AppSecret:     env.GetEnv("APP_SECRET", ""),
		AppPortSocket: env.GetEnv("APP_PORT_SOCKET", ":4001"),
		MongoDBUri:    env.GetEnv("MONGODB_URI", ""),
	}
}
