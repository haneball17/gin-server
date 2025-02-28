package model

import (
	"os"
)

// Config 结构体定义配置项
type Config struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	ServerPort string
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	return &Config{
		DBUser:     getEnv("DB_USER", "gin_user"),
		DBPassword: getEnv("DB_PASSWORD", "your_password"),
		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBName:     getEnv("DB_NAME", "gin_server"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}
}

// getEnv 获取环境变量
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
