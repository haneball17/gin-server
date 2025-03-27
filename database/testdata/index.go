package testdata

import (
	"gin-server/config"
	"log"
	"sync"

	"gorm.io/gorm"
)

var (
	// 单例模式实现
	manager *TestDataManager
	once    sync.Once
)

// InitializeManager 初始化测试数据管理器
func InitializeManager(cfg *config.Config, db *gorm.DB) {
	once.Do(func() {
		manager = NewTestDataManager(cfg, db)
	})
}

// GetManager 获取测试数据管理器实例
func GetManager() *TestDataManager {
	if manager == nil {
		log.Println("警告: 测试数据管理器尚未初始化")
	}
	return manager
}

// Initialize 初始化测试数据
// 此函数应在程序启动时由main函数调用
func Initialize(cfg *config.Config, db *gorm.DB) error {
	// 初始化管理器
	InitializeManager(cfg, db)

	// 初始化测试数据
	if err := manager.InitializeTestData(); err != nil {
		return err
	}

	// 如果开启了实时数据生成，则启动实时数据生成
	if cfg.TestData.EnableRealtimeData {
		if err := manager.StartRealtimeDataGeneration(); err != nil {
			return err
		}
	}

	return nil
}

// Shutdown 关闭测试数据模块
// 此函数应在程序退出时由main函数调用
func Shutdown() {
	if manager != nil {
		manager.StopRealtimeDataGeneration()
	}
}
