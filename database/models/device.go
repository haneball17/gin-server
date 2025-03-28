package models

import (
	"gorm.io/gorm"
)

// Device 设备信息
type Device struct {
	gorm.Model
	DeviceName          string `json:"device_name" gorm:"column:device_name;not null;type:varchar(128)"`
	DeviceType          int    `json:"device_type" gorm:"column:device_type;not null"`
	Password            string `json:"password" gorm:"column:pass_wd;not null;type:varchar(128)"`
	DeviceID            int    `json:"device_id" gorm:"column:device_id;uniqueIndex;not null"`
	SuperiorDeviceID    int    `json:"superior_device_id" gorm:"column:superior_device_id"`
	DeviceStatus        int    `json:"device_status" gorm:"column:device_status;default:2"` // 默认离线状态
	PeakCPUUsage        int    `json:"peak_cpu_usage" gorm:"column:peak_cpu_usage;default:0"`
	PeakMemoryUsage     int    `json:"peak_memory_usage" gorm:"column:peak_memory_usage;default:0"`
	OnlineDuration      int    `json:"online_duration" gorm:"column:online_duration;default:0"`
	CertID              string `json:"cert_id" gorm:"column:cert_id;type:varchar(255)"`
	KeyID               string `json:"key_id" gorm:"column:key_id;type:varchar(255)"`
	RegisterIP          string `json:"register_ip" gorm:"column:register_ip;type:varchar(64)"`
	Email               string `json:"email" gorm:"column:email;type:varchar(128)"`
	HardwareFingerprint string `json:"hardware_fingerprint" gorm:"column:hardware_fingerprint;type:varchar(255)"`
	AnonymousUser       string `json:"anonymous_user" gorm:"column:anonymous_user;type:varchar(128)"`
	LongAddress         string `json:"long_address" gorm:"column:long_address;type:varchar(255)"`
	ShortAddress        string `json:"short_address" gorm:"column:short_address;type:varchar(255)"`
	SESKey              string `json:"ses_key" gorm:"column:ses_key;type:varchar(255)"`
}

// TableName 指定表名
func (Device) TableName() string {
	return "devices"
}
