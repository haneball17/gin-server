package handler

import (
	"log"
	"net/http"

	"gin-server/config"
	"gin-server/database"
	"gin-server/database/models"
	"gin-server/database/repositories"

	"github.com/gin-gonic/gin"
)

// GetAuthRecords 处理获取认证记录的请求
func GetAuthRecords(c *gin.Context) {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Println("接收到获取认证记录的请求")
	}

	var query models.RadPostAuthQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的查询参数",
			"error":   err.Error(),
		})
		return
	}

	// 验证分页参数
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}
	if query.PageSize > 100 {
		query.PageSize = 100
	}

	// 获取Radius数据库连接
	radiusDB, err := database.GetRadiusDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "数据库连接失败",
			"error":   err.Error(),
		})
		return
	}

	// 查询认证记录
	factory := repositories.NewRepositoryFactory(radiusDB)
	authRepo := factory.GetRadiusAuthRepository()
	records, total, err := authRepo.FindByConditions(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询认证记录失败",
			"error":   err.Error(),
		})
		return
	}

	// 计算总页数
	totalPages := (total + int64(query.PageSize) - 1) / int64(query.PageSize)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data": gin.H{
			"total":       total,
			"total_pages": totalPages,
			"page":        query.Page,
			"page_size":   query.PageSize,
			"records":     records,
		},
	})
}
