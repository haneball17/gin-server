package config

import (
	"log"
	"os"
	"strconv"
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

	// 配置管理模块配置
	ConfigManager ConfigManagerConfig
}

// ConfigManagerConfig 配置管理模块配置
type ConfigManagerConfig struct {
	// 日志管理配置
	LogManager LogManagerConfig
	// 策略管理配置
	StrategyManager StrategyManagerConfig
	// 存储配置
	Storage StorageConfig
}

// LogManagerConfig 日志管理配置
type LogManagerConfig struct {
	// 日志生成间隔（分钟）
	GenerateInterval int
	// 是否启用加密
	EnableEncryption bool
	// 加密配置
	Encryption EncryptionConfig
}

// StrategyManagerConfig 策略管理配置
type StrategyManagerConfig struct {
	// 轮询间隔（秒）
	PollInterval int
	// 是否启用加密
	EnableEncryption bool
	// 加密配置
	Encryption EncryptionConfig
}

// StorageConfig 存储配置
type StorageConfig struct {
	// 存储类型 (gitee/ftp)
	Type string
	// Gitee配置
	Gitee GiteeConfig
	// FTP配置
	FTP FTPConfig
}

// GiteeConfig Gitee配置
type GiteeConfig struct {
	// 访问令牌
	AccessToken string
	// 仓库所有者
	Owner string
	// 仓库名称
	Repo string
	// 日志文件夹路径
	LogPath string
	// 策略文件夹路径
	StrategyPath string
}

// FTPConfig FTP配置
type FTPConfig struct {
	// FTP服务器地址
	Host string
	// FTP服务器端口
	Port int
	// 用户名
	Username string
	// 密码
	Password string
	// 日志文件夹路径
	LogPath string
	// 策略文件夹路径
	StrategyPath string
}

// EncryptionConfig 加密配置
type EncryptionConfig struct {
	// AES密钥长度
	AESKeyLength int
	// 公钥加密方式 (RSA/ECDSA/ED25519)
	PublicKeyAlgorithm string
	// 公钥长度
	PublicKeyLength int
	// 运维系统公钥路径
	PublicKeyPath string
	// 本系统私钥路径
	PrivateKeyPath string
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

		// 配置管理模块配置
		ConfigManager: ConfigManagerConfig{
			// 日志管理配置
			LogManager: LogManagerConfig{
				GenerateInterval: getEnvInt("LOG_GENERATE_INTERVAL", 5),
				EnableEncryption: getEnvBool("LOG_ENABLE_ENCRYPTION", false),
				Encryption: EncryptionConfig{
					AESKeyLength:       getEnvInt("LOG_AES_KEY_LENGTH", 256),
					PublicKeyAlgorithm: getEnv("LOG_PUBLIC_KEY_ALGORITHM", "RSA"),
					PublicKeyLength:    getEnvInt("LOG_PUBLIC_KEY_LENGTH", 2048),
					PublicKeyPath:      getEnv("LOG_PUBLIC_KEY_PATH", "keys/ops_public.pem"),
					PrivateKeyPath:     getEnv("LOG_PRIVATE_KEY_PATH", "keys/system_private.pem"),
				},
			},
			// 策略管理配置
			StrategyManager: StrategyManagerConfig{
				PollInterval:     getEnvInt("STRATEGY_POLL_INTERVAL", 60),
				EnableEncryption: getEnvBool("STRATEGY_ENABLE_ENCRYPTION", false),
				Encryption: EncryptionConfig{
					AESKeyLength:       getEnvInt("STRATEGY_AES_KEY_LENGTH", 256),
					PublicKeyAlgorithm: getEnv("STRATEGY_PUBLIC_KEY_ALGORITHM", "RSA"),
					PublicKeyLength:    getEnvInt("STRATEGY_PUBLIC_KEY_LENGTH", 2048),
					PublicKeyPath:      getEnv("STRATEGY_PUBLIC_KEY_PATH", "keys/ops_public.pem"),
					PrivateKeyPath:     getEnv("STRATEGY_PRIVATE_KEY_PATH", "keys/system_private.pem"),
				},
			},
			// 存储配置
			Storage: StorageConfig{
				Type: getEnv("STORAGE_TYPE", "gitee"),
				Gitee: GiteeConfig{
					AccessToken:  getEnv("GITEE_ACCESS_TOKEN", ""),
					Owner:        getEnv("GITEE_OWNER", ""),
					Repo:         getEnv("GITEE_REPO", ""),
					LogPath:      getEnv("GITEE_LOG_PATH", "logs"),
					StrategyPath: getEnv("GITEE_STRATEGY_PATH", "strategy"),
				},
				FTP: FTPConfig{
					Host:         getEnv("FTP_HOST", "127.0.0.1"),
					Port:         getEnvInt("FTP_PORT", 21),
					Username:     getEnv("FTP_USERNAME", ""),
					Password:     getEnv("FTP_PASSWORD", ""),
					LogPath:      getEnv("FTP_LOG_PATH", "logs"),
					StrategyPath: getEnv("FTP_STRATEGY_PATH", "strategy"),
				},
			},
		},
	}

	if globalConfig.DebugLevel == "true" {
		log.Println("配置初始化完成")
		log.Printf("服务器端口: %s\n", globalConfig.ServerPort)
		log.Printf("调试级别: %s\n", globalConfig.DebugLevel)
		log.Printf("主数据库: %s@%s:%s/%s\n", globalConfig.DBUser, globalConfig.DBHost, globalConfig.DBPort, globalConfig.DBName)
		log.Printf("Radius数据库: %s@%s:%s/%s\n", globalConfig.RadiusDBUser, globalConfig.RadiusDBHost, globalConfig.RadiusDBPort, globalConfig.RadiusDBName)
		log.Printf("存储类型: %s\n", globalConfig.ConfigManager.Storage.Type)
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

// getEnvInt 获取环境变量并转换为整数
func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvBool 获取环境变量并转换为布尔值
func getEnvBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
