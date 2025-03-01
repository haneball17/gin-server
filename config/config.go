package config

import (
	"log"
	"os"
)

// Config 结构体定义配置项
type Config struct {
	// 通用配置
	DebugLevel string

	// 服务器配置
	ServerPort string

	// 主数据库配置
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string

	// Radius数据库配置
	RadiusDBUser     string
	RadiusDBPassword string
	RadiusDBHost     string
	RadiusDBPort     string
	RadiusDBName     string
}

// 全局配置实例
var globalConfig *Config

// InitConfig 初始化配置
func InitConfig() {
	log.Println("初始化全局配置...")

	globalConfig = &Config{
		// 通用配置
		DebugLevel: getEnv("DEBUG_LEVEL", "false"),

		// 服务器配置
		ServerPort: getEnv("SERVER_PORT", "8080"),

		// 主数据库配置
		DBUser:     getEnv("DB_USER", "gin_user"),
		DBPassword: getEnv("DB_PASSWORD", "your_password"),
		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBName:     getEnv("DB_NAME", "gin_server"),

		// Radius数据库配置 (默认与主数据库相同的连接信息，但不同的数据库名)
		RadiusDBUser:     getEnv("RADIUS_DB_USER", getEnv("DB_USER", "gin_user")),
		RadiusDBPassword: getEnv("RADIUS_DB_PASSWORD", getEnv("DB_PASSWORD", "your_password")),
		RadiusDBHost:     getEnv("RADIUS_DB_HOST", getEnv("DB_HOST", "127.0.0.1")),
		RadiusDBPort:     getEnv("RADIUS_DB_PORT", getEnv("DB_PORT", "3306")),
		RadiusDBName:     getEnv("RADIUS_DB_NAME", "radius"),
	}

	if globalConfig.DebugLevel == "true" {
		log.Println("配置初始化完成")
		log.Printf("服务器端口: %s\n", globalConfig.ServerPort)
		log.Printf("调试级别: %s\n", globalConfig.DebugLevel)
		log.Printf("主数据库: %s@%s:%s/%s\n", globalConfig.DBUser, globalConfig.DBHost, globalConfig.DBPort, globalConfig.DBName)
		log.Printf("Radius数据库: %s@%s:%s/%s\n", globalConfig.RadiusDBUser, globalConfig.RadiusDBHost, globalConfig.RadiusDBPort, globalConfig.RadiusDBName)
	}
}

// GetConfig 获取配置
func GetConfig() *Config {
	if globalConfig == nil {
		InitConfig()
	}
	return globalConfig
}

// getEnv 获取环境变量
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
