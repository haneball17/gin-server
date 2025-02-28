package router

import (
	"gin-server/regist/handler" // 导入处理器包

	"github.com/gin-gonic/gin" // 导入 Gin 框架
)

// SetupRouter 设置路由
func SetupRouter(r *gin.Engine) {
	// 用户注册
	r.POST("/regist/users", handler.RegisterUser) // 注册用户接口
	r.GET("/search/users", handler.GetUsers)      // 获取所有用户接口
	// 设备注册
	r.POST("/regist/devices", handler.RegisterDevice) // 注册设备接口
	r.GET("/search/devices", handler.GetDevices)      // 获取所有设备接口
}
