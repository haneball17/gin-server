package handler

import (
	"net/http"
<<<<<<< HEAD
	"time" // 导入时间包
=======
>>>>>>> acf2b2b3ad5d317a7af3f00ba17d40574692a5ae

	"gin-server/regist/model"

	"github.com/gin-gonic/gin"
)

// Device 结构体定义设备信息
type Device struct {
<<<<<<< HEAD
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
=======
	DeviceName       string `json:"deviceName" binding:"required,min=2,max=50"` // 设备名称，必填，长度限制
	DeviceType       int    `json:"deviceType" binding:"required"`              // 设备类型，必填
	PassWD           string `json:"passWD" binding:"required,min=8"`            // 设备登录口令，必填，长度限制
	DeviceID         string `json:"deviceID" binding:"required,len=12"`         // 设备唯一标识，必填，固定长度
	RegisterIP       string `json:"registerIP" binding:"required"`              // 上级设备 IP，必填
	SuperiorDeviceID string `json:"superiorDeviceID" binding:"required"`        // 上级设备 ID，必填
	Email            string `json:"email" binding:"email"`                      // 联系邮箱，格式校验
	CertAddress      string `json:"certAddress" binding:"required"`             // 证书地址，必填
	CertDomain       string `json:"certDomain" binding:"required"`              // 证书域名，必填
	CertAuthType     int    `json:"certAuthType" binding:"required"`            // 证书认证类型，必填
	CertKeyLen       int    `json:"certKeyLen" binding:"required"`              // 证书密钥长度，必填
>>>>>>> acf2b2b3ad5d317a7af3f00ba17d40574692a5ae
}

// RegisterDevice 处理设备注册请求
func RegisterDevice(c *gin.Context) {
	var device Device
	if err := c.ShouldBindJSON(&device); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 返回参数错误信息
		return
	}

<<<<<<< HEAD
	// 检查设备 ID 是否存在
	existsID, err := model.CheckDeviceExistsByID(device.DeviceID)
	if err != nil {
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
=======
	// 插入设备信息到数据库
	db := model.GetDB() // 获取数据库连接
	_, err := db.Exec("INSERT INTO devices (deviceName, deviceType, passWD, deviceID, registerIP, superiorDeviceID, email, certAddress, certDomain, certAuthType, certKeyLen) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
>>>>>>> acf2b2b3ad5d317a7af3f00ba17d40574692a5ae
		device.DeviceName, device.DeviceType, device.PassWD, device.DeviceID, device.RegisterIP, device.SuperiorDeviceID, device.Email, device.CertAddress, device.CertDomain, device.CertAuthType, device.CertKeyLen)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法创建设备"}) // 返回创建设备失败信息
		return
	}

<<<<<<< HEAD
	// 获取当前时间并格式化为 ISO 8601
	registeredAt := time.Now().Format(time.RFC3339)

=======
>>>>>>> acf2b2b3ad5d317a7af3f00ba17d40574692a5ae
	c.JSON(http.StatusCreated, gin.H{
		"code":    201,
		"message": "Device registered", // 返回设备注册成功信息
		"data": gin.H{
			"deviceName":    device.DeviceName,
			"email":         device.Email,
<<<<<<< HEAD
			"registered_at": registeredAt, // 返回实际注册时间
=======
			"registered_at": "ISO8601时间戳", // 这里可以用实际时间戳替代
>>>>>>> acf2b2b3ad5d317a7af3f00ba17d40574692a5ae
		},
	})
}

// GetDevices 处理获取所有设备的请求
func GetDevices(c *gin.Context) {
	devices, err := model.GetAllDevices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取设备列表"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"devices": devices})
}
<<<<<<< HEAD

// UpdateDevice 处理设备修改请求
func UpdateDevice(c *gin.Context) {
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
=======
>>>>>>> acf2b2b3ad5d317a7af3f00ba17d40574692a5ae
