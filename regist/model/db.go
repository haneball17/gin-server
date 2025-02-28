package model

import (
	"database/sql" // 导入数据库/sql 包
	"fmt"          // 导入格式化输出包
	"log"          // 导入日志包

	_ "github.com/go-sql-driver/mysql" // 导入 MySQL 驱动
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
	CreatedAt    string `json:"created_at"`                               // 创建时间
}

// Device 结构体定义设备信息
type Device struct {
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
	CreatedAt        string `json:"created_at"`                                 // 创建时间
}

var db *sql.DB // 声明数据库连接变量

// InitDB 初始化数据库连接
func InitDB() {
	var err error
	// 连接数据库，DSN 格式为 "用户名:密码@tcp(主机:端口)/数据库名"
	dsn := "gin_user:your_password@tcp(127.0.0.1:3306)/gin_server"
	db, err = sql.Open("mysql", dsn) // 打开数据库连接
	if err != nil {
		log.Fatal(err) // 如果连接失败，记录错误并退出
	}

	// 测试连接
	if err = db.Ping(); err != nil {
		log.Fatal(err) // 如果连接失败，记录错误并退出
	}

	fmt.Println("数据库连接成功！") // 输出连接成功信息

	// 创建用户表
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		userName VARCHAR(20) NOT NULL,
		passWD VARCHAR(255) NOT NULL,
		email VARCHAR(32),
		userID INT NOT NULL,
		certAddress VARCHAR(32) NOT NULL,
		certDomain VARCHAR(32) NOT NULL,
		certAuthType INT NOT NULL,
		certKeyLen INT NOT NULL,
		secuLevel INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createUsersTable)
	if err != nil {
		log.Fatal("创建用户表失败:", err)
	}

	// 创建设备表
	createDevicesTable := `
	CREATE TABLE IF NOT EXISTS devices (
		id INT AUTO_INCREMENT PRIMARY KEY,
		deviceName VARCHAR(50) NOT NULL,
		deviceType INT NOT NULL,
		passWD VARCHAR(255) NOT NULL,
		deviceID CHAR(12) NOT NULL,
		registerIP VARCHAR(24) NOT NULL,
		superiorDeviceID CHAR(12) NOT NULL,
		email VARCHAR(32),
		certAddress VARCHAR(32) NOT NULL,
		certDomain VARCHAR(32) NOT NULL,
		certAuthType INT NOT NULL,
		certKeyLen INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createDevicesTable)
	if err != nil {
		log.Fatal("创建设备表失败:", err)
	}

	fmt.Println("数据库表创建成功或已存在！") // 输出表创建成功信息
}

// GetDB 返回数据库连接
func GetDB() *sql.DB {
	return db // 返回数据库连接
}

// GetAllUsers 查询所有用户
func GetAllUsers() ([]User, error) {
	var users []User
	rows, err := db.Query("SELECT userName, email, userID, certAddress, certDomain, certAuthType, certKeyLen, secuLevel, created_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.UserName, &user.Email, &user.UserID, &user.CertAddress, &user.CertDomain, &user.CertAuthType, &user.CertKeyLen, &user.SecuLevel, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// GetAllDevices 查询所有设备
func GetAllDevices() ([]Device, error) {
	var devices []Device
	rows, err := db.Query("SELECT deviceName, deviceType, deviceID, registerIP, superiorDeviceID, email, certAddress, certDomain, certAuthType, certKeyLen, created_at FROM devices")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var device Device
		if err := rows.Scan(&device.DeviceName, &device.DeviceType, &device.DeviceID, &device.RegisterIP, &device.SuperiorDeviceID, &device.Email, &device.CertAddress, &device.CertDomain, &device.CertAuthType, &device.CertKeyLen, &device.CreatedAt); err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}
	return devices, nil
}
