package repositories

import (
	"gorm.io/gorm"
)

// Repository 仓库接口
type Repository interface {
	// WithTx 使用事务进行操作
	WithTx(tx *gorm.DB) Repository
}

// BaseRepository 基础仓库实现
type BaseRepository struct {
	db *gorm.DB
}

// NewBaseRepository 创建基础仓库实例
func NewBaseRepository(db *gorm.DB) *BaseRepository {
	return &BaseRepository{
		db: db,
	}
}

// WithTx 使用事务进行操作
func (r *BaseRepository) WithTx(tx *gorm.DB) *BaseRepository {
	return &BaseRepository{
		db: tx,
	}
}

// GetDB 获取数据库实例
func (r *BaseRepository) GetDB() *gorm.DB {
	return r.db
}
