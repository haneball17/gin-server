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

// User 结构体定义用户信息
type User struct {
	UserName        string `json:"userName" binding:"required,min=4,max=20"` // 用户名，必填，长度限制，注册时需要
	PassWD          string `json:"passWD" binding:"required,min=8"`          // 密码，必填，长度限制，注册时需要
	UserID          int    `json:"userID" binding:"required"`                // 用户唯一标识，必填，注册时需要
	UserType        int    `json:"userType" binding:"required"`              // 用户类型，注册时需要
	GatewayDeviceID string `json:"gatewayDeviceID" binding:"required"`       // 用户所属网关设备ID，注册时需要，作为外键关联到设备表
	CertID          string `json:"certID"`                                   // 证书ID，允许为 NULL
	KeyID           string `json:"keyID"`                                    // 密钥ID，允许为 NULL
	Email           string `json:"email"`                                    // 邮箱，允许为 NULL
}

// RegisterUser 处理用户注册请求
func RegisterUser(c *gin.Context) {
	cfg := config.GetConfig() // 获取全局配置

	if cfg.DebugLevel == "true" {
		log.Println("接收到用户注册请求")
	}

	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
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
	userRepo := repoFactory.GetUserRepository()

	// 检查用户 ID 是否存在
	if _, err := userRepo.FindByID(uint(user.UserID)); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "用户 ID 已存在"})
		return
	}

	// 检查用户名是否存在
	if _, err := userRepo.FindByUsername(user.UserName); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
		return
	}

	// 创建新用户模型
	newUser := &models.User{
		Username:        user.UserName,
		Password:        user.PassWD,
		UserID:          user.UserID,
		UserType:        user.UserType,
		GatewayDeviceID: user.GatewayDeviceID,
		Status:          nil, // 默认为 NULL
		OnlineDuration:  0,   // 默认为 0
		CertID:          user.CertID,
		KeyID:           user.KeyID,
		Email:           user.Email,
	}

	// 创建用户
	if err := userRepo.Create(newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法创建用户"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户注册成功"})
}

// GetUsers 处理获取所有用户的请求
func GetUsers(c *gin.Context) {
	cfg := config.GetConfig() // 获取全局配置

	if cfg.DebugLevel == "true" {
		log.Println("接收到获取所有用户的请求")
	}

	// 获取数据库连接和仓库
	db, err := database.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库连接失败"})
		return
	}
	repoFactory := repositories.NewRepositoryFactory(db)
	userRepo := repoFactory.GetUserRepository()

	// 获取所有用户
	users, err := userRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取用户列表"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// UpdateUser 处理用户修改请求
func UpdateUser(c *gin.Context) {
	cfg := config.GetConfig() // 获取全局配置

	if cfg.DebugLevel == "true" {
		log.Println("接收到更新用户的请求")
	}

	var requestUser User
	if err := c.ShouldBindJSON(&requestUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 返回参数错误信息
		return
	}

	userIDStr := c.Param("id")             // 获取路径参数中的用户 ID
	userID, err := strconv.Atoi(userIDStr) // 将字符串转换为整数
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户 ID"}) // 返回无效用户 ID 错误信息
		return
	}

	// 获取数据库连接和仓库
	db, err := database.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库连接失败"})
		return
	}
	repoFactory := repositories.NewRepositoryFactory(db)
	userRepo := repoFactory.GetUserRepository()

	// 查找用户
	existingUser, err := userRepo.FindByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 更新用户字段
	existingUser.Username = requestUser.UserName
	if requestUser.PassWD != "" {
		existingUser.Password = requestUser.PassWD
	}
	existingUser.UserID = requestUser.UserID
	existingUser.UserType = requestUser.UserType
	existingUser.GatewayDeviceID = requestUser.GatewayDeviceID
	existingUser.CertID = requestUser.CertID
	existingUser.KeyID = requestUser.KeyID
	existingUser.Email = requestUser.Email
	// 不更新其他字段，保持原值

	// 保存更新
	if err := userRepo.Update(existingUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法更新用户信息"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户信息更新成功"})
}

// DeleteUser 处理删除用户的请求
func DeleteUser(c *gin.Context) {
	cfg := config.GetConfig() // 获取全局配置

	if cfg.DebugLevel == "true" {
		log.Println("接收到删除用户的请求")
	}

	userIDStr := c.Param("id")             // 获取路径参数中的用户 ID
	userID, err := strconv.Atoi(userIDStr) // 将字符串转换为整数
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户 ID"}) // 返回无效用户 ID 错误信息
		return
	}

	// 获取数据库连接和仓库
	db, err := database.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库连接失败"})
		return
	}
	repoFactory := repositories.NewRepositoryFactory(db)
	userRepo := repoFactory.GetUserRepository()

	// 删除用户
	if err := userRepo.Delete(uint(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法删除用户"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户删除成功"})
}
