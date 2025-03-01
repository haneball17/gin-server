package handler

import (
	"log"
	"net/http"

	"gin-server/auth/model"
	configModel "gin-server/regist/model"

	"github.com/gin-gonic/gin"
)

// GetAuthRecords 处理获取认证记录的请求
func GetAuthRecords(c *gin.Context) {
	config := configModel.LoadConfig()
	if config.DebugLevel == "true" {
		log.Println("接收到获取认证记录的请求")
	}

	var query model.AuthRecordQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		if config.DebugLevel == "true" {
			log.Printf("绑定查询参数失败: %v\n", err)
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的查询参数",
			"error":   err.Error(),
		})
		return
	}

	// 设置默认分页参数
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}
	// 限制每页最大记录数
	if query.PageSize > 100 {
		query.PageSize = 100
	}

	if config.DebugLevel == "true" {
		log.Printf("查询参数: %+v\n", query)
	}

	// 查询认证记录
	records, total, err := model.GetAuthRecords(query)
	if err != nil {
		if config.DebugLevel == "true" {
			log.Printf("查询认证记录失败: %v\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询认证记录失败",
			"error":   err.Error(),
		})
		return
	}

	// 计算总页数
	totalPages := (total + query.PageSize - 1) / query.PageSize

	if config.DebugLevel == "true" {
		log.Printf("成功获取认证记录，总记录数: %d, 总页数: %d\n", total, totalPages)
	}

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
