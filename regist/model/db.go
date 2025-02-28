package model

import (
	"database/sql" // 导入数据库/sql 包
	"fmt"          // 导入格式化输出包
	"log"          // 导入日志包
	"strings"      // 导入 strings 包

	_ "github.com/go-sql-driver/mysql" // 导入 MySQL 驱动
)

// User 结构体定义用户信息
type User struct {
	UserName           string `json:"userName" binding:"min=4,max=20"` // 用户名，长度限制
	PassWD             string `json:"passWD" binding:"min=8"`          // 密码，长度限制
	Email              string `json:"email" binding:"email"`           // 邮箱，格式校验
	UserID             int    `json:"userID"`                          // 用户唯一标识
	CertAddress        string `json:"certAddress"`                     // 证书地址
	CertDomain         string `json:"certDomain"`                      // 证书域名
	CertAuthType       int    `json:"certAuthType"`                    // 证书认证类型
	CertKeyLen         int    `json:"certKeyLen"`                      // 证书密钥长度
	SecuLevel          int    `json:"secuLevel"`                       // 安全级别
	Status             int    `json:"status"`                          // 账户状态
	PermissionMask     string `json:"permissionMask"`                  // 权限位掩码
	LastLoginTimeStamp string `json:"lastLoginTimeStamp"`              // 登录时间戳
	OffLineTimeStamp   string `json:"offLineTimeStamp"`                // 离线时间戳
	LoginIP            string `json:"loginIP"`                         // 用户登录 IP
	IllegalLoginTimes  int    `json:"illegalLoginTimes"`               // 用户本次的非法登录次数
	CreatedAt          string `json:"created_at"`                      // 创建时间
}

// Device 结构体定义设备信息
type Device struct {
	DeviceName                string  `json:"deviceName"`                // 设备名称，长度限制
	DeviceType                int     `json:"deviceType"`                // 设备类型
	PassWD                    string  `json:"passWD"`                    // 设备登录口令
	DeviceID                  string  `json:"deviceID"`                  // 设备唯一标识
	RegisterIP                string  `json:"registerIP"`                // 上级设备 IP
	SuperiorDeviceID          string  `json:"superiorDeviceID"`          // 上级设备 ID
	Email                     string  `json:"email"`                     // 联系邮箱
	CertAddress               string  `json:"certAddress"`               // 证书地址
	CertDomain                string  `json:"certDomain"`                // 证书域名
	CertAuthType              int     `json:"certAuthType"`              // 证书认证类型
	CertKeyLen                int     `json:"certKeyLen"`                // 证书密钥长度
	DeviceHardwareFingerprint *string `json:"deviceHardwareFingerprint"` // 用户的硬件指纹信息
	CreatedAt                 string  `json:"created_at"`                // 创建时间
}

var db *sql.DB // 声明数据库连接变量

// InitDB 初始化数据库连接
func InitDB() {
	config := LoadConfig() // 加载配置

	var err error
	// 连接数据库，DSN 格式为 "用户名:密码@tcp(主机:端口)/数据库名"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
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
		status INT,
		permissionMask CHAR(8),
		lastLoginTimeStamp DATETIME,
		offLineTimeStamp DATETIME,
		loginIP CHAR(24),
		illegalLoginTimes INT,
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
		deviceHardwareFingerprint CHAR(128),
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

// UpdateUser 更新用户信息
func UpdateUser(userID int, user User) (map[string]interface{}, error) {
	// 获取当前用户信息
	var existingUser User
	err := db.QueryRow("SELECT userName, passWD, email, Status, PermissionMask, LastLoginTimeStamp, OffLineTimeStamp, LoginIP, IllegalLoginTimes, created_at FROM users WHERE userID = ?", userID).
		Scan(&existingUser.UserName, &existingUser.PassWD, &existingUser.Email, &existingUser.Status, &existingUser.PermissionMask, &existingUser.LastLoginTimeStamp, &existingUser.OffLineTimeStamp, &existingUser.LoginIP, &existingUser.IllegalLoginTimes, &existingUser.CreatedAt)

	if err != nil {
		return nil, err // 如果查询失败，返回错误
	}

	// 构建更新 SQL 语句
	updateFields := []string{}
	updateValues := []interface{}{}
	updatedFields := make(map[string]interface{}) // 用于存储更新的字段

	if user.UserName != "" {
		updateFields = append(updateFields, "userName=?")
		updateValues = append(updateValues, user.UserName)
		updatedFields["userName"] = user.UserName
	}
	if user.PassWD != "" {
		updateFields = append(updateFields, "passWD=?")
		updateValues = append(updateValues, user.PassWD)
		updatedFields["passWD"] = user.PassWD
	}
	if user.Email != "" {
		updateFields = append(updateFields, "email=?")
		updateValues = append(updateValues, user.Email)
		updatedFields["email"] = user.Email
	}
	if user.Status != 0 {
		updateFields = append(updateFields, "Status=?")
		updateValues = append(updateValues, user.Status)
		updatedFields["status"] = user.Status
	}
	if user.PermissionMask != "" {
		updateFields = append(updateFields, "PermissionMask=?")
		updateValues = append(updateValues, user.PermissionMask)
		updatedFields["permissionMask"] = user.PermissionMask
	}
	if user.LastLoginTimeStamp != "" {
		updateFields = append(updateFields, "LastLoginTimeStamp=?")
		updateValues = append(updateValues, user.LastLoginTimeStamp)
		updatedFields["lastLoginTimeStamp"] = user.LastLoginTimeStamp
	}
	if user.OffLineTimeStamp != "" {
		updateFields = append(updateFields, "OffLineTimeStamp=?")
		updateValues = append(updateValues, user.OffLineTimeStamp)
		updatedFields["offLineTimeStamp"] = user.OffLineTimeStamp
	}
	if user.LoginIP != "" {
		updateFields = append(updateFields, "LoginIP=?")
		updateValues = append(updateValues, user.LoginIP)
		updatedFields["loginIP"] = user.LoginIP
	}
	if user.IllegalLoginTimes != 0 {
		updateFields = append(updateFields, "IllegalLoginTimes=?")
		updateValues = append(updateValues, user.IllegalLoginTimes)
		updatedFields["illegalLoginTimes"] = user.IllegalLoginTimes
	}

	// 如果没有字段需要更新，直接返回
	if len(updateFields) == 0 {
		return nil, nil
	}

	// 添加 userID 到更新值的最后
	updateValues = append(updateValues, userID)

	// 构建完整的 SQL 语句
	updateSQL := "UPDATE users SET " + strings.Join(updateFields, ", ") + " WHERE userID=?"
	_, err = db.Exec(updateSQL, updateValues...)
	if err != nil {
		return nil, err
	}

	// 返回更新的字段
	return updatedFields, nil
}

// UpdateDevice 更新设备信息
func UpdateDevice(deviceID string, device Device) (map[string]interface{}, error) {
	// 获取当前设备信息
	var existingDevice Device
	var deviceHardwareFingerprint sql.NullString // 使用 sql.NullString 来处理可能为 NULL 的字段
	err := db.QueryRow("SELECT deviceName, deviceType, passWD, registerIP, superiorDeviceID, email, certAddress, certDomain, certAuthType, certKeyLen, DeviceHardwareFingerprint, created_at FROM devices WHERE deviceID = ?", deviceID).
		Scan(&existingDevice.DeviceName, &existingDevice.DeviceType, &existingDevice.PassWD, &existingDevice.RegisterIP, &existingDevice.SuperiorDeviceID, &existingDevice.Email, &existingDevice.CertAddress, &existingDevice.CertDomain, &existingDevice.CertAuthType, &existingDevice.CertKeyLen, &deviceHardwareFingerprint, &existingDevice.CreatedAt)

	if err != nil {
		log.Println("获取设备信息失败:", err) // 打印错误信息
		return nil, err               // 如果查询失败，返回错误
	}

	// 构建更新 SQL 语句
	updateFields := []string{}
	updateValues := []interface{}{}
	updatedFields := make(map[string]interface{}) // 用于存储更新的字段

	if device.DeviceName != "" {
		updateFields = append(updateFields, "deviceName=?")
		updateValues = append(updateValues, device.DeviceName)
		updatedFields["deviceName"] = device.DeviceName
	}
	if device.DeviceType != 0 {
		updateFields = append(updateFields, "deviceType=?")
		updateValues = append(updateValues, device.DeviceType)
		updatedFields["deviceType"] = device.DeviceType
	}
	if device.PassWD != "" {
		updateFields = append(updateFields, "passWD=?")
		updateValues = append(updateValues, device.PassWD)
		updatedFields["passWD"] = device.PassWD
	}
	if device.RegisterIP != "" {
		updateFields = append(updateFields, "registerIP=?")
		updateValues = append(updateValues, device.RegisterIP)
		updatedFields["registerIP"] = device.RegisterIP
	}
	if device.SuperiorDeviceID != "" {
		updateFields = append(updateFields, "superiorDeviceID=?")
		updateValues = append(updateValues, device.SuperiorDeviceID)
		updatedFields["superiorDeviceID"] = device.SuperiorDeviceID
	}
	if device.Email != "" {
		updateFields = append(updateFields, "email=?")
		updateValues = append(updateValues, device.Email)
		updatedFields["email"] = device.Email
	}
	if device.CertAddress != "" {
		updateFields = append(updateFields, "certAddress=?")
		updateValues = append(updateValues, device.CertAddress)
		updatedFields["certAddress"] = device.CertAddress
	}
	if device.CertDomain != "" {
		updateFields = append(updateFields, "certDomain=?")
		updateValues = append(updateValues, device.CertDomain)
		updatedFields["certDomain"] = device.CertDomain
	}
	if device.CertAuthType != 0 {
		updateFields = append(updateFields, "certAuthType=?")
		updateValues = append(updateValues, device.CertAuthType)
		updatedFields["certAuthType"] = device.CertAuthType
	}
	if device.CertKeyLen != 0 {
		updateFields = append(updateFields, "certKeyLen=?")
		updateValues = append(updateValues, device.CertKeyLen)
		updatedFields["certKeyLen"] = device.CertKeyLen
	}
	if device.DeviceHardwareFingerprint != nil {
		updateFields = append(updateFields, "DeviceHardwareFingerprint=?")
		updateValues = append(updateValues, device.DeviceHardwareFingerprint)
		updatedFields["deviceHardwareFingerprint"] = device.DeviceHardwareFingerprint
	}

	// 如果没有字段需要更新，直接返回
	if len(updateFields) == 0 {
		return nil, nil
	}

	// 添加 deviceID 到更新值的最后
	updateValues = append(updateValues, deviceID)

	// 构建完整的 SQL 语句
	updateSQL := "UPDATE devices SET " + strings.Join(updateFields, ", ") + " WHERE deviceID=?"
	_, err = db.Exec(updateSQL, updateValues...)
	if err != nil {
		log.Println("更新设备信息失败:", err) // 打印错误信息
		return nil, err
	}

	// 返回更新的字段
	return updatedFields, nil
}

// CheckUserExists 检查用户是否存在
func CheckUserExists(userName string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE userName = ?", userName).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CheckDeviceExists 检查设备是否存在
func CheckDeviceExists(deviceID string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM devices WHERE deviceID = ?", deviceID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CheckUserExistsByID 检查用户 ID 是否存在
func CheckUserExistsByID(userID int) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE userID = ?", userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CheckUserExistsByName 检查用户名是否存在
func CheckUserExistsByName(userName string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE userName = ?", userName).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CheckDeviceExistsByID 检查设备 ID 是否存在
func CheckDeviceExistsByID(deviceID string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM devices WHERE deviceID = ?", deviceID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CheckDeviceExistsByName 检查设备名称是否存在
func CheckDeviceExistsByName(deviceName string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM devices WHERE deviceName = ?", deviceName).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
