package testdata

import (
	"fmt"
	"gin-server/config"
	"gin-server/database/models"
	"log"

	"gorm.io/gorm"
)

// DeviceGenerator 设备数据生成器
type DeviceGenerator struct {
	*BaseGenerator
}

// NewDeviceGenerator 创建设备数据生成器
func NewDeviceGenerator(cfg *config.Config) *DeviceGenerator {
	return &DeviceGenerator{
		BaseGenerator: NewBaseGenerator(cfg),
	}
}

// Generate 生成设备测试数据
func (g *DeviceGenerator) Generate(db *gorm.DB, count int) error {
	// 检查是否已存在设备数据
	exists, err := IsDataExists(db, &models.Device{})
	if err != nil {
		return fmt.Errorf("检查设备数据失败: %w", err)
	}

	// 如果设备表中已有数据，且设置了不覆盖数据，则跳过
	if exists {
		g.LogInfo("设备数据已存在，跳过生成")
		return nil
	}

	g.LogInfo("开始生成设备测试数据，数量: %d", count)

	// 创建根设备 (安全接入管理设备，设备类型为4，上级设备ID为0)
	rootDevice := &models.Device{
		DeviceName:       "安全接入管理设备",
		DeviceType:       4,
		Password:         "admin123456",
		DeviceID:         1000, // 给根设备一个特定ID
		SuperiorDeviceID: 0,    // 根设备没有上级
		DeviceStatus:     1,    // 1表示在线
		RegisterIP:       g.RandomIP(),
		Email:            g.RandomEmail(),
		// 安全接入管理设备这三个字段为空
		LongAddress:  "",
		ShortAddress: "",
		SESKey:       "",
	}

	if err := db.Create(rootDevice).Error; err != nil {
		return fmt.Errorf("创建根设备失败: %w", err)
	}

	// 生成其他设备（网关设备，设备类型为1-3，上级设备ID指向根设备）
	for i := 0; i < count-1; i++ {
		deviceID := 1001 + i            // 从1001开始，避免与根设备ID冲突
		deviceType := g.RandomInt(1, 3) // 随机网关设备类型1-3
		deviceName := fmt.Sprintf("网关设备-%d-%d", deviceType, i+1)

		// 生成IPv6格式的长地址
		longAddress := g.RandomIPv6()

		// 生成2字节的短地址（以十六进制格式表示）
		shortAddress := fmt.Sprintf("%02X%02X",
			g.RandomInt(0, 255), g.RandomInt(0, 255))

		// 生成SES密钥
		sesKey := g.RandomString(16)

		device := &models.Device{
			DeviceName:          deviceName,
			DeviceType:          deviceType,
			Password:            "device" + g.RandomString(6),
			DeviceID:            deviceID,
			SuperiorDeviceID:    1000,              // 指向根设备
			DeviceStatus:        g.RandomInt(1, 2), // 状态随机在线/离线
			PeakCPUUsage:        g.RandomInt(10, 90),
			PeakMemoryUsage:     g.RandomInt(20, 80),
			OnlineDuration:      g.RandomInt(100, 10000),
			RegisterIP:          g.RandomIP(),
			Email:               g.RandomEmail(),
			HardwareFingerprint: g.RandomString(32),
			// 设置网关设备的三个新字段
			LongAddress:  longAddress,
			ShortAddress: shortAddress,
			SESKey:       sesKey,
		}

		if err := db.Create(device).Error; err != nil {
			log.Printf("创建设备 %s 失败: %v", deviceName, err)
			continue
		}
	}

	g.LogInfo("设备测试数据生成完成，成功创建 %d 个设备", count)
	return nil
}
