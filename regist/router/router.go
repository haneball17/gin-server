package router

import (
	"gin-server/config"
	"gin-server/regist/handler" // 导入处理器包

	"log"

	"github.com/gin-gonic/gin" // 导入 Gin 框架
)

// SetupRouter 设置路由
func SetupRouter(r *gin.Engine) {
	cfg := config.GetConfig() // 获取全局配置

	r.Use(func(c *gin.Context) {
		if cfg.DebugLevel == "true" {
			log.Printf("注册管理 - 请求路径: %s, 方法: %s\n", c.Request.URL.Path, c.Request.Method) // 记录请求路径和方法
		}
		c.Next() // 继续处理请求
		if cfg.DebugLevel == "true" {
			log.Printf("注册管理 - 响应状态: %d\n", c.Writer.Status()) // 记录响应状态
		}
	})

	// 用户注册
	r.POST("/regist/users", handler.RegisterUser)  // 注册用户接口
	r.GET("/search/users", handler.GetUsers)       // 获取所有用户接口
	r.PUT("/update/users/:id", handler.UpdateUser) // 更新用户接口
	// 设备注册
	r.POST("/regist/devices", handler.RegisterDevice)  // 注册设备接口
	r.GET("/search/devices", handler.GetDevices)       // 获取所有设备接口
	r.PUT("/update/devices/:id", handler.UpdateDevice) // 更新设备接口

	// 新增：指定用户和设备查询接口
	r.GET("/search/user", handler.GetUserByID)     // 根据ID查询用户接口
	r.GET("/search/device", handler.GetDeviceByID) // 根据ID查询设备接口

	// 证书管理路由
	// 用户证书绑定
	r.POST("/bind/users/:id/cert", handler.BindUserCert) // 用户证书绑定接口
	r.POST("/bind/users/:id/key", handler.BindUserKey)   // 用户密钥绑定接口

	// 设备证书绑定
	r.POST("/bind/devices/:id/cert", handler.BindDeviceCert) // 设备证书绑定接口
	r.POST("/bind/devices/:id/key", handler.BindDeviceKey)   // 设备密钥绑定接口

	// 获取证书信息
	r.GET("/cert/info", handler.GetCertInfo) // 获取证书信息接口
}
