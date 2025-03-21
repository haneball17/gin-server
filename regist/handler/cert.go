package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"gin-server/config"
	"gin-server/regist/model"

	"github.com/gin-gonic/gin"
)

// 最大文件大小（8MB）
const MaxFileSize = 8 * 1024 * 1024

// BindUserCert 处理用户证书绑定
func BindUserCert(c *gin.Context) {
	cfg := config.GetConfig()
	userID := c.Param("id")

	// 检查用户ID是否有效
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	// 检查用户是否存在
	exists, err := model.CheckUserExistsByID(userIDInt)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("检查用户是否存在失败: %v\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查用户是否存在失败"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("cert")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未找到证书文件"})
		return
	}

	// 检查文件大小
	if file.Size > MaxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件大小超过限制"})
		return
	}

	// 检查文件扩展名
	if file.Filename[len(file.Filename)-4:] != ".pem" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件必须是.pem格式"})
		return
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("打开上传的文件失败: %v\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法读取上传的文件"})
		return
	}
	defer src.Close()

	// 保存文件
	filePath, err := model.SaveCertFile("user", userID, src, false)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("保存证书文件失败: %v\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存证书文件失败"})
		return
	}

	// 添加记录到证书表
	err = model.AddCertRecord("user", userID, filePath, "")
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("添加证书记录失败: %v\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加证书记录失败"})
		return
	}

	// 更新用户表中的证书信息
	err = model.UpdateUserCertInfo(userID, filePath, "")
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("更新用户证书信息失败: %v\n", err)
		}
		// 不中断请求，因为证书已经保存成功
		log.Printf("警告: 更新用户证书信息失败: %v\n", err)
	}

	// 返回成功信息
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "证书绑定成功",
		"data": gin.H{
			"userID":   userID,
			"certPath": filePath,
		},
	})
}

// BindUserKey 处理用户密钥绑定
func BindUserKey(c *gin.Context) {
	cfg := config.GetConfig()
	userID := c.Param("id")

	// 检查用户ID是否有效
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	// 检查用户是否存在
	exists, err := model.CheckUserExistsByID(userIDInt)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("检查用户是否存在失败: %v\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查用户是否存在失败"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("key")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未找到密钥文件"})
		return
	}

	// 检查文件大小
	if file.Size > MaxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件大小超过限制"})
		return
	}

	// 检查文件扩展名
	if file.Filename[len(file.Filename)-4:] != ".pem" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件必须是.pem格式"})
		return
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("打开上传的文件失败: %v\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法读取上传的文件"})
		return
	}
	defer src.Close()

	// 保存文件
	filePath, err := model.SaveCertFile("user", userID, src, true)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("保存密钥文件失败: %v\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存密钥文件失败"})
		return
	}

	// 获取当前证书信息
	cert, err := model.GetCertInfo("user", userID)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("获取证书信息失败: %v\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取证书信息失败"})
		return
	}

	// 添加或更新记录到证书表
	certPath := ""
	if cert != nil {
		certPath = cert.CertPath
	}
	err = model.AddCertRecord("user", userID, certPath, filePath)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("添加密钥记录失败: %v\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加密钥记录失败"})
		return
	}

	// 更新用户表中的密钥信息
	err = model.UpdateUserCertInfo(userID, "", filePath)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("更新用户密钥信息失败: %v\n", err)
		}
		// 不中断请求，因为密钥已经保存成功
		log.Printf("警告: 更新用户密钥信息失败: %v\n", err)
	}

	// 返回成功信息
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "密钥绑定成功",
		"data": gin.H{
			"userID":  userID,
			"keyPath": filePath,
		},
	})
}

// BindDeviceCert 处理设备证书绑定
func BindDeviceCert(c *gin.Context) {
	cfg := config.GetConfig()
	deviceID := c.Param("id")

	// 检查设备是否存在
	exists, err := model.CheckDeviceExistsByID(deviceID)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("检查设备是否存在失败: %v\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查设备是否存在失败"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "设备不存在"})
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("cert")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未找到证书文件"})
		return
	}

	// 检查文件大小
	if file.Size > MaxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件大小超过限制"})
		return
	}

	// 检查文件扩展名
	if file.Filename[len(file.Filename)-4:] != ".pem" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件必须是.pem格式"})
		return
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("打开上传的文件失败: %v\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法读取上传的文件"})
		return
	}
	defer src.Close()

	// 保存文件
	filePath, err := model.SaveCertFile("device", deviceID, src, false)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("保存证书文件失败: %v\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存证书文件失败"})
		return
	}

	// 添加记录到证书表
	err = model.AddCertRecord("device", deviceID, filePath, "")
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("添加证书记录失败: %v\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加证书记录失败"})
		return
	}

	// 更新设备表中的证书信息
	err = model.UpdateDeviceCertInfo(deviceID, filePath, "")
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("更新设备证书信息失败: %v\n", err)
		}
		// 不中断请求，因为证书已经保存成功
		log.Printf("警告: 更新设备证书信息失败: %v\n", err)
	}

	// 返回成功信息
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "证书绑定成功",
		"data": gin.H{
			"deviceID": deviceID,
			"certPath": filePath,
		},
	})
}

// BindDeviceKey 处理设备密钥绑定
func BindDeviceKey(c *gin.Context) {
	cfg := config.GetConfig()
	deviceID := c.Param("id")

	// 检查设备是否存在
	exists, err := model.CheckDeviceExistsByID(deviceID)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("检查设备是否存在失败: %v\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查设备是否存在失败"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "设备不存在"})
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("key")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未找到密钥文件"})
		return
	}

	// 检查文件大小
	if file.Size > MaxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件大小超过限制"})
		return
	}

	// 检查文件扩展名
	if file.Filename[len(file.Filename)-4:] != ".pem" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件必须是.pem格式"})
		return
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("打开上传的文件失败: %v\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法读取上传的文件"})
		return
	}
	defer src.Close()

	// 保存文件
	filePath, err := model.SaveCertFile("device", deviceID, src, true)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("保存密钥文件失败: %v\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存密钥文件失败"})
		return
	}

	// 获取当前证书信息
	cert, err := model.GetCertInfo("device", deviceID)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("获取证书信息失败: %v\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取证书信息失败"})
		return
	}

	// 添加或更新记录到证书表
	certPath := ""
	if cert != nil {
		certPath = cert.CertPath
	}
	err = model.AddCertRecord("device", deviceID, certPath, filePath)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("添加密钥记录失败: %v\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加密钥记录失败"})
		return
	}

	// 更新设备表中的密钥信息
	err = model.UpdateDeviceCertInfo(deviceID, "", filePath)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("更新设备密钥信息失败: %v\n", err)
		}
		// 不中断请求，因为密钥已经保存成功
		log.Printf("警告: 更新设备密钥信息失败: %v\n", err)
	}

	// 返回成功信息
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "密钥绑定成功",
		"data": gin.H{
			"deviceID": deviceID,
			"keyPath":  filePath,
		},
	})
}

// GetCertInfo 获取证书信息
func GetCertInfo(c *gin.Context) {
	cfg := config.GetConfig()
	entityType := c.Query("type")
	entityID := c.Query("id")

	// 检查参数
	if entityType == "" || entityID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少必要参数"})
		return
	}

	// 检查实体类型
	if entityType != "user" && entityType != "device" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的实体类型"})
		return
	}

	// 检查实体是否存在
	var exists bool
	var err error
	if entityType == "user" {
		userID, err := strconv.Atoi(entityID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
			return
		}
		exists, err = model.CheckUserExistsByID(userID)
	} else {
		exists, err = model.CheckDeviceExistsByID(entityID)
	}

	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("检查实体是否存在失败: %v\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("检查%s是否存在失败", entityType)})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("%s不存在", entityType)})
		return
	}

	// 获取证书信息
	cert, err := model.GetCertInfo(entityType, entityID)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("获取证书信息失败: %v\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取证书信息失败"})
		return
	}

	if cert == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到证书信息"})
		return
	}

	// 返回证书信息
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取证书信息成功",
		"data":    cert,
	})
}
