package router

import (
	"gin-server/regist/handler" // 导入处理器包

	"log"

	"github.com/gin-gonic/gin" // 导入 Gin 框架
)

// SetupRouter 设置路由
func SetupRouter(r *gin.Engine) {
	r.Use(func(c *gin.Context) {
		log.Printf("请求路径: %s, 方法: %s\n", c.Request.URL.Path, c.Request.Method) // 记录请求路径和方法
		c.Next()                                                               // 继续处理请求
		log.Printf("响应状态: %d\n", c.Writer.Status())                            // 记录响应状态
	})

	// 用户注册
	r.POST("/regist/users", handler.RegisterUser)  // 注册用户接口
	r.GET("/search/users", handler.GetUsers)       // 获取所有用户接口
	r.PUT("/update/users/:id", handler.UpdateUser) // 更新用户接口
	// 设备注册
	r.POST("/regist/devices", handler.RegisterDevice)  // 注册设备接口
	r.GET("/search/devices", handler.GetDevices)       // 获取所有设备接口
	r.PUT("/update/devices/:id", handler.UpdateDevice) // 更新设备接口
}
