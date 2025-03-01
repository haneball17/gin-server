package main

import (
	"gin-server/auth/model" as authModel
	"gin-server/auth/router" as authRouter
	"gin-server/config"
	"gin-server/regist/model"
	"gin-server/regist/router"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化全局配置
	config.InitConfig()
	cfg := config.GetConfig()
	
	// 初始化数据库连接
	model.InitDB() // 初始化主数据库连接
	
	// 初始化Radius数据库连接
	if err := authModel.InitRadiusDB(); err != nil {
		log.Fatalf("初始化Radius数据库失败: %v", err)
	}
	
	// 创建Gin路由引擎
	r := gin.Default() // 创建一个默认的 Gin 路由引擎
	
	// 设置路由
	router.SetupRouter(r)     // 设置注册模块路由
	authRouter.SetupRouter(r) // 设置认证管理模块路由
	
	// 启动服务
	log.Printf("服务器启动，监听端口: %s\n", cfg.ServerPort)
	r.Run(":" + cfg.ServerPort) // 启动服务，监听配置的端口
}
