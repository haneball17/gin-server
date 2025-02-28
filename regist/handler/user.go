package handler

import (
	"net/http"

	"gin-server/regist/model"

	"github.com/gin-gonic/gin"
)

// User 结构体定义用户信息
type User struct {
	UserName     string `json:"userName" binding:"required,min=4,max=20"` // 用户名，必填，长度限制
	PassWD       string `json:"passWD" binding:"required,min=8"`          // 密码，必填，长度限制
	Email        string `json:"email" binding:"email"`                    // 邮箱，格式校验
	UserID       int    `json:"userID" binding:"required"`                // 用户唯一标识，必填
	CertAddress  string `json:"certAddress" binding:"required"`           // 证书地址，必填
	CertDomain   string `json:"certDomain" binding:"required"`            // 证书域名，必填
	CertAuthType int    `json:"certAuthType" binding:"required"`          // 证书认证类型，必填
	CertKeyLen   int    `json:"certKeyLen" binding:"required"`            // 证书密钥长度，必填
	SecuLevel    int    `json:"secuLevel" binding:"required"`             // 安全级别，必填
}

// RegisterUser 处理用户注册请求
func RegisterUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 返回参数错误信息
		return
	}

	// 插入用户信息到数据库
	db := model.GetDB() // 获取数据库连接
	_, err := db.Exec("INSERT INTO users (userName, passWD, email, userID, certAddress, certDomain, certAuthType, certKeyLen, secuLevel) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		user.UserName, user.PassWD, user.Email, user.UserID, user.CertAddress, user.CertDomain, user.CertAuthType, user.CertKeyLen, user.SecuLevel)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法创建用户"}) // 返回创建用户失败信息
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    201,
		"message": "User created", // 返回用户创建成功信息
		"data": gin.H{
			"userName":   user.UserName,
			"email":      user.Email,
			"created_at": "ISO8601时间戳", // 这里可以用实际时间戳替代
		},
	})
}

// GetUsers 处理获取所有用户的请求
func GetUsers(c *gin.Context) {
	users, err := model.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取用户列表"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}
