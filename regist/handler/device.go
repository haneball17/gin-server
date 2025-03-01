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
	DeviceName                string `json:"deviceName" binding:"required,min=2,max=50"` // 设备名称，必填，长度限制
	DeviceType                int    `json:"deviceType"`                                 // 设备类型
	PassWD                    string `json:"passWD" binding:"min=8"`                     // 设备登录口令，长度限制
	DeviceID                  string `json:"deviceID" binding:"required,len=12"`         // 设备唯一标识，必填，固定长度
	RegisterIP                string `json:"registerIP"`                                 // 上级设备 IP
	SuperiorDeviceID          string `json:"superiorDeviceID"`                           // 上级设备 ID
	Email                     string `json:"email"`                                      // 联系邮箱
	CertAddress               string `json:"certAddress"`                                // 证书地址
	CertDomain                string `json:"certDomain"`                                 // 证书域名
	CertAuthType              int    `json:"certAuthType"`                               // 证书认证类型
	CertKeyLen                int    `json:"certKeyLen"`                                 // 证书密钥长度
	DeviceHardwareFingerprint string `json:"deviceHardwareFingerprint"`                  // 用户的硬件指纹信息
	CreatedAt                 string `json:"created_at"`                                 // 创建时间
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
	_, err = db.Exec("INSERT INTO devices (deviceName, deviceType, passWD, deviceID, registerIP, superiorDeviceID, email, certAddress, certDomain, certAuthType, certKeyLen) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		device.DeviceName, device.DeviceType, device.PassWD, device.DeviceID, device.RegisterIP, device.SuperiorDeviceID, device.Email, device.CertAddress, device.CertDomain, device.CertAuthType, device.CertKeyLen)

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
			"email":         device.Email,
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
