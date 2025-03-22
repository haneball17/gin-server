package database

import (
	"fmt"
	"gin-server/config"
	"gin-server/database/connections"
	"gin-server/database/migrations"
	"log"
	"sync"

	"gorm.io/gorm"
)

var (
	// 单例模式实现
	manager *connections.ConnectionManager
	once    sync.Once
)

// Initialize 初始化数据库连接管理器和表结构
// 此函数将完成以下工作：
// 1. 创建数据库连接管理器
// 2. 检查并创建所有必要的表结构
// 3. 执行必要的迁移任务
func Initialize(cfg *config.Config) error {
	// 初始化连接管理器（使用once确保只初始化一次）
	once.Do(func() {
		manager = connections.NewConnectionManager(cfg)
	})

	// 获取数据库连接
	db, err := GetDB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %w", err)
	}

	// 执行自动迁移（包含所有模型的表创建/更新）
	if err := migrations.AutoMigrate(db); err != nil {
		log.Printf("数据库迁移失败: %v", err)

		// 尝试错误恢复 - 确保关键表存在
		if recoverErr := ensureAllCriticalTablesExist(db); recoverErr != nil {
			log.Printf("错误恢复失败，无法确保关键表存在: %v", recoverErr)
			// 返回原始错误和恢复错误的组合信息
			return fmt.Errorf("数据库迁移失败: %w (恢复也失败: %v)", err, recoverErr)
		}

		// 恢复成功，但仍返回原始错误，提示用户迁移未完全成功
		log.Printf("尽管迁移失败，但已确保关键表存在，系统可能部分功能受限")
		return fmt.Errorf("数据库迁移部分失败: %w (已执行恢复措施)", err)
	}

	if cfg.DebugLevel == "true" {
		log.Println("数据库迁移成功")
	}

	// 特别确保关键表存在（如 log_files 表）
	// 此步骤是显式的，以确保关键表即使在 AutoMigrate 失败的情况下也能创建
	if err := ensureAllCriticalTablesExist(db); err != nil {
		log.Printf("确保关键表存在失败: %v", err)
		return fmt.Errorf("确保关键表存在失败: %w", err)
	}

	if cfg.DebugLevel == "true" {
		log.Println("数据库初始化完成")
	}

	return nil
}

// ensureAllCriticalTablesExist 确保所有关键表存在
// 关键表是指系统运行必须依赖的表
func ensureAllCriticalTablesExist(db *gorm.DB) error {
	// 使用新的表检查机制
	return migrations.EnsureAllTablesExist(db)
}

// GetManager 获取连接管理器实例
func GetManager() *connections.ConnectionManager {
	if manager == nil {
		// 如果管理器未初始化，使用全局配置初始化
		cfg := config.GetConfig()
		err := Initialize(cfg)
		if err != nil {
			// 记录错误，但尝试继续（可能只有部分功能可用）
			log.Printf("警告：数据库初始化失败：%v，部分功能可能无法正常工作", err)
		}
	}
	return manager
}

// GetDB 获取默认数据库连接的GORM实例
func GetDB() (*gorm.DB, error) {
	conn, err := GetManager().Default()
	if err != nil {
		return nil, err
	}
	return conn.GetDB(), nil
}

// GetRadiusDB 获取Radius数据库连接的GORM实例
func GetRadiusDB() (*gorm.DB, error) {
	conn, err := GetManager().Radius()
	if err != nil {
		return nil, err
	}
	return conn.GetDB(), nil
}

// CloseAll 关闭所有数据库连接
func CloseAll() error {
	if manager != nil {
		return manager.CloseAll()
	}
	return nil
}

// 保留这个函数作为向后兼容，以防其他模块已经在使用它
// 但实际上它现在会在Initialize中自动调用
func InitLogFilesTable() error {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Println("通过兼容接口初始化log_files表...")
	}

	db, err := GetDB()
	if err != nil {
		log.Printf("获取数据库连接失败: %v", err)
		return err
	}

	if err := migrations.EnsureLogFilesTableExists(db); err != nil {
		log.Printf("确保log_files表存在失败: %v", err)
		return err
	}

	if cfg.DebugLevel == "true" {
		log.Println("log_files表初始化成功")
	}

	return nil
}
