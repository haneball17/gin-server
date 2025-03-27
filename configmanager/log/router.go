package log

import (
	"net/http"
	"strconv"
	"time"

	"gin-server/database/models"

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
			c.Data(http.StatusOK, "application/json", content)
		})

		// 生成日志 "/logs/generate"
		logGroup.POST("/generate", func(c *gin.Context) {
			if err := logManager.GenerateLog(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"message": "日志生成成功",
			})
		})

		// 获取远程日志文件列表 "/logs/files"
		logGroup.GET("/files", func(c *gin.Context) {
			files, err := logManager.ListRemoteLogFiles()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, files)
		})

		// 根据时间范围查询日志文件 "/logs/files/search"
		logGroup.GET("/files/search", func(c *gin.Context) {
			startTimeStr := c.Query("start_time")
			endTimeStr := c.Query("end_time")

			var startTime, endTime time.Time
			var err error

			if startTimeStr != "" {
				startTime, err = time.Parse(time.RFC3339, startTimeStr)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": "无效的开始时间格式，请使用RFC3339格式",
					})
					return
				}
			} else {
				// 默认为24小时前
				startTime = time.Now().Add(-24 * time.Hour)
			}

			if endTimeStr != "" {
				endTime, err = time.Parse(time.RFC3339, endTimeStr)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": "无效的结束时间格式，请使用RFC3339格式",
					})
					return
				}
			} else {
				// 默认为当前时间
				endTime = time.Now()
			}

			logFiles, count, err := logManager.GetLogFilesByTimeRange(startTime, endTime)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"total": count,
				"files": logFiles,
			})
		})

		// 创建事件 "/logs/events"
		logGroup.POST("/events", func(c *gin.Context) {
			type EventRequest struct {
				EventCode string `json:"event_code" binding:"required"`
				EventDesc string `json:"event_desc" binding:"required"`
				DeviceID  int    `json:"device_id"`
				EventType int    `json:"event_type" binding:"required"`
			}

			var req EventRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// 将整数类型转换为EventType枚举
			eventType := models.EventType(req.EventType)

			// 验证事件类型
			if eventType != models.EventTypeSecurity && eventType != models.EventTypeFault {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "无效的事件类型，有效值: 1 (安全事件), 2 (故障事件)",
				})
				return
			}

			event, err := logManager.CreateEvent(req.EventCode, req.EventDesc, req.DeviceID, eventType)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusCreated, event)
		})

		// 查询事件 "/logs/events/search"
		logGroup.GET("/events/search", func(c *gin.Context) {
			startTimeStr := c.Query("start_time")
			endTimeStr := c.Query("end_time")

			var startTime, endTime time.Time
			var err error

			if startTimeStr != "" {
				startTime, err = time.Parse(time.RFC3339, startTimeStr)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": "无效的开始时间格式，请使用RFC3339格式",
					})
					return
				}
			} else {
				// 默认为1小时前
				startTime = time.Now().Add(-1 * time.Hour)
			}

			if endTimeStr != "" {
				endTime, err = time.Parse(time.RFC3339, endTimeStr)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": "无效的结束时间格式，请使用RFC3339格式",
					})
					return
				}
			} else {
				// 默认为当前时间
				endTime = time.Now()
			}

			events, count, err := logManager.GetEventsByTimeRange(startTime, endTime)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"total":  count,
				"events": events,
			})
		})

		// 记录用户行为 "/logs/behaviors"
		logGroup.POST("/behaviors", func(c *gin.Context) {
			type BehaviorRequest struct {
				UserID       int   `json:"user_id" binding:"required"`
				BehaviorType int   `json:"behavior_type" binding:"required"`
				DataType     int   `json:"data_type"`
				DataSize     int64 `json:"data_size"`
			}

			var req BehaviorRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// 验证行为类型
			if req.BehaviorType < 1 || req.BehaviorType > 2 {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "无效的行为类型，有效值: 1 (发送), 2 (接收)",
				})
				return
			}

			behavior, err := logManager.LogUserBehavior(req.UserID, req.BehaviorType, req.DataType, req.DataSize)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusCreated, behavior)
		})

		// 查询用户行为 "/logs/behaviors/:user_id"
		logGroup.GET("/behaviors/:user_id", func(c *gin.Context) {
			userIDStr := c.Param("user_id")
			userID, err := strconv.Atoi(userIDStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "无效的用户ID",
				})
				return
			}

			behaviors, count, err := logManager.GetUserBehaviorsByUserID(userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"total":     count,
				"behaviors": behaviors,
				"user_id":   userID,
			})
		})
	}
}
