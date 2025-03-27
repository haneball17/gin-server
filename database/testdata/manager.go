package testdata

import (
	"fmt"
	"gin-server/config"
	"log"
	"sync"
	"time"

	"gorm.io/gorm"
)

// TestDataManager 测试数据管理器
type TestDataManager struct {
	cfg               *config.Config
	db                *gorm.DB
	deviceGenerator   *DeviceGenerator
	userGenerator     *UserGenerator
	behaviorGenerator *BehaviorGenerator
	realtimeRunning   bool
	realtimeStopChan  chan struct{}
	realtimeWaitGroup sync.WaitGroup
}

// NewTestDataManager 创建测试数据管理器
func NewTestDataManager(cfg *config.Config, db *gorm.DB) *TestDataManager {
	return &TestDataManager{
		cfg:               cfg,
		db:                db,
		deviceGenerator:   NewDeviceGenerator(cfg),
		userGenerator:     NewUserGenerator(cfg),
		behaviorGenerator: NewBehaviorGenerator(cfg),
		realtimeStopChan:  make(chan struct{}),
	}
}

// InitializeTestData 初始化测试数据
func (m *TestDataManager) InitializeTestData() error {
	log.Println("检查是否需要初始化测试数据...")

	// 检查配置是否启用初始数据生成
	if !m.cfg.TestData.EnableInitialData {
		log.Println("初始测试数据生成未启用，跳过")
		return nil
	}

	log.Println("开始生成初始测试数据...")

	// 1. 先生成设备数据
	if err := m.deviceGenerator.Generate(m.db, m.cfg.TestData.DeviceCount); err != nil {
		log.Printf("生成设备数据失败: %v", err)
		return fmt.Errorf("生成设备数据失败: %w", err)
	}

	// 2. 基于设备生成用户数据
	if err := m.userGenerator.Generate(m.db, m.cfg.TestData.UsersPerDevice); err != nil {
		log.Printf("生成用户数据失败: %v", err)
		return fmt.Errorf("生成用户数据失败: %w", err)
	}

	// 3. 基于用户生成行为数据
	if err := m.behaviorGenerator.Generate(m.db, m.cfg.TestData.BehaviorsPerUser); err != nil {
		log.Printf("生成用户行为数据失败: %v", err)
		return fmt.Errorf("生成用户行为数据失败: %w", err)
	}

	log.Println("初始测试数据生成完成")
	return nil
}

// StartRealtimeDataGeneration 启动实时数据生成
func (m *TestDataManager) StartRealtimeDataGeneration() error {
	// 检查配置是否启用实时数据生成
	if !m.cfg.TestData.EnableRealtimeData {
		log.Println("实时测试数据生成未启用，跳过")
		return nil
	}

	// 防止重复启动
	if m.realtimeRunning {
		log.Println("实时数据生成已经在运行中")
		return nil
	}

	log.Println("启动实时测试数据生成...")
	m.realtimeRunning = true

	// 启动goroutine执行定时生成
	m.realtimeWaitGroup.Add(1)
	go m.realtimeDataGenerationLoop()

	log.Println("实时测试数据生成已启动")
	return nil
}

// StopRealtimeDataGeneration 停止实时数据生成
func (m *TestDataManager) StopRealtimeDataGeneration() {
	if !m.realtimeRunning {
		return
	}

	log.Println("正在停止实时测试数据生成...")
	close(m.realtimeStopChan)
	m.realtimeWaitGroup.Wait()

	// 重置状态
	m.realtimeRunning = false
	m.realtimeStopChan = make(chan struct{})

	log.Println("实时测试数据生成已停止")
}

// realtimeDataGenerationLoop 实时数据生成循环
func (m *TestDataManager) realtimeDataGenerationLoop() {
	defer m.realtimeWaitGroup.Done()

	// 获取配置的时间间隔（秒）
	interval := time.Duration(m.cfg.TestData.RealtimeInterval) * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Printf("实时数据生成循环已启动，间隔: %v", interval)

	for {
		select {
		case <-ticker.C:
			// 执行实时数据生成
			if err := m.behaviorGenerator.GenerateRealtime(
				m.db,
				m.cfg.TestData.RealtimeBehaviorsPerInterval,
			); err != nil {
				log.Printf("生成实时行为数据失败: %v", err)
			}
		case <-m.realtimeStopChan:
			log.Println("收到停止信号，实时数据生成循环退出")
			return
		}
	}
}
