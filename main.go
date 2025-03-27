package main

import (
	"fmt"
	"io"
	stdlog "log"
	"net/http"
	"os"
	"runtime"
	"time"

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
	// 添加自定义的恢复中间件
	r.Use(customRecoveryMiddleware())

	// 设置注册模块路由
	router.SetupRouter(r)

	// 设置认证管理模块路由
	authRouter.SetupRouter(r)

	// 注册日志路由
	if logManager != nil {
		log.RegisterRoutes(r, logManager)
	}
}

// customRecoveryMiddleware 自定义恢复中间件，捕获请求处理过程中的panic
func customRecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 获取堆栈信息
				buf := make([]byte, 2048)
				n := runtime.Stack(buf, false)
				stackInfo := fmt.Sprintf("%s", buf[:n])

				// 记录错误和堆栈信息
				stdlog.Printf("[严重错误] 处理请求时发生panic: %v\n堆栈: %s", err, stackInfo)

				// 向客户端返回500错误
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "服务器内部错误",
					"error":   fmt.Sprintf("%v", err),
				})
			}
		}()
		c.Next()
	}
}

func main() {
	// 设置日志记录到文件
	setupFileLogger()

	// 初始化全局配置
	config.InitConfig()
	cfg := config.GetConfig()

	// 初始化密钥管理器 (非致命错误，允许继续)
	if err := initKeyManager(cfg); err != nil {
		stdlog.Printf("警告: 密钥管理器初始化失败，部分加密功能可能不可用")
	}

	// 初始化数据库连接管理器（致命错误，必须停止）
	if err := database.Initialize(cfg); err != nil {
		stdlog.Printf("数据库初始化失败: %v", err)
		waitForExit() // 在错误退出前等待用户输入
		stdlog.Fatalf("程序终止")
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

	// 设置恢复函数，防止服务器因panic而崩溃
	defer func() {
		if err := recover(); err != nil {
			stdlog.Printf("服务器发生严重错误: %v", err)
			waitForExit() // 在panic发生后等待用户输入
		}
	}()

	if err := r.Run(":" + cfg.ServerPort); err != nil {
		stdlog.Printf("服务器启动失败: %v", err)
		waitForExit() // 在服务器启动失败后等待用户输入
	}
}

// waitForExit 等待用户输入后再退出程序
func waitForExit() {
	stdlog.Println("按任意键退出程序...")
	var input string
	fmt.Scanln(&input)
}

// setupFileLogger 设置文件日志记录器
func setupFileLogger() {
	// 确保logs目录存在
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", 0755)
	}

	// 创建日志文件，使用当前日期作为文件名
	currentTime := time.Now()
	logFileName := fmt.Sprintf("logs/server_%s.log", currentTime.Format("2006-01-02"))
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		stdlog.Printf("无法创建日志文件: %v", err)
		return
	}

	// 创建多重写入器，同时写入文件和标准输出
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	stdlog.SetOutput(multiWriter)

	// 设置日志格式
	stdlog.SetFlags(stdlog.Ldate | stdlog.Ltime | stdlog.Lshortfile)

	stdlog.Printf("日志系统初始化完成，日志将同时记录到控制台和文件: %s", logFileName)
}
