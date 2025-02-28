package main

import (
	"gin-server/regist/model"  // 导入数据库模型包
	"gin-server/regist/router" // 导入路由设置包

	"github.com/gin-gonic/gin" // 导入 Gin 框架
)

func main() {
	model.InitDB()        // 初始化数据库连接
	r := gin.Default()    // 创建一个默认的 Gin 路由引擎
	router.SetupRouter(r) // 设置路由
	r.Run(":8080")        // 启动服务，监听 8080 端口
}
