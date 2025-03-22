package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"gin-server/config"
	"gin-server/database"
	"gin-server/database/repositories"

	// 保留model包在第二阶段，但计划在第三阶段完全移除它

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 最大文件大小（8MB）
const MaxFileSize = 8 * 1024 * 1024

// 证书文件存储的基础路径
const (
	BaseCertPath = "regist/certs"
	CertsDir     = "certs"
	KeysDir      = "keys"
)

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

	// 获取数据库连接和仓库
	db, err := database.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库连接失败"})
		return
	}
	repoFactory := repositories.NewRepositoryFactory(db)
	userRepo := repoFactory.GetUserRepository()

	// 检查用户是否存在
	user, err := userRepo.FindByUserID(userIDInt)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("检查用户是否存在失败: %v\n", err)
		}
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "检查用户是否存在失败"})
		}
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
	filePath, err := saveFileToDisk("user", userID, src, false)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("保存证书文件失败: %v\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存证书文件失败"})
		return
	}

	// 更新证书记录
	certRepo := repoFactory.GetCertRepository()
	err = certRepo.UpdateCertPath("user", userID, filePath)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("更新证书记录失败: %v\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新证书记录失败"})
		return
	}

	// 更新用户表中的证书信息
	user.CertID = filePath
	err = userRepo.Update(user)
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

	// 获取数据库连接和仓库
	db, err := database.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库连接失败"})
		return
	}
	repoFactory := repositories.NewRepositoryFactory(db)
	userRepo := repoFactory.GetUserRepository()
	certRepo := repoFactory.GetCertRepository()

	// 检查用户是否存在
	user, err := userRepo.FindByUserID(userIDInt)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("检查用户是否存在失败: %v\n", err)
		}
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "检查用户是否存在失败"})
		}
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
	filePath, err := saveFileToDisk("user", userID, src, true)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("保存密钥文件失败: %v\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存密钥文件失败"})
		return
	}

	// 更新密钥记录
	err = certRepo.UpdateKeyPath("user", userID, filePath)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("更新密钥记录失败: %v\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新密钥记录失败"})
		return
	}

	// 更新用户表中的密钥信息
	user.KeyID = filePath
	err = userRepo.Update(user)
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

	// 获取数据库连接和仓库
	db, err := database.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库连接失败"})
		return
	}
	repoFactory := repositories.NewRepositoryFactory(db)
	deviceRepo := repoFactory.GetDeviceRepository()
	certRepo := repoFactory.GetCertRepository()

	// 检查设备是否存在
	device, err := deviceRepo.FindByDeviceID(deviceID)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("检查设备是否存在失败: %v\n", err)
		}
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "设备不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "检查设备是否存在失败"})
		}
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
	filePath, err := saveFileToDisk("device", deviceID, src, false)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("保存证书文件失败: %v\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存证书文件失败"})
		return
	}

	// 更新证书记录
	err = certRepo.UpdateCertPath("device", deviceID, filePath)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("更新证书记录失败: %v\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新证书记录失败"})
		return
	}

	// 更新设备表中的证书信息
	device.CertID = filePath
	err = deviceRepo.Update(device)
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

	// 获取数据库连接和仓库
	db, err := database.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库连接失败"})
		return
	}
	repoFactory := repositories.NewRepositoryFactory(db)
	deviceRepo := repoFactory.GetDeviceRepository()
	certRepo := repoFactory.GetCertRepository()

	// 检查设备是否存在
	device, err := deviceRepo.FindByDeviceID(deviceID)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("检查设备是否存在失败: %v\n", err)
		}
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "设备不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "检查设备是否存在失败"})
		}
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
	filePath, err := saveFileToDisk("device", deviceID, src, true)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("保存密钥文件失败: %v\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存密钥文件失败"})
		return
	}

	// 更新密钥记录
	err = certRepo.UpdateKeyPath("device", deviceID, filePath)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("更新密钥记录失败: %v\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新密钥记录失败"})
		return
	}

	// 更新设备表中的密钥信息
	device.KeyID = filePath
	err = deviceRepo.Update(device)
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

	// 获取数据库连接和仓库
	db, err := database.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库连接失败"})
		return
	}
	repoFactory := repositories.NewRepositoryFactory(db)
	certRepo := repoFactory.GetCertRepository()

	// 检查实体是否存在
	var exists bool
	if entityType == "user" {
		userID, err := strconv.Atoi(entityID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
			return
		}

		userRepo := repoFactory.GetUserRepository()
		_, err = userRepo.FindByUserID(userID)
		exists = (err == nil)
	} else {
		deviceRepo := repoFactory.GetDeviceRepository()
		_, err = deviceRepo.FindByDeviceID(entityID)
		exists = (err == nil)
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("%s不存在", entityType)})
		return
	}

	// 获取证书信息
	cert, err := certRepo.FindByEntity(entityType, entityID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "未找到证书信息"})
		} else {
			if cfg.DebugLevel == "true" {
				log.Printf("获取证书信息失败: %v\n", err)
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取证书信息失败"})
		}
		return
	}

	// 返回证书信息
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取证书信息成功",
		"data":    cert,
	})
}

// 确保证书目录存在
func ensureCertDirsExist() error {
	cfg := config.GetConfig()

	// 创建基础目录
	basePath := BaseCertPath
	if err := os.MkdirAll(basePath, 0755); err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("创建基础证书目录失败: %v\n", err)
		}
		return fmt.Errorf("创建基础证书目录失败: %w", err)
	}

	// 创建证书目录
	certPath := filepath.Join(basePath, CertsDir)
	if err := os.MkdirAll(certPath, 0755); err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("创建证书目录失败: %v\n", err)
		}
		return fmt.Errorf("创建证书目录失败: %w", err)
	}

	// 创建密钥目录
	keyPath := filepath.Join(basePath, KeysDir)
	if err := os.MkdirAll(keyPath, 0700); err != nil { // 密钥目录权限更严格
		if cfg.DebugLevel == "true" {
			log.Printf("创建密钥目录失败: %v\n", err)
		}
		return fmt.Errorf("创建密钥目录失败: %w", err)
	}

	if cfg.DebugLevel == "true" {
		log.Println("证书目录结构已创建")
	}

	return nil
}

// 获取文件的绝对路径
func getAbsolutePath(relativePath string) string {
	// 获取当前工作目录
	workDir, err := os.Getwd()
	if err != nil {
		log.Printf("获取工作目录失败: %v\n", err)
		return relativePath // 如果失败，返回相对路径
	}

	// 构建绝对路径
	absPath := filepath.Join(workDir, relativePath)

	// 根据操作系统规范化路径
	if runtime.GOOS == "windows" {
		// Windows 下使用反斜杠
		absPath = strings.ReplaceAll(absPath, "/", "\\")
	}

	return absPath
}

// saveFileToDisk 保存文件到磁盘
func saveFileToDisk(entityType, entityID string, fileReader io.Reader, isKey bool) (string, error) {
	cfg := config.GetConfig()

	// 确保目录存在
	if err := ensureCertDirsExist(); err != nil {
		return "", err
	}

	// 确定文件类型和目录
	dirType := CertsDir
	if isKey {
		dirType = KeysDir
	}

	// 构建文件名和路径
	fileName := fmt.Sprintf("%s_%s.pem", entityType, entityID)
	filePath := filepath.Join(BaseCertPath, dirType, fileName)

	// 创建文件
	out, err := os.Create(filePath)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("创建文件失败: %v\n", err)
		}
		return "", fmt.Errorf("创建文件失败: %w", err)
	}
	defer out.Close()

	// 写入文件内容
	_, err = io.Copy(out, fileReader)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("写入文件失败: %v\n", err)
		}
		return "", fmt.Errorf("写入文件失败: %w", err)
	}

	// 设置合适的权限
	if isKey {
		// 密钥文件权限更严格
		if err := os.Chmod(filePath, 0600); err != nil {
			if cfg.DebugLevel == "true" {
				log.Printf("设置密钥文件权限失败: %v\n", err)
			}
			// 不中断流程，只记录错误
			log.Printf("警告: 设置密钥文件权限失败: %v\n", err)
		}
	}

	// 获取并返回绝对路径
	absPath := getAbsolutePath(filePath)
	if cfg.DebugLevel == "true" {
		log.Printf("文件保存成功: %s\n", absPath)
	}

	return absPath, nil
}
