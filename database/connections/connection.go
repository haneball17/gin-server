package connections

import (
	"gorm.io/gorm"
)

// Connection 数据库连接接口
type Connection interface {
	// GetDB 获取GORM数据库实例
	GetDB() *gorm.DB

	// Close 关闭数据库连接
	Close() error
}
