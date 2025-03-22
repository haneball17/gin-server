package model

import (
	"fmt"
	"gin-server/config"
	"gin-server/database"
	"gin-server/database/models"
	"log"
)

// 兼容旧代码的类型别名
type AuthRecord = models.RadPostAuth
type AuthRecordQuery = models.RadPostAuthQuery

// InitRadiusDB 初始化Radius数据库连接
func InitRadiusDB() error {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Println("初始化Radius数据库连接...")
	}

	// 使用database模块获取Radius数据库连接
	radiusDB, err := database.GetRadiusDB()
	if err != nil {
		return fmt.Errorf("连接Radius数据库失败: %w", err)
	}

	// 测试连接
	sqlDB, err := radiusDB.DB()
	if err != nil {
		return fmt.Errorf("获取底层数据库连接失败: %w", err)
	}

	if err = sqlDB.Ping(); err != nil {
		return fmt.Errorf("Radius数据库Ping失败: %w", err)
	}

	// 确保表结构存在
	if err := EnsureRadiusTablesExist(); err != nil {
		return fmt.Errorf("确保Radius表存在失败: %w", err)
	}

	if cfg.DebugLevel == "true" {
		log.Println("Radius数据库连接成功")
	}
	return nil
}

// EnsureRadiusTablesExist 确保Radius数据库的表存在
func EnsureRadiusTablesExist() error {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Println("检查Radius数据库表...")
	}

	// 获取Radius数据库连接
	radiusDB, err := database.GetRadiusDB()
	if err != nil {
		return fmt.Errorf("获取Radius数据库连接失败: %w", err)
	}

	// 使用GORM自动迁移RadPostAuth表
	if err := radiusDB.AutoMigrate(&models.RadPostAuth{}); err != nil {
		return fmt.Errorf("自动迁移RadPostAuth表失败: %w", err)
	}

	if cfg.DebugLevel == "true" {
		log.Println("Radius数据库表检查完成")
	}

	return nil
}
