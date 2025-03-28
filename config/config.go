package config

import (
	"log"
	"os"
	"strconv"
)

// Config 系统全局配置结构体
type Config struct {
	// DebugLevel 调试级别
	// 可选值: "true" - 输出详细调试日志, "false" - 仅输出基本日志
	DebugLevel string

	// Debug 是否开启调试模式
	// 控制数据库日志输出等详细信息
	Debug bool

	// ServerPort 服务器监听端口
	// 格式: "端口号", 例如: "8080"
	ServerPort string

	// DBUser 主数据库用户名
	// 用于访问系统主数据库的用户名
	DBUser string

	// DBPassword 主数据库密码
	// 用于访问系统主数据库的密码
	DBPassword string

	// DBHost 主数据库主机地址
	// 格式: "主机名或IP地址"
	DBHost string

	// DBPort 主数据库端口
	// 格式: "端口号"
	DBPort string

	// DBName 主数据库名称
	// 系统主数据库的名称
	DBName string

	// RadiusDBUser Radius数据库用户名
	// 用于访问Radius认证数据库的用户名
	RadiusDBUser string

	// RadiusDBPassword Radius数据库密码
	// 用于访问Radius认证数据库的密码
	RadiusDBPassword string

	// RadiusDBHost Radius数据库主机地址
	// 格式: "主机名或IP地址"
	RadiusDBHost string

	// RadiusDBPort Radius数据库端口
	// 格式: "端口号"
	RadiusDBPort string

	// RadiusDBName Radius数据库名称
	// Radius认证数据库的名称
	RadiusDBName string

	// ConfigManager 配置管理模块配置
	// 包含日志管理、策略管理和存储配置
	ConfigManager ConfigManagerConfig

	// Gitee Gitee仓库配置
	// 用于远程存储和版本控制
	Gitee *GiteeConfig `json:"gitee" yaml:"gitee"`

	// FTP FTP服务器配置
	// 用于远程文件传输
	FTP *FTPConfig `json:"ftp" yaml:"ftp"`

	// TestData 测试数据配置
	// 用于控制测试数据的生成和维护
	TestData TestDataConfig
}

// ConfigManagerConfig 配置管理模块配置结构体
type ConfigManagerConfig struct {
	// LogManager 日志管理配置
	// 控制日志生成和加密
	LogManager LogManagerConfig

	// StrategyManager 策略管理配置
	// 控制策略更新和加密
	StrategyManager StrategyManagerConfig

	// Storage 存储配置
	// 配置远程存储方式和参数
	Storage StorageConfig

	// Compress 压缩配置
	// 配置文件压缩和解压参数
	Compress CompressConfig
}

// LogManagerConfig 日志管理配置结构体
type LogManagerConfig struct {
	// GenerateInterval 日志生成间隔（分钟）
	// 指定系统生成新日志文件的时间间隔
	GenerateInterval int

	// EnableEncryption 是否启用加密
	EnableEncryption bool `yaml:"enable_encryption"`

	// LogDir 日志目录
	LogDir string `yaml:"log_dir"`

	// UploadDir 远程上传目录
	// 指定日志文件上传到远程存储的目录路径
	UploadDir string `yaml:"upload_dir"`

	// Encryption 加密配置
	Encryption EncryptionConfig `yaml:"encryption"`

	// ProcessedLogPath 处理后的日志文件路径
	ProcessedLogPath string `json:"-" yaml:"-"`

	// ProcessedKeyPath 处理后的密钥文件路径
	ProcessedKeyPath string `json:"-" yaml:"-"`
}

// StrategyManagerConfig 策略管理配置结构体
type StrategyManagerConfig struct {
	// PollInterval 轮询间隔（秒）
	// 指定检查策略更新的时间间隔
	PollInterval int

	// EnableEncryption 是否启用加密
	// true: 启用策略加密, false: 不加密
	EnableEncryption bool

	// Encryption 加密配置
	// 指定策略加密的相关参数
	Encryption EncryptionConfig
}

// StorageConfig 存储配置结构体
type StorageConfig struct {
	// Type 存储类型
	// 可选值: "gitee" - 使用Gitee仓库, "ftp" - 使用FTP服务器
	Type string

	// Gitee Gitee配置
	// 当Type为"gitee"时使用
	Gitee GiteeConfig

	// FTP FTP配置
	// 当Type为"ftp"时使用
	FTP FTPConfig
}

// CompressConfig 压缩配置结构体
type CompressConfig struct {
	// Format 压缩格式
	// 可选值: "tar.gz" - tar.gz格式
	Format string

	// CompressionLevel 压缩级别
	// 可选值: 1-9，1最快但压缩率最低，9最慢但压缩率最高
	CompressionLevel int

	// BufferSize 缓冲区大小（字节）
	// 建议值: 32768 (32KB)
	BufferSize int

	// IncludePatterns 包含的文件模式
	// 例如: ["*.json", "*.txt"]
	IncludePatterns []string

	// ExcludePatterns 排除的文件模式
	// 例如: ["*.tmp", "*.bak"]
	ExcludePatterns []string
}

// GiteeConfig Gitee配置结构体
type GiteeConfig struct {
	// AccessToken Gitee访问令牌
	// 用于访问Gitee API的授权令牌
	AccessToken string `json:"access_token" yaml:"access_token"`

	// Owner 仓库所有者
	// Gitee仓库所有者的用户名或组织名
	Owner string `json:"owner" yaml:"owner"`

	// Repo 仓库名称
	// Gitee仓库的名称
	Repo string `json:"repo" yaml:"repo"`

	// Branch 分支名称
	// 要操作的Git分支名称，默认为"master"
	Branch string `json:"branch" yaml:"branch"`
}

// FTPConfig FTP配置结构体
type FTPConfig struct {
	// Host FTP服务器地址
	// 格式: "主机名或IP地址"
	Host string `json:"host" yaml:"host"`

	// Port FTP服务器端口
	// 默认为21
	Port int `json:"port" yaml:"port"`

	// Username FTP用户名
	// 用于FTP服务器认证
	Username string `json:"username" yaml:"username"`

	// Password FTP密码
	// 用于FTP服务器认证
	Password string `json:"password" yaml:"password"`
}

// TestDataConfig 测试数据配置结构体
type TestDataConfig struct {
	// EnableInitialData 是否在启动时初始化测试数据
	// true: 启动时自动添加测试数据, false: 不添加测试数据
	EnableInitialData bool `json:"enable_initial_data" yaml:"enable_initial_data"`

	// EnableRealtimeData 是否启用实时测试数据生成
	// true: 启用实时测试数据生成, false: 不生成实时测试数据
	EnableRealtimeData bool `json:"enable_realtime_data" yaml:"enable_realtime_data"`

	// DeviceCount 生成的设备数量
	// 至少需要生成10条设备数据
	DeviceCount int `json:"device_count" yaml:"device_count"`

	// UsersPerDevice 每个设备的用户数量
	// 每个网关至少需要5个用户
	UsersPerDevice int `json:"users_per_device" yaml:"users_per_device"`

	// BehaviorsPerUser 每个用户的行为数据量
	// 每个用户至少需要10条行为数据
	BehaviorsPerUser int `json:"behaviors_per_user" yaml:"behaviors_per_user"`

	// RealtimeInterval 实时数据生成间隔(秒)
	// 默认60秒，即每分钟生成一次
	RealtimeInterval int `json:"realtime_interval" yaml:"realtime_interval"`

	// RealtimeBehaviorsPerInterval 每个间隔生成的行为数据量
	// 默认10条，即每分钟为每个测试用户添加10条数据
	RealtimeBehaviorsPerInterval int `json:"realtime_behaviors_per_interval" yaml:"realtime_behaviors_per_interval"`

	// RealtimeStartTimeOffset 实时数据生成的起始时间偏移(分钟)
	// 默认值为30，表示从当前时间往前30分钟作为起始时间
	// 负值表示过去的时间，正值表示未来的时间（通常不建议使用正值）
	RealtimeStartTimeOffset int `json:"realtime_start_time_offset" yaml:"realtime_start_time_offset"`

	// RealtimeEndTimeOffset 实时数据生成的结束时间偏移(分钟)
	// 默认值为0，表示直到当前时间
	// 负值表示过去的时间，正值表示未来的时间（通常不建议使用正值）
	RealtimeEndTimeOffset int `json:"realtime_end_time_offset" yaml:"realtime_end_time_offset"`
}

// EncryptionConfig 加密配置结构体
type EncryptionConfig struct {
	// AESKeyLength AES密钥长度
	// 可选值: 128, 192, 256
	AESKeyLength int

	// PublicKeyAlgorithm 公钥加密算法
	// 可选值: "RSA", "ECDSA", "ED25519"
	PublicKeyAlgorithm string

	// PublicKeyLength 公钥长度
	// RSA推荐: 2048或4096
	// ECDSA推荐: 256或384
	PublicKeyLength int

	// PublicKeyPath 运维系统公钥路径
	// 用于加密的公钥文件路径
	PublicKeyPath string

	// PrivateKeyPath 本系统私钥路径
	// 用于解密的私钥文件路径
	PrivateKeyPath string
}

// 全局配置实例
var globalConfig *Config

// InitConfig 初始化配置
func InitConfig() {
	log.Println("初始化全局配置...")

	globalConfig = &Config{
		// 通用配置
		DebugLevel: getEnv("DEBUG_LEVEL", "true"),

		// 服务器配置
		ServerPort: getEnv("SERVER_PORT", "8123"),

		// 主数据库配置
		DBUser:     getEnv("DB_USER", "gin_user"),
		DBPassword: getEnv("DB_PASSWORD", "P@ssw0rd123!"),
		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBName:     getEnv("DB_NAME", "gin_server"),

		// Radius数据库配置 (默认与主数据库相同的连接信息，但不同的数据库名)
		RadiusDBUser:     getEnv("RADIUS_DB_USER", getEnv("DB_USER", "gin_user")),
		RadiusDBPassword: getEnv("RADIUS_DB_PASSWORD", getEnv("DB_PASSWORD", "P@ssw0rd123!")),
		RadiusDBHost:     getEnv("RADIUS_DB_HOST", getEnv("DB_HOST", "127.0.0.1")),
		RadiusDBPort:     getEnv("RADIUS_DB_PORT", getEnv("DB_PORT", "3306")),
		RadiusDBName:     getEnv("RADIUS_DB_NAME", "radius"),

		// 配置管理模块配置
		ConfigManager: ConfigManagerConfig{
			// 日志管理配置
			LogManager: LogManagerConfig{
				GenerateInterval: getEnvInt("LOG_GENERATE_INTERVAL", 10),
				EnableEncryption: getEnvBool("LOG_ENABLE_ENCRYPTION", true),
				LogDir:           getEnv("LOG_DIR", "logs"),
				UploadDir:        getEnv("LOG_UPLOAD_DIR", "log"),
				Encryption: EncryptionConfig{
					AESKeyLength:       getEnvInt("LOG_AES_KEY_LENGTH", 256),
					PublicKeyAlgorithm: getEnv("LOG_PUBLIC_KEY_ALGORITHM", "RSA"),
					PublicKeyLength:    getEnvInt("LOG_PUBLIC_KEY_LENGTH", 2048),
					PublicKeyPath:      getEnv("LOG_PUBLIC_KEY_PATH", "keys/public.pem"),
					PrivateKeyPath:     getEnv("LOG_PRIVATE_KEY_PATH", "keys/private.pem"),
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
					PublicKeyPath:      getEnv("STRATEGY_PUBLIC_KEY_PATH", "keys/public.pem"),
					PrivateKeyPath:     getEnv("STRATEGY_PRIVATE_KEY_PATH", "keys/private.pem"),
				},
			},
			// 存储配置
			Storage: StorageConfig{
				Type: getEnv("STORAGE_TYPE", "gitee"),
				Gitee: GiteeConfig{
					AccessToken: getEnv("GITEE_ACCESS_TOKEN", "684ed853186d45b825b415b00f16bf27"),
					Owner:       getEnv("GITEE_OWNER", "xxxLogTest"),
					Repo:        getEnv("GITEE_REPO", "logtest"),
					Branch:      getEnv("GITEE_BRANCH", "master"),
				},
				FTP: FTPConfig{
					Host:     getEnv("FTP_HOST", "127.0.0.1"),
					Port:     getEnvInt("FTP_PORT", 21),
					Username: getEnv("FTP_USERNAME", ""),
					Password: getEnv("FTP_PASSWORD", ""),
				},
			},
			Compress: CompressConfig{
				Format:           getEnv("COMPRESS_FORMAT", "tar.gz"),
				CompressionLevel: getEnvInt("COMPRESS_LEVEL", 6),
				BufferSize:       getEnvInt("COMPRESS_BUFFER_SIZE", 32*1024),
				IncludePatterns:  []string{},
				ExcludePatterns:  []string{},
			},
		},

		// 测试数据配置
		TestData: TestDataConfig{
			EnableInitialData:            getEnvBool("TEST_ENABLE_INITIAL_DATA", true),
			EnableRealtimeData:           getEnvBool("TEST_ENABLE_REALTIME_DATA", false),
			DeviceCount:                  getEnvInt("TEST_DEVICE_COUNT", 10),
			UsersPerDevice:               getEnvInt("TEST_USERS_PER_DEVICE", 5),
			BehaviorsPerUser:             getEnvInt("TEST_BEHAVIORS_PER_USER", 10),
			RealtimeInterval:             getEnvInt("TEST_REALTIME_INTERVAL", 60),
			RealtimeBehaviorsPerInterval: getEnvInt("TEST_REALTIME_BEHAVIORS_PER_INTERVAL", 10),
			RealtimeStartTimeOffset:      getEnvInt("TEST_REALTIME_START_TIME_OFFSET", 2),
			RealtimeEndTimeOffset:        getEnvInt("TEST_REALTIME_END_TIME_OFFSET", 0),
		},
	}

	// 设置Gitee配置
	if globalConfig.ConfigManager.Storage.Type == "gitee" {
		globalConfig.Gitee = &GiteeConfig{
			AccessToken: getEnv("GITEE_ACCESS_TOKEN", "684ed853186d45b825b415b00f16bf27"),
			Owner:       getEnv("GITEE_OWNER", "xxxLogTest"),
			Repo:        getEnv("GITEE_REPO", "logtest"),
			Branch:      getEnv("GITEE_BRANCH", "master"),
		}
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

// SetConfig 设置全局配置
func SetConfig(cfg *Config) {
	globalConfig = cfg
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

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		DebugLevel: "false",
		Debug:      false,
		ServerPort: "8123",
		DBUser:     "root",
		DBPassword: "",
		DBHost:     "localhost",
		DBPort:     "3306",
		DBName:     "gin_server",
		ConfigManager: ConfigManagerConfig{
			LogManager: LogManagerConfig{
				GenerateInterval: 1,
				EnableEncryption: true,
				LogDir:           "logs",
				UploadDir:        "log",
				Encryption: EncryptionConfig{
					AESKeyLength:       256,
					PublicKeyAlgorithm: "RSA",
					PublicKeyLength:    2048,
					PublicKeyPath:      "keys/public.pem",
					PrivateKeyPath:     "keys/private.pem",
				},
			},
			StrategyManager: StrategyManagerConfig{
				PollInterval:     60,
				EnableEncryption: false,
				Encryption: EncryptionConfig{
					AESKeyLength:       256,
					PublicKeyAlgorithm: "RSA",
					PublicKeyLength:    2048,
					PublicKeyPath:      "keys/ops_public.pem",
					PrivateKeyPath:     "keys/system_private.pem",
				},
			},
			Storage: StorageConfig{
				Type: "gitee",
				Gitee: GiteeConfig{
					AccessToken: "",
					Owner:       "",
					Repo:        "",
					Branch:      "master",
				},
				FTP: FTPConfig{
					Host:     "127.0.0.1",
					Port:     21,
					Username: "",
					Password: "",
				},
			},
			Compress: CompressConfig{
				Format:           "tar.gz",
				CompressionLevel: 6,
				BufferSize:       32 * 1024,
				IncludePatterns:  []string{},
				ExcludePatterns:  []string{},
			},
		},
		TestData: TestDataConfig{
			EnableInitialData:            true,
			EnableRealtimeData:           false,
			DeviceCount:                  10,
			UsersPerDevice:               5,
			BehaviorsPerUser:             10,
			RealtimeInterval:             60,
			RealtimeBehaviorsPerInterval: 10,
			RealtimeStartTimeOffset:      2,
			RealtimeEndTimeOffset:        0,
		},
	}
}
