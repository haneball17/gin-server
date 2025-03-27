package handler

import (
	"log"
	"net/http"
	"strconv"

	"gin-server/config"
	"gin-server/database"
	"gin-server/database/models"
	"gin-server/database/repositories"

	// 临时保留，后续完全迁移后可删除
	"github.com/gin-gonic/gin"
)

// Device 结构体定义设备信息
type Device struct {
	DeviceName          string  `json:"device_name" binding:"required,min=4,max=50"` // 设备名称，长度限制，注册时需要
	DeviceType          int     `json:"device_type" binding:"required"`              // 设备类型，1代表网关设备A型，2代表网关设备B型，3代表网关设备C型，4代表安全接入管理设备，注册时需要
	PassWD              string  `json:"pass_wd" binding:"required,min=8"`            // 设备登录口令，注册时需要
	DeviceID            int     `json:"device_id" binding:"required"`                // 设备唯一标识，注册时需要
	SuperiorDeviceID    int     `json:"superior_device_id" binding:"required"`       // 上级设备ID，注册时需要，当设备为安全接入管理设备时，上级设备ID为0
	CertID              string  `json:"cert_id"`                                     // 证书ID，允许为 NULL
	KeyID               string  `json:"key_id"`                                      // 密钥ID，允许为 NULL
	DeviceStatus        int     `json:"device_status"`                               // 设备状态，注册时需要
	RegisterIP          string  `json:"register_ip"`                                 // 注册IP，注册时需要
	Email               string  `json:"email"`                                       // 邮箱，注册时需要
	HardwareFingerprint *string `json:"hardware_fingerprint"`                        // 设备硬件指纹，允许为 NULL
	AnonymousUser       *string `json:"anonymous_user"`                              // 匿名用户，允许为 NULL
}

// DeviceRegisterRequest 简化的设备注册请求结构体
type DeviceRegisterRequest struct {
	DeviceName       string `json:"device_name" binding:"required,min=4,max=50"` // 设备名称，长度限制，注册时需要
	DeviceType       int    `json:"device_type" binding:"required"`              // 设备类型，1代表网关设备A型，2代表网关设备B型，3代表网关设备C型，4代表安全接入管理设备，注册时需要
	PassWD           string `json:"pass_wd" binding:"required,min=8"`            // 设备登录口令，注册时需要
	DeviceID         int    `json:"device_id" binding:"required"`                // 设备唯一标识，注册时需要
	SuperiorDeviceID int    `json:"superior_device_id" `                         // 上级设备ID，注册时需要，当设备为安全接入管理设备时，上级设备ID为0
}

// RegisterDevice 处理设备注册请求
func RegisterDevice(c *gin.Context) {
	cfg := config.GetConfig() // 获取全局配置

	if cfg.DebugLevel == "true" {
		log.Println("接收到设备注册请求")
	}

	var request DeviceRegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 返回参数错误信息
		return
	}

	// 获取数据库连接和仓库
	db, err := database.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库连接失败"})
		return
	}
	repoFactory := repositories.NewRepositoryFactory(db)
	deviceRepo := repoFactory.GetDeviceRepository()

	// 检查设备 ID 是否存在
	existingDevice, err := deviceRepo.FindByDeviceID(request.DeviceID)
	if err == nil && existingDevice != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "设备 ID 已存在"})
		return
	}

	// 检查设备名称是否存在
	existingDeviceByName, err := deviceRepo.FindByDeviceName(request.DeviceName)
	if err == nil && existingDeviceByName != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "设备名称已存在"})
		return
	}

	// 创建新设备模型
	newDevice := &models.Device{
		DeviceName:       request.DeviceName,
		DeviceType:       request.DeviceType,
		Password:         request.PassWD,
		DeviceID:         request.DeviceID,
		SuperiorDeviceID: request.SuperiorDeviceID,
		DeviceStatus:     2, // 默认离线状态
		CertID:           "",
		KeyID:            "",
		RegisterIP:       c.ClientIP(), // 自动获取客户端IP
		Email:            "",
	}

	// 创建设备
	if err := deviceRepo.Create(newDevice); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法创建设备"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "设备注册成功",
		"data":    newDevice,
	})
}

// GetDevices 处理获取所有设备的请求
func GetDevices(c *gin.Context) {
	cfg := config.GetConfig() // 获取全局配置

	if cfg.DebugLevel == "true" {
		log.Println("接收到获取所有设备的请求")
	}

	// 获取数据库连接和仓库
	db, err := database.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库连接失败"})
		return
	}
	repoFactory := repositories.NewRepositoryFactory(db)
	deviceRepo := repoFactory.GetDeviceRepository()

	// 获取所有设备
	devices, err := deviceRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取设备列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取设备列表成功",
		"data":    devices,
	})
}

// UpdateDevice 处理更新设备的请求
func UpdateDevice(c *gin.Context) {
	cfg := config.GetConfig() // 获取全局配置

	if cfg.DebugLevel == "true" {
		log.Println("接收到更新设备的请求")
	}

	// 获取路径参数中的设备 ID
	deviceIDStr := c.Param("id")

	// 转换设备ID为整数
	deviceID, err := strconv.Atoi(deviceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的设备ID格式"})
		return
	}

	var device Device
	if err := c.ShouldBindJSON(&device); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 返回参数错误信息
		return
	}

	// 获取数据库连接和仓库
	db, err := database.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库连接失败"})
		return
	}
	repoFactory := repositories.NewRepositoryFactory(db)
	deviceRepo := repoFactory.GetDeviceRepository()

	// 查找设备
	existingDevice, err := deviceRepo.FindByDeviceID(deviceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "设备不存在"})
		return
	}

	// 更新设备字段
	existingDevice.DeviceName = device.DeviceName
	existingDevice.DeviceType = device.DeviceType
	existingDevice.Password = device.PassWD
	existingDevice.SuperiorDeviceID = device.SuperiorDeviceID
	existingDevice.DeviceStatus = device.DeviceStatus
	existingDevice.CertID = device.CertID
	existingDevice.KeyID = device.KeyID
	existingDevice.RegisterIP = device.RegisterIP
	existingDevice.Email = device.Email

	// 处理可能为nil的指针字段
	if device.HardwareFingerprint != nil {
		existingDevice.HardwareFingerprint = *device.HardwareFingerprint
	}

	if device.AnonymousUser != nil {
		existingDevice.AnonymousUser = *device.AnonymousUser
	}

	// 保存更新
	if err := deviceRepo.Update(existingDevice); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法更新设备信息"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "设备信息更新成功",
		"data":    existingDevice,
	})
}

// DeleteDevice 处理删除设备的请求
func DeleteDevice(c *gin.Context) {
	cfg := config.GetConfig() // 获取全局配置

	if cfg.DebugLevel == "true" {
		log.Println("接收到删除设备的请求")
	}

	// 获取路径参数中的设备 ID
	deviceIDStr := c.Param("id")

	// 转换设备ID为整数
	deviceID, err := strconv.Atoi(deviceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的设备ID格式"})
		return
	}

	// 获取数据库连接和仓库
	db, err := database.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库连接失败"})
		return
	}
	repoFactory := repositories.NewRepositoryFactory(db)
	deviceRepo := repoFactory.GetDeviceRepository()

	// 查找设备
	existingDevice, err := deviceRepo.FindByDeviceID(deviceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "设备不存在"})
		return
	}

	// 删除设备
	if err := deviceRepo.Delete(existingDevice.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法删除设备"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "设备删除成功",
	})
}
