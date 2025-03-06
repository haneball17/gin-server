package model

import (
	"database/sql" // 导入数据库/sql 包
	"fmt"          // 导入格式化输出包
	"log"          // 导入日志包
	"strings"      // 导入 strings 包
	"time"         // 导入 time 包

	"gin-server/config" // 导入全局配置包

	_ "github.com/go-sql-driver/mysql" // 导入 MySQL 驱动
)

// User 结构体定义用户信息
type User struct {
	UserName        string `json:"userName"`        // 用户名，必填，长度限制，注册时需要
	PassWD          string `json:"passWD"`          // 密码，必填，长度限制，注册时需要
	UserID          int    `json:"userID"`          // 用户唯一标识，必填，注册时需要
	UserType        int    `json:"userType"`        // 用户类型，注册时需要
	GatewayDeviceID string `json:"gatewayDeviceID"` // 用户所属网关设备ID，注册时需要，作为外键关联到设备表

	Status         sql.NullInt64 `json:"status"`         // 账户状态，允许为 NULL，上报时需要
	OnlineDuration int           `json:"onlineDuration"` // 在线时长，上报时需要，允许为 NULL

	CertID             string         `json:"certID"`             // 证书ID，允许为 NULL
	KeyID              string         `json:"keyID"`              // 密钥ID，允许为 NULL
	Email              string         `json:"email"`              // 邮箱，格式校验，允许为 NULL
	PermissionMask     sql.NullString `json:"permissionMask"`     // 权限位掩码，允许为 NULL
	LastLoginTimeStamp sql.NullString `json:"lastLoginTimeStamp"` // 登录时间戳，允许为 NULL
	OffLineTimeStamp   sql.NullString `json:"offLineTimeStamp"`   // 离线时间戳，允许为 NULL
	LoginIP            sql.NullString `json:"loginIP"`            // 用户登录 IP，允许为 NULL
	IllegalLoginTimes  sql.NullInt64  `json:"illegalLoginTimes"`  // 用户本次的非法登录次数，允许为 NULL

	CreatedAt string `json:"created_at"` // 创建时间
}

// UserBehavior 结构体定义用户行为信息
type UserBehavior struct {
	UserID       int       `json:"userID"`       // 用户ID，作为外键关联到用户表
	BehaviorID   int       `json:"behaviorID"`   // 行为ID
	BehaviorTime time.Time `json:"behaviorTime"` // 行为开始时间
	BehaviorType int       `json:"behaviorType"` // 行为类型，1代表发送，2代表输出
	DataType     int       `json:"dataType"`     // 数据类型，1代表文件，2代表消息
	DataSize     int64     `json:"dataSize"`     // 数据大小
}

// Device 结构体定义设备信息
type Device struct {
	DeviceName       string `json:"deviceName"`       // 设备名称，长度限制，注册时需要
	DeviceType       int    `json:"deviceType"`       // 设备类型，1代表网关设备A型，2代表网关设备B型，3代表网关设备C型，4代表安全接入管理设备，注册时需要
	PassWD           string `json:"passWD"`           // 设备登录口令，注册时需要
	DeviceID         string `json:"deviceID"`         // 设备唯一标识，注册时需要
	SuperiorDeviceID string `json:"superiorDeviceID"` // 上级设备ID，注册时需要，当设备为安全接入管理设备时，上级设备ID为空

	DeviceStatus    int `json:"deviceStatus"`    // 设备状态，1代表在线，2代表离线，3代表冻结，4代表注销，上报时需要,允许为 NULL
	PeakCPUUsage    int `json:"peakCPUUsage"`    // 峰值CPU使用率，上报时需要，允许为 NULL
	PeakMemoryUsage int `json:"peakMemoryUsage"` // 峰值内存使用率，上报时需要，允许为 NULL
	OnlineDuration  int `json:"onlineDuration"`  // 在线时长，上报时需要，允许为 NULL

	CertID                    string  `json:"certID"`                    // 证书ID，允许为 NULL
	KeyID                     string  `json:"keyID"`                     // 密钥ID，允许为 NULL
	RegisterIP                string  `json:"registerIP"`                // 上级设备 IP，允许为 NULL
	Email                     string  `json:"email"`                     // 联系邮箱，允许为 NULL
	DeviceHardwareFingerprint *string `json:"deviceHardwareFingerprint"` // 用户的硬件指纹信息，允许为 NULL
	AnonymousUser             *string `json:"anonymousUser"`             // 匿名用户，允许为 NULL

	CreatedAt string `json:"created_at"` // 创建时间
}

var db *sql.DB // 声明数据库连接变量

// ResetDB 删除并重新创建数据库表
func ResetDB() error {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Println("开始重置数据库表...")
	}

	// 删除现有表
	dropTables := []string{
		"DROP TABLE IF EXISTS users",
		"DROP TABLE IF EXISTS devices",
	}

	for _, query := range dropTables {
		if cfg.DebugLevel == "true" {
			log.Printf("执行SQL: %s\n", query)
		}
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("删除表失败: %w", err)
		}
	}

	if cfg.DebugLevel == "true" {
		log.Println("已删除现有表，开始创建新表...")
	}

	// 创建用户表
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		userName VARCHAR(20) NOT NULL COLLATE utf8mb4_unicode_ci,
		passWD VARCHAR(255) NOT NULL,
		userID INT NOT NULL,
		userType INT NOT NULL,
		gatewayDeviceID VARCHAR(12) NOT NULL,
		status INT NULL,
		onlineDuration INT NULL DEFAULT 0,
		certID VARCHAR(64) NULL,
		keyID VARCHAR(64) NULL,
		email VARCHAR(32) NULL COLLATE utf8mb4_unicode_ci,
		permissionMask CHAR(8) NULL,
		lastLoginTimeStamp DATETIME(3) NULL,
		offLineTimeStamp DATETIME(3) NULL,
		loginIP CHAR(24) NULL,
		illegalLoginTimes INT NULL,
		created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3),
		INDEX idx_userid (userID),
		INDEX idx_username (userName),
		INDEX idx_email (email),
		INDEX idx_gatewaydeviceid (gatewayDeviceID),
		FOREIGN KEY (gatewayDeviceID) REFERENCES devices(deviceID) ON DELETE CASCADE ON UPDATE CASCADE
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`

	// 创建设备表
	createDevicesTable := `
	CREATE TABLE IF NOT EXISTS devices (
		id INT AUTO_INCREMENT PRIMARY KEY,
		deviceName VARCHAR(50) NOT NULL COLLATE utf8mb4_unicode_ci,
		deviceType INT NOT NULL,
		passWD VARCHAR(255) NOT NULL,
		deviceID CHAR(12) NOT NULL,
		superiorDeviceID CHAR(12) NOT NULL,
		deviceStatus INT NULL DEFAULT 2,
		peakCPUUsage INT NULL DEFAULT 0,
		peakMemoryUsage INT NULL DEFAULT 0,
		onlineDuration INT NULL DEFAULT 0,
		certID VARCHAR(64) NULL,
		keyID VARCHAR(64) NULL,
		registerIP VARCHAR(24) NULL,
		email VARCHAR(32) NULL COLLATE utf8mb4_unicode_ci,
		deviceHardwareFingerprint CHAR(128) NULL,
		anonymousUser VARCHAR(50) NULL,
		created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3),
		INDEX idx_deviceid (deviceID),
		INDEX idx_devicename (deviceName),
		INDEX idx_superiordeviceid (superiorDeviceID)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`

	// 先创建设备表，因为用户表有外键引用
	if _, err := db.Exec(createDevicesTable); err != nil {
		return fmt.Errorf("创建设备表失败: %w", err)
	}

	// 再创建用户表
	if _, err := db.Exec(createUsersTable); err != nil {
		return fmt.Errorf("创建用户表失败: %w", err)
	}

	if cfg.DebugLevel == "true" {
		log.Println("数据库表重置成功！")
	}
	return nil
}

// InitDB 初始化数据库连接
func InitDB() {
	cfg := config.GetConfig() // 获取全局配置

	var err error
	// 连接数据库，DSN 格式为 "用户名:密码@tcp(主机:端口)/数据库名"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local&allowNativePasswords=true",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName)

	db, err = sql.Open("mysql", dsn) // 打开数据库连接
	if err != nil {
		log.Fatal(err) // 如果连接失败，记录错误并退出
	}

	// 测试连接
	if err = db.Ping(); err != nil {
		log.Fatal(err) // 如果连接失败，记录错误并退出
	}

	if cfg.DebugLevel == "true" {
		log.Println("主数据库连接成功！") // 输出连接成功信息
	}

	// 检查并创建表（如果不存在）
	if err := ensureTablesExist(); err != nil {
		log.Fatal("确保数据库表存在时发生错误:", err)
	}
}

// ensureTablesExist 确保必要的表存在
func ensureTablesExist() error {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Println("检查并创建必要的数据库表...")
	}

	// 创建设备表
	createDevicesTable := `
	CREATE TABLE IF NOT EXISTS devices (
		id INT AUTO_INCREMENT PRIMARY KEY,
		deviceName VARCHAR(50) NOT NULL COLLATE utf8mb4_unicode_ci,
		deviceType INT NOT NULL,
		passWD VARCHAR(255) NOT NULL,
		deviceID CHAR(12) NOT NULL,
		superiorDeviceID CHAR(12) NOT NULL,
		deviceStatus INT NULL DEFAULT 2,
		peakCPUUsage INT NULL DEFAULT 0,
		peakMemoryUsage INT NULL DEFAULT 0,
		onlineDuration INT NULL DEFAULT 0,
		certID VARCHAR(64) NULL,
		keyID VARCHAR(64) NULL,
		registerIP VARCHAR(24) NULL,
		email VARCHAR(32) NULL COLLATE utf8mb4_unicode_ci,
		deviceHardwareFingerprint CHAR(128) NULL,
		anonymousUser VARCHAR(50) NULL,
		created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3),
		INDEX idx_deviceid (deviceID),
		INDEX idx_devicename (deviceName),
		INDEX idx_superiordeviceid (superiorDeviceID)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`

	// 创建用户表
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		userName VARCHAR(20) NOT NULL COLLATE utf8mb4_unicode_ci,
		passWD VARCHAR(255) NOT NULL,
		userID INT NOT NULL,
		userType INT NOT NULL,
		gatewayDeviceID VARCHAR(12) NOT NULL,
		status INT NULL,
		onlineDuration INT NULL DEFAULT 0,
		certID VARCHAR(64) NULL,
		keyID VARCHAR(64) NULL,
		email VARCHAR(32) NULL COLLATE utf8mb4_unicode_ci,
		permissionMask CHAR(8) NULL,
		lastLoginTimeStamp DATETIME(3) NULL,
		offLineTimeStamp DATETIME(3) NULL,
		loginIP CHAR(24) NULL,
		illegalLoginTimes INT NULL,
		created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3),
		INDEX idx_userid (userID),
		INDEX idx_username (userName),
		INDEX idx_email (email),
		INDEX idx_gatewaydeviceid (gatewayDeviceID),
		FOREIGN KEY (gatewayDeviceID) REFERENCES devices(deviceID) ON DELETE CASCADE ON UPDATE CASCADE
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`

	// 创建用户行为表
	createUserBehaviorTable := `
	CREATE TABLE IF NOT EXISTS user_behaviors (
		behaviorID INT AUTO_INCREMENT PRIMARY KEY,
		userID INT NOT NULL,
		behaviorTime DATETIME(3) NOT NULL,
		behaviorType INT NOT NULL COMMENT '1:发送 2:输出',
		dataType INT NOT NULL COMMENT '1:文件 2:消息',
		dataSize BIGINT NOT NULL,
		created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3),
		INDEX idx_userid (userID),
		INDEX idx_behaviortime (behaviorTime),
		INDEX idx_behaviortype (behaviorType),
		FOREIGN KEY (userID) REFERENCES users(userID) ON DELETE CASCADE ON UPDATE CASCADE
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`

	// 先创建设备表
	if _, err := db.Exec(createDevicesTable); err != nil {
		return fmt.Errorf("创建设备表失败: %w", err)
	}

	// 再创建用户表
	if _, err := db.Exec(createUsersTable); err != nil {
		return fmt.Errorf("创建用户表失败: %w", err)
	}

	// 最后创建用户行为表
	if _, err := db.Exec(createUserBehaviorTable); err != nil {
		return fmt.Errorf("创建用户行为表失败: %w", err)
	}

	if cfg.DebugLevel == "true" {
		log.Println("数据库表检查完成！")
	}
	return nil
}

// GetDB 返回数据库连接
func GetDB() *sql.DB {
	return db // 返回数据库连接
}

// GetAllUsers 获取所有用户的全部信息
func GetAllUsers() ([]User, error) {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Println("开始获取所有用户信息")
	}

	db := GetDB()
	query := `SELECT 
		userName, passWD, userID, userType, gatewayDeviceID,
		status, certID, keyID, email, onlineDuration,
		permissionMask, lastLoginTimeStamp, offLineTimeStamp,
		loginIP, illegalLoginTimes, created_at 
	FROM users`

	rows, err := db.Query(query)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("查询用户信息失败: %v\n", err)
		}
		return nil, fmt.Errorf("查询用户信息失败: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		var status sql.NullInt64
		var permissionMask, lastLoginTimeStamp, offLineTimeStamp, loginIP sql.NullString
		var illegalLoginTimes sql.NullInt64
		var certID, keyID, email sql.NullString

		err := rows.Scan(
			&user.UserName, &user.PassWD, &user.UserID, &user.UserType, &user.GatewayDeviceID,
			&status, &certID, &keyID, &email, &user.OnlineDuration,
			&permissionMask, &lastLoginTimeStamp, &offLineTimeStamp,
			&loginIP, &illegalLoginTimes, &user.CreatedAt,
		)
		if err != nil {
			if cfg.DebugLevel == "true" {
				log.Printf("扫描用户信息失败: %v\n", err)
			}
			return nil, fmt.Errorf("扫描用户信息失败: %w", err)
		}

		// 处理可能为 NULL 的字段
		if status.Valid {
			user.Status.Int64 = status.Int64
			user.Status.Valid = true
		}
		if permissionMask.Valid {
			user.PermissionMask.String = permissionMask.String
			user.PermissionMask.Valid = true
		}
		if lastLoginTimeStamp.Valid {
			user.LastLoginTimeStamp.String = lastLoginTimeStamp.String
			user.LastLoginTimeStamp.Valid = true
		}
		if offLineTimeStamp.Valid {
			user.OffLineTimeStamp.String = offLineTimeStamp.String
			user.OffLineTimeStamp.Valid = true
		}
		if loginIP.Valid {
			user.LoginIP.String = loginIP.String
			user.LoginIP.Valid = true
		}
		if illegalLoginTimes.Valid {
			user.IllegalLoginTimes.Int64 = illegalLoginTimes.Int64
			user.IllegalLoginTimes.Valid = true
		}
		if certID.Valid {
			user.CertID = certID.String
		}
		if keyID.Valid {
			user.KeyID = keyID.String
		}
		if email.Valid {
			user.Email = email.String
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("遍历用户信息时发生错误: %v\n", err)
		}
		return nil, fmt.Errorf("遍历用户信息时发生错误: %w", err)
	}

	if cfg.DebugLevel == "true" {
		log.Printf("成功获取 %d 个用户信息\n", len(users))
	}

	return users, nil
}

// GetAllDevices 获取所有设备的全部信息
func GetAllDevices() ([]Device, error) {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Println("开始获取所有设备信息")
	}

	db := GetDB()
	query := `SELECT 
		deviceName, deviceType, passWD, deviceID, superiorDeviceID,
		deviceStatus, peakCPUUsage, peakMemoryUsage, onlineDuration,
		certID, keyID, registerIP, email,
		deviceHardwareFingerprint, anonymousUser, created_at
	FROM devices`

	rows, err := db.Query(query)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("查询设备信息失败: %v\n", err)
		}
		return nil, fmt.Errorf("查询设备信息失败: %w", err)
	}
	defer rows.Close()

	var devices []Device
	for rows.Next() {
		var device Device
		var deviceHardwareFingerprint, anonymousUser sql.NullString
		var certID, keyID, registerIP, email sql.NullString
		var deviceStatus, peakCPUUsage, peakMemoryUsage, onlineDuration sql.NullInt64

		err := rows.Scan(
			&device.DeviceName, &device.DeviceType, &device.PassWD, &device.DeviceID, &device.SuperiorDeviceID,
			&deviceStatus, &peakCPUUsage, &peakMemoryUsage, &onlineDuration,
			&certID, &keyID, &registerIP, &email,
			&deviceHardwareFingerprint, &anonymousUser, &device.CreatedAt,
		)
		if err != nil {
			if cfg.DebugLevel == "true" {
				log.Printf("扫描设备信息失败: %v\n", err)
			}
			return nil, fmt.Errorf("扫描设备信息失败: %w", err)
		}

		// 处理可能为 NULL 的字段
		if deviceHardwareFingerprint.Valid {
			device.DeviceHardwareFingerprint = &deviceHardwareFingerprint.String
		}
		if anonymousUser.Valid {
			device.AnonymousUser = &anonymousUser.String
		}
		if certID.Valid {
			device.CertID = certID.String
		}
		if keyID.Valid {
			device.KeyID = keyID.String
		}
		if registerIP.Valid {
			device.RegisterIP = registerIP.String
		}
		if email.Valid {
			device.Email = email.String
		}
		if deviceStatus.Valid {
			device.DeviceStatus = int(deviceStatus.Int64)
		}
		if peakCPUUsage.Valid {
			device.PeakCPUUsage = int(peakCPUUsage.Int64)
		}
		if peakMemoryUsage.Valid {
			device.PeakMemoryUsage = int(peakMemoryUsage.Int64)
		}
		if onlineDuration.Valid {
			device.OnlineDuration = int(onlineDuration.Int64)
		}

		devices = append(devices, device)
	}

	if err = rows.Err(); err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("遍历设备信息时发生错误: %v\n", err)
		}
		return nil, fmt.Errorf("遍历设备信息时发生错误: %w", err)
	}

	if cfg.DebugLevel == "true" {
		log.Printf("成功获取 %d 个设备信息\n", len(devices))
	}

	return devices, nil
}

// UpdateUser 更新某个用户的除了用户ID以外的其他信息
func UpdateUser(userID int, user User) (map[string]interface{}, error) {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Printf("开始更新用户信息，用户ID: %d\n", userID)
	}

	// 获取当前用户信息
	var existingUser User
	query := `SELECT 
		userName, passWD, email, userType, gatewayDeviceID,
		status, certID, keyID, onlineDuration,
		permissionMask, lastLoginTimeStamp, offLineTimeStamp,
		loginIP, illegalLoginTimes, created_at 
	FROM users WHERE userID = ?`

	var email, certID, keyID sql.NullString
	var status sql.NullInt64
	var permissionMask, lastLoginTimeStamp, offLineTimeStamp, loginIP sql.NullString
	var illegalLoginTimes sql.NullInt64

	err := db.QueryRow(query, userID).Scan(
		&existingUser.UserName, &existingUser.PassWD, &email,
		&existingUser.UserType, &existingUser.GatewayDeviceID,
		&status, &certID, &keyID,
		&existingUser.OnlineDuration, &permissionMask,
		&lastLoginTimeStamp, &offLineTimeStamp,
		&loginIP, &illegalLoginTimes,
		&existingUser.CreatedAt,
	)

	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("获取用户信息失败: %v\n", err)
		}
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 处理可能为 NULL 的字段
	if email.Valid {
		existingUser.Email = email.String
	}
	if certID.Valid {
		existingUser.CertID = certID.String
	}
	if keyID.Valid {
		existingUser.KeyID = keyID.String
	}
	if status.Valid {
		existingUser.Status.Int64 = status.Int64
		existingUser.Status.Valid = true
	}
	if permissionMask.Valid {
		existingUser.PermissionMask.String = permissionMask.String
		existingUser.PermissionMask.Valid = true
	}
	if lastLoginTimeStamp.Valid {
		existingUser.LastLoginTimeStamp.String = lastLoginTimeStamp.String
		existingUser.LastLoginTimeStamp.Valid = true
	}
	if offLineTimeStamp.Valid {
		existingUser.OffLineTimeStamp.String = offLineTimeStamp.String
		existingUser.OffLineTimeStamp.Valid = true
	}
	if loginIP.Valid {
		existingUser.LoginIP.String = loginIP.String
		existingUser.LoginIP.Valid = true
	}
	if illegalLoginTimes.Valid {
		existingUser.IllegalLoginTimes.Int64 = illegalLoginTimes.Int64
		existingUser.IllegalLoginTimes.Valid = true
	}

	// 构建更新 SQL 语句
	updateFields := []string{}
	updateValues := []interface{}{}
	updatedFields := make(map[string]interface{})

	// 检查并添加必填字段的更新
	if user.UserName != "" && user.UserName != existingUser.UserName {
		updateFields = append(updateFields, "userName=?")
		updateValues = append(updateValues, user.UserName)
		updatedFields["userName"] = user.UserName
	}
	if user.PassWD != "" && user.PassWD != existingUser.PassWD {
		updateFields = append(updateFields, "passWD=?")
		updateValues = append(updateValues, user.PassWD)
		updatedFields["passWD"] = user.PassWD
	}
	if user.UserType != 0 && user.UserType != existingUser.UserType {
		updateFields = append(updateFields, "userType=?")
		updateValues = append(updateValues, user.UserType)
		updatedFields["userType"] = user.UserType
	}
	if user.GatewayDeviceID != "" && user.GatewayDeviceID != existingUser.GatewayDeviceID {
		updateFields = append(updateFields, "gatewayDeviceID=?")
		updateValues = append(updateValues, user.GatewayDeviceID)
		updatedFields["gatewayDeviceID"] = user.GatewayDeviceID
	}

	// 检查并添加可选字段的更新
	if user.Email != "" && user.Email != existingUser.Email {
		updateFields = append(updateFields, "email=?")
		updateValues = append(updateValues, user.Email)
		updatedFields["email"] = user.Email
	}
	if user.Status.Valid {
		updateFields = append(updateFields, "status=?")
		updateValues = append(updateValues, user.Status.Int64)
		updatedFields["status"] = user.Status.Int64
	}
	if user.CertID != "" && user.CertID != existingUser.CertID {
		updateFields = append(updateFields, "certID=?")
		updateValues = append(updateValues, user.CertID)
		updatedFields["certID"] = user.CertID
	}
	if user.KeyID != "" && user.KeyID != existingUser.KeyID {
		updateFields = append(updateFields, "keyID=?")
		updateValues = append(updateValues, user.KeyID)
		updatedFields["keyID"] = user.KeyID
	}
	if user.OnlineDuration != 0 && user.OnlineDuration != existingUser.OnlineDuration {
		updateFields = append(updateFields, "onlineDuration=?")
		updateValues = append(updateValues, user.OnlineDuration)
		updatedFields["onlineDuration"] = user.OnlineDuration
	}
	if user.PermissionMask.Valid {
		updateFields = append(updateFields, "permissionMask=?")
		updateValues = append(updateValues, user.PermissionMask.String)
		updatedFields["permissionMask"] = user.PermissionMask.String
	}
	if user.LastLoginTimeStamp.Valid {
		updateFields = append(updateFields, "lastLoginTimeStamp=?")
		updateValues = append(updateValues, user.LastLoginTimeStamp.String)
		updatedFields["lastLoginTimeStamp"] = user.LastLoginTimeStamp.String
	}
	if user.OffLineTimeStamp.Valid {
		updateFields = append(updateFields, "offLineTimeStamp=?")
		updateValues = append(updateValues, user.OffLineTimeStamp.String)
		updatedFields["offLineTimeStamp"] = user.OffLineTimeStamp.String
	}
	if user.LoginIP.Valid {
		updateFields = append(updateFields, "loginIP=?")
		updateValues = append(updateValues, user.LoginIP.String)
		updatedFields["loginIP"] = user.LoginIP.String
	}
	if user.IllegalLoginTimes.Valid {
		updateFields = append(updateFields, "illegalLoginTimes=?")
		updateValues = append(updateValues, user.IllegalLoginTimes.Int64)
		updatedFields["illegalLoginTimes"] = user.IllegalLoginTimes.Int64
	}

	// 如果没有字段需要更新，直接返回
	if len(updateFields) == 0 {
		if cfg.DebugLevel == "true" {
			log.Println("没有字段需要更新")
		}
		return nil, nil
	}

	// 添加 userID 到更新值的最后
	updateValues = append(updateValues, userID)

	// 构建完整的 SQL 语句
	updateSQL := fmt.Sprintf("UPDATE users SET %s WHERE userID=?", strings.Join(updateFields, ", "))

	// 执行更新
	result, err := db.Exec(updateSQL, updateValues...)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("更新用户信息失败: %v\n", err)
		}
		return nil, fmt.Errorf("更新用户信息失败: %w", err)
	}

	// 检查更新结果
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("获取更新影响行数失败: %v\n", err)
		}
		return nil, fmt.Errorf("获取更新影响行数失败: %w", err)
	}

	if rowsAffected == 0 {
		if cfg.DebugLevel == "true" {
			log.Printf("未找到要更新的用户: %d\n", userID)
		}
		return nil, fmt.Errorf("未找到要更新的用户: %d", userID)
	}

	if cfg.DebugLevel == "true" {
		log.Printf("成功更新用户信息，用户ID: %d\n", userID)
	}

	return updatedFields, nil
}

// UpdateDevice 更新某个设备的除了设备ID以外的其他信息
func UpdateDevice(deviceID string, device Device) (map[string]interface{}, error) {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Printf("开始更新设备信息: %s\n", deviceID)
	}

	db := GetDB()
	var existingDevice Device
	query := `SELECT 
		deviceName, deviceType, passWD, superiorDeviceID,
		deviceStatus, peakCPUUsage, peakMemoryUsage, onlineDuration,
		certID, keyID, registerIP, email,
		deviceHardwareFingerprint, anonymousUser, created_at
	FROM devices WHERE deviceID = ?`

	var deviceHardwareFingerprint, anonymousUser sql.NullString
	var certID, keyID, registerIP, email sql.NullString
	var deviceStatus, peakCPUUsage, peakMemoryUsage, onlineDuration sql.NullInt64

	err := db.QueryRow(query, deviceID).Scan(
		&existingDevice.DeviceName, &existingDevice.DeviceType, &existingDevice.PassWD, &existingDevice.SuperiorDeviceID,
		&deviceStatus, &peakCPUUsage, &peakMemoryUsage, &onlineDuration,
		&certID, &keyID, &registerIP, &email,
		&deviceHardwareFingerprint, &anonymousUser, &existingDevice.CreatedAt,
	)

	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("查询设备信息失败: %v\n", err)
		}
		return nil, fmt.Errorf("查询设备信息失败: %w", err)
	}

	// 处理可能为 NULL 的字段
	if deviceHardwareFingerprint.Valid {
		existingDevice.DeviceHardwareFingerprint = &deviceHardwareFingerprint.String
	}
	if anonymousUser.Valid {
		existingDevice.AnonymousUser = &anonymousUser.String
	}
	if certID.Valid {
		existingDevice.CertID = certID.String
	}
	if keyID.Valid {
		existingDevice.KeyID = keyID.String
	}
	if registerIP.Valid {
		existingDevice.RegisterIP = registerIP.String
	}
	if email.Valid {
		existingDevice.Email = email.String
	}
	if deviceStatus.Valid {
		existingDevice.DeviceStatus = int(deviceStatus.Int64)
	}
	if peakCPUUsage.Valid {
		existingDevice.PeakCPUUsage = int(peakCPUUsage.Int64)
	}
	if peakMemoryUsage.Valid {
		existingDevice.PeakMemoryUsage = int(peakMemoryUsage.Int64)
	}
	if onlineDuration.Valid {
		existingDevice.OnlineDuration = int(onlineDuration.Int64)
	}

	// 构建更新语句
	var updateFields []string
	var updateValues []interface{}
	updatedFields := make(map[string]interface{})

	// 检查并添加必填字段的更新
	if device.DeviceName != "" && device.DeviceName != existingDevice.DeviceName {
		updateFields = append(updateFields, "deviceName=?")
		updateValues = append(updateValues, device.DeviceName)
		updatedFields["deviceName"] = device.DeviceName
	}
	if device.DeviceType != 0 && device.DeviceType != existingDevice.DeviceType {
		updateFields = append(updateFields, "deviceType=?")
		updateValues = append(updateValues, device.DeviceType)
		updatedFields["deviceType"] = device.DeviceType
	}
	if device.PassWD != "" && device.PassWD != existingDevice.PassWD {
		updateFields = append(updateFields, "passWD=?")
		updateValues = append(updateValues, device.PassWD)
		updatedFields["passWD"] = device.PassWD
	}
	if device.SuperiorDeviceID != "" && device.SuperiorDeviceID != existingDevice.SuperiorDeviceID {
		updateFields = append(updateFields, "superiorDeviceID=?")
		updateValues = append(updateValues, device.SuperiorDeviceID)
		updatedFields["superiorDeviceID"] = device.SuperiorDeviceID
	}

	// 检查并添加可选字段的更新
	if device.DeviceStatus != 0 && device.DeviceStatus != existingDevice.DeviceStatus {
		updateFields = append(updateFields, "deviceStatus=?")
		updateValues = append(updateValues, device.DeviceStatus)
		updatedFields["deviceStatus"] = device.DeviceStatus
	}
	if device.PeakCPUUsage != 0 && device.PeakCPUUsage != existingDevice.PeakCPUUsage {
		updateFields = append(updateFields, "peakCPUUsage=?")
		updateValues = append(updateValues, device.PeakCPUUsage)
		updatedFields["peakCPUUsage"] = device.PeakCPUUsage
	}
	if device.PeakMemoryUsage != 0 && device.PeakMemoryUsage != existingDevice.PeakMemoryUsage {
		updateFields = append(updateFields, "peakMemoryUsage=?")
		updateValues = append(updateValues, device.PeakMemoryUsage)
		updatedFields["peakMemoryUsage"] = device.PeakMemoryUsage
	}
	if device.OnlineDuration != 0 && device.OnlineDuration != existingDevice.OnlineDuration {
		updateFields = append(updateFields, "onlineDuration=?")
		updateValues = append(updateValues, device.OnlineDuration)
		updatedFields["onlineDuration"] = device.OnlineDuration
	}
	if device.CertID != "" && device.CertID != existingDevice.CertID {
		updateFields = append(updateFields, "certID=?")
		updateValues = append(updateValues, device.CertID)
		updatedFields["certID"] = device.CertID
	}
	if device.KeyID != "" && device.KeyID != existingDevice.KeyID {
		updateFields = append(updateFields, "keyID=?")
		updateValues = append(updateValues, device.KeyID)
		updatedFields["keyID"] = device.KeyID
	}
	if device.RegisterIP != "" && device.RegisterIP != existingDevice.RegisterIP {
		updateFields = append(updateFields, "registerIP=?")
		updateValues = append(updateValues, device.RegisterIP)
		updatedFields["registerIP"] = device.RegisterIP
	}
	if device.Email != "" && device.Email != existingDevice.Email {
		updateFields = append(updateFields, "email=?")
		updateValues = append(updateValues, device.Email)
		updatedFields["email"] = device.Email
	}
	if device.DeviceHardwareFingerprint != nil {
		updateFields = append(updateFields, "deviceHardwareFingerprint=?")
		updateValues = append(updateValues, *device.DeviceHardwareFingerprint)
		updatedFields["deviceHardwareFingerprint"] = *device.DeviceHardwareFingerprint
	}
	if device.AnonymousUser != nil {
		updateFields = append(updateFields, "anonymousUser=?")
		updateValues = append(updateValues, *device.AnonymousUser)
		updatedFields["anonymousUser"] = *device.AnonymousUser
	}

	// 如果没有需要更新的字段，直接返回
	if len(updateFields) == 0 {
		if cfg.DebugLevel == "true" {
			log.Println("没有需要更新的字段")
		}
		return updatedFields, nil
	}

	// 构建完整的更新语句
	updateValues = append(updateValues, deviceID)
	query = fmt.Sprintf("UPDATE devices SET %s WHERE deviceID=?", strings.Join(updateFields, ", "))

	// 执行更新
	result, err := db.Exec(query, updateValues...)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("更新设备信息失败: %v\n", err)
		}
		return nil, fmt.Errorf("更新设备信息失败: %w", err)
	}

	// 检查更新结果
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("获取更新影响行数失败: %v\n", err)
		}
		return nil, fmt.Errorf("获取更新影响行数失败: %w", err)
	}

	if rowsAffected == 0 {
		if cfg.DebugLevel == "true" {
			log.Printf("未找到要更新的设备: %s\n", deviceID)
		}
		return nil, fmt.Errorf("未找到要更新的设备: %s", deviceID)
	}

	if cfg.DebugLevel == "true" {
		log.Printf("成功更新设备信息: %s\n", deviceID)
	}

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
