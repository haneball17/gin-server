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
	DeviceName                string  `json:"deviceName" binding:"required,min=4,max=50"` // 设备名称，长度限制，注册时需要
	DeviceType                int     `json:"deviceType" binding:"required"`              // 设备类型，1代表网关设备A型，2代表网关设备B型，3代表网关设备C型，4代表安全接入管理设备，注册时需要
	PassWD                    string  `json:"passWD" binding:"required,min=8"`            // 设备登录口令，注册时需要
	DeviceID                  int     `json:"deviceID" binding:"required"`                // 设备唯一标识，注册时需要
	SuperiorDeviceID          int     `json:"superiorDeviceID" binding:"required"`        // 上级设备ID，注册时需要，当设备为安全接入管理设备时，上级设备ID为0
	CertID                    string  `json:"certID"`                                     // 证书ID，允许为 NULL
	KeyID                     string  `json:"keyID"`                                      // 密钥ID，允许为 NULL
	DeviceStatus              int     `json:"deviceStatus"`                               // 设备状态，注册时需要
	RegisterIP                string  `json:"registerIP"`                                 // 注册IP，注册时需要
	Email                     string  `json:"email"`                                      // 邮箱，注册时需要
	DeviceHardwareFingerprint *string `json:"deviceHardwareFingerprint"`                  // 设备硬件指纹，允许为 NULL
	AnonymousUser             *string `json:"anonymousUser"`                              // 匿名用户，允许为 NULL
}

// RegisterDevice 处理设备注册请求
func RegisterDevice(c *gin.Context) {
	cfg := config.GetConfig() // 获取全局配置

	if cfg.DebugLevel == "true" {
		log.Println("接收到设备注册请求")
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

	// 检查设备 ID 是否存在
	existingDevice, err := deviceRepo.FindByDeviceID(device.DeviceID)
	if err == nil && existingDevice != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "设备 ID 已存在"})
		return
	}

	// 检查设备名称是否存在
	existingDeviceByName, err := deviceRepo.FindByDeviceName(device.DeviceName)
	if err == nil && existingDeviceByName != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "设备名称已存在"})
		return
	}

	// 创建新设备模型
	newDevice := &models.Device{
		DeviceName:       device.DeviceName,
		DeviceType:       device.DeviceType,
		Password:         device.PassWD,
		DeviceID:         device.DeviceID,
		SuperiorDeviceID: device.SuperiorDeviceID,
		DeviceStatus:     device.DeviceStatus,
		CertID:           device.CertID,
		KeyID:            device.KeyID,
		RegisterIP:       device.RegisterIP,
		Email:            device.Email,
	}

	// 处理可能为nil的指针字段
	if device.DeviceHardwareFingerprint != nil {
		newDevice.HardwareFingerprint = *device.DeviceHardwareFingerprint
	}

	if device.AnonymousUser != nil {
		newDevice.AnonymousUser = *device.AnonymousUser
	}

	// 创建设备
	if err := deviceRepo.Create(newDevice); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法创建设备"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "设备注册成功"})
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
	if device.DeviceHardwareFingerprint != nil {
		existingDevice.HardwareFingerprint = *device.DeviceHardwareFingerprint
	}

	if device.AnonymousUser != nil {
		existingDevice.AnonymousUser = *device.AnonymousUser
	}

	// 保存更新
	if err := deviceRepo.Update(existingDevice); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法更新设备信息"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "设备信息更新成功"})
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

	c.JSON(http.StatusOK, gin.H{"message": "设备删除成功"})
}
