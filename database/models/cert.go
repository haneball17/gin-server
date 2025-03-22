package models

import (
	"time"

	"gorm.io/gorm"
)

// Cert 证书信息
type Cert struct {
	gorm.Model
	EntityType string    `json:"entity_type" gorm:"column:entity_type;not null;index:idx_entity;type:varchar(32)"` // 实体类型：user或device
	EntityID   string    `json:"entity_id" gorm:"column:entity_id;not null;index:idx_entity;type:varchar(128)"`    // 用户ID或设备ID
	CertPath   string    `json:"cert_path" gorm:"column:cert_path;type:varchar(255)"`                              // 证书文件路径
	KeyPath    string    `json:"key_path" gorm:"column:key_path;type:varchar(255)"`                                // 密钥文件路径
	UploadTime time.Time `json:"upload_time" gorm:"column:upload_time;not null"`                                   // 上传时间
}

// TableName 指定表名
func (Cert) TableName() string {
	return "certs"
}
