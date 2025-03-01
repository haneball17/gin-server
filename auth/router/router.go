package router

import (
	"gin-server/auth/handler" // 导入处理器包
	"log"

	"github.com/gin-gonic/gin" // 导入 Gin 框架
)

// SetupRouter 设置认证管理模块的路由
func SetupRouter(r *gin.Engine) {
	// 添加中间件记录请求日志
	authGroup := r.Group("/auth")
	authGroup.Use(func(c *gin.Context) {
		log.Printf("认证管理 - 请求路径: %s, 方法: %s\n", c.Request.URL.Path, c.Request.Method) // 记录请求路径和方法
		c.Next()                                                                      // 继续处理请求
		log.Printf("认证管理 - 响应状态: %d\n", c.Writer.Status())                            // 记录响应状态
	})

	// 认证记录查询接口
	authGroup.GET("/records", handler.GetAuthRecords) // 获取认证记录接口
}
