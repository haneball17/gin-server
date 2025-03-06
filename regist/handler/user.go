package handler

import (
	"log"
	"net/http"
	"strconv"
	"time" // 导入时间包

	"gin-server/config"
	"gin-server/regist/model"

	"github.com/gin-gonic/gin"
)

// User 结构体定义用户信息
type User struct {
	UserName        string `json:"userName" binding:"required,min=4,max=20"` // 用户名，必填，长度限制，注册时需要
	PassWD          string `json:"passWD" binding:"required,min=8"`          // 密码，必填，长度限制，注册时需要
	UserID          int    `json:"userID" binding:"required"`                // 用户唯一标识，必填，注册时需要
	UserType        int    `json:"userType" binding:"required"`              // 用户类型，注册时需要
	GatewayDeviceID string `json:"gatewayDeviceID" binding:"required"`       // 用户所属网关设备ID，注册时需要，作为外键关联到设备表
	CertID          string `json:"certID"`                                   // 证书ID，允许为 NULL
	KeyID           string `json:"keyID"`                                    // 密钥ID，允许为 NULL
}

// RegisterUser 处理用户注册请求
func RegisterUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 返回参数错误信息
		return
	}

	cfg := config.GetConfig() // 获取全局配置

	// 检查用户 ID 是否存在
	existsID, err := model.CheckUserExistsByID(user.UserID)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("无法检查用户 ID 是否存在: %v\n", err) // 记录错误信息
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法检查用户 ID 是否存在"}) // 返回检查失败信息
		return
	}
	if existsID {
		c.JSON(http.StatusConflict, gin.H{"error": "用户 ID 已存在"}) // 返回冲突错误信息
		return
	}

	// 检查用户名是否存在
	existsName, err := model.CheckUserExistsByName(user.UserName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法检查用户名是否存在"}) // 返回检查失败信息
		return
	}
	if existsName {
		c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"}) // 返回冲突错误信息
		return
	}

	// 插入用户信息到数据库
	db := model.GetDB() // 获取数据库连接
	_, err = db.Exec("INSERT INTO users (userName, passWD, userID, userType, gatewayDeviceID, certID, keyID) VALUES (?, ?, ?, ?, ?, ?, ?)",
		user.UserName, user.PassWD, user.UserID, user.UserType, user.GatewayDeviceID, user.CertID, user.KeyID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法创建用户"}) // 返回创建用户失败信息
		return
	}

	// 获取当前时间并格式化为 ISO 8601
	createdAt := time.Now().Format(time.RFC3339)

	c.JSON(http.StatusCreated, gin.H{
		"code":    201,
		"message": "User created", // 返回用户创建成功信息
		"data": gin.H{
			"userName":   user.UserName,
			"userID":     user.UserID,
			"created_at": createdAt, // 返回实际创建时间
		},
	})
}

// GetUsers 处理获取所有用户的请求
func GetUsers(c *gin.Context) {
	cfg := config.GetConfig() // 获取全局配置

	if cfg.DebugLevel == "true" {
		log.Println("接收到获取所有用户的请求")
	}

	users, err := model.GetAllUsers()
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("获取用户列表失败: %v\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取用户列表"})
		return
	}

	if cfg.DebugLevel == "true" {
		log.Printf("成功获取 %d 个用户信息\n", len(users))
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// UpdateUser 处理用户修改请求
func UpdateUser(c *gin.Context) {
	cfg := config.GetConfig() // 获取全局配置

	if cfg.DebugLevel == "true" {
		log.Println("接收到更新用户的请求")
	}

	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 返回参数错误信息
		return
	}

	userIDStr := c.Param("id")             // 获取路径参数中的用户 ID
	userID, err := strconv.Atoi(userIDStr) // 将字符串转换为整数
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户 ID"}) // 返回无效用户 ID 错误信息
		return
	}

	// 更新用户信息
	updatedFields, err := model.UpdateUser(userID, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法更新用户信息"}) // 返回更新失败信息
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "User updated successfully", // 返回用户更新成功信息
		"data":    updatedFields,               // 返回更新的字段
	})
}
