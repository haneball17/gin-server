package log

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册日志相关路由
func RegisterRoutes(router *gin.Engine, logManager *LogManager) {
	// 直接使用 "/logs" 作为路由组前缀
	logGroup := router.Group("/logs")
	{
		// 最终路径将是 "/logs/latest"
		logGroup.GET("/latest", func(c *gin.Context) {
			content, err := logManager.GetLatestLogContent()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			// 直接返回JSON内容
			c.Data(http.StatusOK, "application/json", content)
		})
	}
}
