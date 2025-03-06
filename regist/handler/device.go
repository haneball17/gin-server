package handler

import (
	"log"
	"net/http"
	"time" // 导入时间包

	"gin-server/config"
	"gin-server/regist/model"

	"github.com/gin-gonic/gin"
)

// Device 结构体定义设备信息
type Device struct {
	DeviceName       string `json:"deviceName" binding:"required,min=4,max=50"` // 设备名称，长度限制，注册时需要
	DeviceType       int    `json:"deviceType" binding:"required"`              // 设备类型，1代表网关设备A型，2代表网关设备B型，3代表网关设备C型，4代表安全接入管理设备，注册时需要
	PassWD           string `json:"passWD" binding:"required,min=8"`            // 设备登录口令，注册时需要
	DeviceID         string `json:"deviceID" binding:"required"`                // 设备唯一标识，注册时需要
	SuperiorDeviceID string `json:"superiorDeviceID" binding:"required"`        // 上级设备ID，注册时需要，当设备为安全接入管理设备时，上级设备ID为空
	CertID           string `json:"certID"`                                     // 证书ID，允许为 NULL
	KeyID            string `json:"keyID"`                                      // 密钥ID，允许为 NULL
}

// RegisterDevice 处理设备注册请求
func RegisterDevice(c *gin.Context) {
	var device Device
	if err := c.ShouldBindJSON(&device); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 返回参数错误信息
		return
	}

	cfg := config.GetConfig() // 获取全局配置

	// 检查设备 ID 是否存在
	existsID, err := model.CheckDeviceExistsByID(device.DeviceID)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("无法检查设备 ID 是否存在: %v\n", err) // 记录错误信息
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法检查设备 ID 是否存在"}) // 返回检查失败信息
		return
	}
	if existsID {
		c.JSON(http.StatusConflict, gin.H{"error": "设备 ID 已存在"}) // 返回冲突错误信息
		return
	}

	// 检查设备名称是否存在
	existsName, err := model.CheckDeviceExistsByName(device.DeviceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法检查设备名称是否存在"}) // 返回检查失败信息
		return
	}
	if existsName {
		c.JSON(http.StatusConflict, gin.H{"error": "设备名称已存在"}) // 返回冲突错误信息
		return
	}

	// 插入设备信息到数据库
	db := model.GetDB() // 获取数据库连接
	_, err = db.Exec("INSERT INTO devices (deviceName, deviceType, passWD, deviceID, superiorDeviceID, certID, keyID) VALUES (?, ?, ?, ?, ?, ?, ?)",
		device.DeviceName, device.DeviceType, device.PassWD, device.DeviceID, device.SuperiorDeviceID, device.CertID, device.KeyID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法创建设备"}) // 返回创建设备失败信息
		return
	}

	// 获取当前时间并格式化为 ISO 8601
	registeredAt := time.Now().Format(time.RFC3339)

	c.JSON(http.StatusCreated, gin.H{
		"code":    201,
		"message": "Device registered", // 返回设备注册成功信息
		"data": gin.H{
			"deviceName":    device.DeviceName,
			"deviceID":      device.DeviceID,
			"registered_at": registeredAt, // 返回实际注册时间
		},
	})
}

// GetDevices 处理获取所有设备的请求
func GetDevices(c *gin.Context) {
	cfg := config.GetConfig() // 获取全局配置

	if cfg.DebugLevel == "true" {
		log.Println("接收到获取所有设备的请求")
	}

	devices, err := model.GetAllDevices()
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("获取设备列表失败: %v\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取设备列表"})
		return
	}

	if cfg.DebugLevel == "true" {
		log.Printf("成功获取 %d 个设备信息\n", len(devices))
	}

	c.JSON(http.StatusOK, gin.H{"devices": devices})
}

// UpdateDevice 处理设备修改请求
func UpdateDevice(c *gin.Context) {
	cfg := config.GetConfig() // 获取全局配置

	if cfg.DebugLevel == "true" {
		log.Println("接收到更新设备的请求")
	}

	var device model.Device
	if err := c.ShouldBindJSON(&device); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 返回参数错误信息
		return
	}

	deviceID := c.Param("id") // 获取路径参数中的设备 ID

	// 更新设备信息
	updatedFields, err := model.UpdateDevice(deviceID, device)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法更新设备信息"}) // 返回更新失败信息
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Device updated successfully", // 返回设备更新成功信息
		"data":    updatedFields,                 // 返回更新的字段
	})
}
