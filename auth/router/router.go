package router

import (
	"gin-server/auth/handler" // 导入处理器包
	"gin-server/config"
	"log"

	"github.com/gin-gonic/gin" // 导入 Gin 框架
)

// SetupRouter 设置认证管理模块的路由
func SetupRouter(r *gin.Engine) {
	cfg := config.GetConfig()

	// 添加中间件记录请求日志
	authGroup := r.Group("/auth")

	// 仅在调试模式下启用日志中间件
	if cfg.DebugLevel == "true" {
		authGroup.Use(func(c *gin.Context) {
			log.Printf("认证管理 - 请求路径: %s, 方法: %s\n", c.Request.URL.Path, c.Request.Method)
			c.Next()
		})
	}

	// 认证记录查询接口
	authGroup.GET("/records", handler.GetAuthRecords)
}
