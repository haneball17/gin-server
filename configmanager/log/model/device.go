package model

import "time"

// Device 设备信息
type Device struct {
	ID                  int       `json:"id" gorm:"column:id"`
	DeviceName          string    `json:"deviceName" gorm:"column:deviceName"`
	DeviceType          int       `json:"deviceType" gorm:"column:deviceType"`
	Password            string    `json:"password" gorm:"column:passWD"`
	DeviceID            string    `json:"deviceId" gorm:"column:deviceID"`
	SuperiorDeviceID    string    `json:"superiorDeviceId" gorm:"column:superiorDeviceID"`
	DeviceStatus        int       `json:"deviceStatus" gorm:"column:deviceStatus"`
	PeakCPUUsage        int       `json:"peakCPUUsage" gorm:"column:peakCPUUsage"`
	PeakMemoryUsage     int       `json:"peakMemoryUsage" gorm:"column:peakMemoryUsage"`
	OnlineDuration      int       `json:"onlineDuration" gorm:"column:onlineDuration"`
	CertID              string    `json:"certId" gorm:"column:certID"`
	KeyID               string    `json:"keyId" gorm:"column:keyID"`
	RegisterIP          string    `json:"registerIP" gorm:"column:registerIP"`
	Email               string    `json:"email" gorm:"column:email"`
	HardwareFingerprint string    `json:"hardwareFingerprint" gorm:"column:deviceHardwareFingerprint"`
	AnonymousUser       string    `json:"anonymousUser" gorm:"column:anonymousUser"`
	CreatedAt           time.Time `json:"createdAt" gorm:"column:created_at"`
}

// TableName 指定表名
func (Device) TableName() string {
	return "devices"
}
