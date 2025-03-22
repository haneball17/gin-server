package models

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型，包含共享字段
type BaseModel struct {
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deleted_at;index"`
}

// Model 模型接口
type Model interface {
	TableName() string
}

// TableName 获取表名
func (m *BaseModel) TableName() string {
	return ""
}
