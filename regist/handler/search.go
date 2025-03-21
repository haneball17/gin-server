package handler

import (
	"log"
	"net/http"
	"strconv"

	"gin-server/config"
	"gin-server/regist/model"

	"github.com/gin-gonic/gin"
)

// GetUserByID 处理根据ID查询用户的请求
func GetUserByID(c *gin.Context) {
	cfg := config.GetConfig() // 获取全局配置

	// 获取查询参数中的用户ID
	userIDStr := c.Query("id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少必要的id参数"})
		return
	}

	// 将用户ID从字符串转换为整数
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("无效的用户ID格式: %s, %v\n", userIDStr, err)
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID格式"})
		return
	}

	// 查询用户信息
	user, err := model.GetUserByID(userID)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("查询用户失败: %v\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户失败"})
		return
	}

	// 如果用户不存在
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	if cfg.DebugLevel == "true" {
		log.Printf("成功查询到用户ID为 %d 的信息\n", userID)
	}

	// 返回用户信息
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "用户查询成功",
		"data":    user,
	})
}

// GetDeviceByID 处理根据ID查询设备的请求
func GetDeviceByID(c *gin.Context) {
	cfg := config.GetConfig() // 获取全局配置

	// 获取查询参数中的设备ID
	deviceID := c.Query("id")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少必要的id参数"})
		return
	}

	// 查询设备信息
	device, err := model.GetDeviceByID(deviceID)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("查询设备失败: %v\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询设备失败"})
		return
	}

	// 如果设备不存在
	if device == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "设备不存在"})
		return
	}

	if cfg.DebugLevel == "true" {
		log.Printf("成功查询到设备ID为 %s 的信息\n", deviceID)
	}

	// 返回设备信息
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "设备查询成功",
		"data":    device,
	})
}
