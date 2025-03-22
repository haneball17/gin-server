package main

import (
	stdlog "log"

	authModel "gin-server/auth/model"
	authRouter "gin-server/auth/router"
	"gin-server/config"
	"gin-server/configmanager/common/crypto"
	"gin-server/configmanager/log"
	"gin-server/database"
	"gin-server/regist/router"

	"github.com/gin-gonic/gin"
)

// initKeyManager 初始化密钥管理器
func initKeyManager(cfg *config.Config) error {
	keyManager := crypto.NewKeyManager(cfg)
	if err := keyManager.EnsureKeyPair(); err != nil {
		stdlog.Printf("初始化密钥对失败: %v", err)
		return err
	}

	stdlog.Println("密钥对初始化成功")
	return nil
}

// initLogManager 初始化并启动日志管理器
func initLogManager(cfg *config.Config) (*log.LogManager, error) {
	// 获取主数据库连接
	db, err := database.GetDB()
	if err != nil {
		return nil, err
	}

	// 初始化日志管理器
	logManager, err := log.NewLogManager(cfg, db)
	if err != nil {
		return nil, err
	}

	// 启动日志管理器
	if err := logManager.Start(); err != nil {
		return nil, err
	}

	stdlog.Println("日志管理器启动成功")
	return logManager, nil
}

// setupRoutes 设置所有API路由
func setupRoutes(r *gin.Engine, logManager *log.LogManager) {
	// 设置注册模块路由
	router.SetupRouter(r)

	// 设置认证管理模块路由
	authRouter.SetupRouter(r)

	// 注册日志路由
	if logManager != nil {
		log.RegisterRoutes(r, logManager)
	}
}

func main() {
	// 初始化全局配置
	config.InitConfig()
	cfg := config.GetConfig()

	// 初始化密钥管理器 (非致命错误，允许继续)
	if err := initKeyManager(cfg); err != nil {
		stdlog.Printf("警告: 密钥管理器初始化失败，部分加密功能可能不可用")
	}

	// 初始化数据库连接管理器（致命错误，必须停止）
	if err := database.Initialize(cfg); err != nil {
		stdlog.Fatalf("数据库初始化失败: %v", err)
	}
	defer database.CloseAll() // 程序结束时关闭所有数据库连接

	// 初始化Radius数据库（非致命错误，允许继续）
	if err := authModel.InitRadiusDB(); err != nil {
		stdlog.Printf("警告: Radius数据库初始化失败，认证功能可能不可用: %v", err)
	}

	// 初始化日志管理器（非致命错误，允许继续）
	var logManager *log.LogManager
	if manager, err := initLogManager(cfg); err != nil {
		stdlog.Printf("警告: 日志管理器初始化失败，日志记录功能将不可用: %v", err)
	} else {
		logManager = manager
		// 程序退出时停止日志管理器
		defer func() {
			if err := logManager.Stop(); err != nil {
				stdlog.Printf("停止日志管理器失败: %v", err)
			} else {
				stdlog.Println("日志管理器已停止")
			}
		}()
	}

	// 创建Gin路由引擎
	r := gin.Default()

	// 设置所有路由
	setupRoutes(r, logManager)

	// 启动服务
	stdlog.Printf("服务器启动，监听端口: %s\n", cfg.ServerPort)
	r.Run(":" + cfg.ServerPort)
}
