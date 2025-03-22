package repositories

import (
	"gorm.io/gorm"
)

// RepositoryFactory 仓库工厂接口
type RepositoryFactory interface {
	// GetEventRepository 获取事件仓库
	GetEventRepository() EventRepository

	// GetLogFileRepository 获取日志文件仓库
	GetLogFileRepository() LogFileRepository

	// GetUserBehaviorRepository 获取用户行为仓库
	GetUserBehaviorRepository() UserBehaviorRepository

	// GetUserRepository 获取用户仓库
	GetUserRepository() UserRepository

	// GetDeviceRepository 获取设备仓库
	GetDeviceRepository() DeviceRepository

	// GetCertRepository 获取证书仓库
	GetCertRepository() CertRepository

	// GetRadiusAuthRepository 获取Radius认证仓库
	GetRadiusAuthRepository() RadiusAuthRepository

	// WithTx 使用事务创建仓库工厂
	WithTx(tx *gorm.DB) RepositoryFactory
}

// repositoryFactory 仓库工厂实现
type repositoryFactory struct {
	db *gorm.DB
}

// NewRepositoryFactory 创建仓库工厂实例
func NewRepositoryFactory(db *gorm.DB) RepositoryFactory {
	return &repositoryFactory{
		db: db,
	}
}

// GetEventRepository 获取事件仓库
func (f *repositoryFactory) GetEventRepository() EventRepository {
	return NewEventRepository(f.db)
}

// GetLogFileRepository 获取日志文件仓库
func (f *repositoryFactory) GetLogFileRepository() LogFileRepository {
	return NewLogFileRepository(f.db)
}

// GetUserBehaviorRepository 获取用户行为仓库
func (f *repositoryFactory) GetUserBehaviorRepository() UserBehaviorRepository {
	return NewUserBehaviorRepository(f.db)
}

// GetUserRepository 获取用户仓库
func (f *repositoryFactory) GetUserRepository() UserRepository {
	return NewUserRepository(f.db)
}

// GetDeviceRepository 获取设备仓库
func (f *repositoryFactory) GetDeviceRepository() DeviceRepository {
	return NewDeviceRepository(f.db)
}

// GetCertRepository 获取证书仓库
func (f *repositoryFactory) GetCertRepository() CertRepository {
	return NewCertRepository(f.db)
}

// GetRadiusAuthRepository 获取Radius认证仓库
func (f *repositoryFactory) GetRadiusAuthRepository() RadiusAuthRepository {
	return NewRadiusAuthRepository(f.db)
}

// WithTx 使用事务创建仓库工厂
func (f *repositoryFactory) WithTx(tx *gorm.DB) RepositoryFactory {
	return &repositoryFactory{
		db: tx,
	}
}
