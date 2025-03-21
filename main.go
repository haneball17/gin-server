package main

import (
	stdlog "log"

	authModel "gin-server/auth/model"
	authRouter "gin-server/auth/router"
	"gin-server/config"
	"gin-server/configmanager/common/crypto"
	"gin-server/configmanager/log"
	logModel "gin-server/configmanager/log/model"
	"gin-server/regist/model"
	"gin-server/regist/router"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 初始化全局配置
	config.InitConfig()
	cfg := config.GetConfig()

	// 初始化密钥管理器
	keyManager := crypto.NewKeyManager(cfg)
	if err := keyManager.EnsureKeyPair(); err != nil {
		stdlog.Printf("初始化密钥对失败: %v", err)
	} else {
		stdlog.Println("密钥对初始化成功")
	}

	// 初始化数据库连接
	model.InitDB() // 初始化主数据库连接

	// 初始化Radius数据库连接
	if err := authModel.InitRadiusDB(); err != nil {
		stdlog.Fatalf("初始化Radius数据库失败: %v", err)
	}

	// 初始化GORM数据库连接
	dsn := cfg.DBUser + ":" + cfg.DBPassword + "@tcp(" + cfg.DBHost + ":" + cfg.DBPort + ")/" + cfg.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		stdlog.Fatalf("连接GORM数据库失败: %v", err)
	}

	// 检查并迁移数据库表结构
	if err := logModel.MigrateDatabase(db); err != nil {
		stdlog.Printf("检查并迁移数据库表结构失败: %v", err)
	} else {
		stdlog.Println("数据库表结构检查和迁移完成")
	}

	// 初始化日志管理器
	logManager, err := log.NewLogManager(cfg, db)
	if err != nil {
		stdlog.Printf("创建日志管理器失败: %v", err)
	} else {
		// 启动日志管理器
		if err := logManager.Start(); err != nil {
			stdlog.Printf("启动日志管理器失败: %v", err)
		} else {
			stdlog.Println("日志管理器启动成功")
			// 程序退出时停止日志管理器
			defer func() {
				if err := logManager.Stop(); err != nil {
					stdlog.Printf("停止日志管理器失败: %v", err)
				} else {
					stdlog.Println("日志管理器已停止")
				}
			}()
		}
	}

	// 创建Gin路由引擎
	r := gin.Default()

	// 设置路由
	router.SetupRouter(r)     // 设置注册模块路由
	authRouter.SetupRouter(r) // 设置认证管理模块路由

	// 注册日志路由
	log.RegisterRoutes(r, logManager)

	// 启动服务
	stdlog.Printf("服务器启动，监听端口: %s\n", cfg.ServerPort)
	r.Run(":" + cfg.ServerPort)
}
