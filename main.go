package main

import (
	"gin-server/regist/model"  // 导入数据库模型包
	"gin-server/regist/router" // 导入路由设置包

	"github.com/gin-gonic/gin" // 导入 Gin 框架
)

func main() {
<<<<<<< HEAD
	model.InitDB()                 // 初始化数据库连接
	r := gin.Default()             // 创建一个默认的 Gin 路由引擎
	router.SetupRouter(r)          // 设置路由
	config := model.LoadConfig()   // 加载配置
	r.Run(":" + config.ServerPort) // 启动服务，监听配置的端口
=======
	model.InitDB()        // 初始化数据库连接
	r := gin.Default()    // 创建一个默认的 Gin 路由引擎
	router.SetupRouter(r) // 设置路由
	r.Run(":8080")        // 启动服务，监听 8080 端口
>>>>>>> acf2b2b3ad5d317a7af3f00ba17d40574692a5ae
}
