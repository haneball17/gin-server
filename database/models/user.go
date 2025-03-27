package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户信息
type User struct {
	gorm.Model
	Username           string     `json:"username" gorm:"column:user_name;not null;type:varchar(64)"` // 用户名
	Password           string     `json:"-" gorm:"column:pass_wd;not null;type:varchar(128)"`         // 密码
	UserID             int        `json:"user_id" gorm:"column:user_id;not null;uniqueIndex"`         // 用户唯一标识
	UserType           int        `json:"user_type" gorm:"column:user_type;not null"`                 // 用户类型
	GatewayDeviceID    int        `json:"gateway_device_id" gorm:"column:gateway_device_id;not null"` // 用户所属网关设备ID
	Status             *int       `json:"status" gorm:"column:status;default:null"`                   // 用户状态，1:在线，2:离线，3:冻结，4:注销
	OnlineDuration     int        `json:"online_duration" gorm:"column:online_duration;default:0"`    // 在线时长
	CertID             string     `json:"cert_id" gorm:"column:cert_id;type:varchar(255)"`            // 证书ID
	KeyID              string     `json:"key_id" gorm:"column:key_id;type:varchar(255)"`              // 密钥ID
	Email              string     `json:"email" gorm:"column:email;type:varchar(128)"`                // 邮箱
	PermissionMask     string     `json:"permission_mask" gorm:"column:permission_mask;type:char(8)"` // 权限位掩码
	LastLoginTimeStamp *time.Time `json:"last_login_timestamp" gorm:"column:last_login_time_stamp"`   // 登录时间戳
	OffLineTimeStamp   *time.Time `json:"offline_timestamp" gorm:"column:off_line_time_stamp"`        // 离线时间戳
	LoginIP            string     `json:"login_ip" gorm:"column:login_ip;type:char(24)"`              // 用户登录IP
	IllegalLoginTimes  *int       `json:"illegal_login_times" gorm:"column:illegal_login_times"`      // 用户本次的非法登录次数
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
