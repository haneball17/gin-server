package testdata

import (
	"fmt"
	"gin-server/config"
	"gin-server/database/models"
	"log"
	"time"

	"gorm.io/gorm"
)

// UserGenerator 用户数据生成器
type UserGenerator struct {
	*BaseGenerator
}

// NewUserGenerator 创建用户数据生成器
func NewUserGenerator(cfg *config.Config) *UserGenerator {
	return &UserGenerator{
		BaseGenerator: NewBaseGenerator(cfg),
	}
}

// Generate 生成用户测试数据
// 参数 count 表示每个设备需要生成的用户数
func (g *UserGenerator) Generate(db *gorm.DB, count int) error {
	// 检查是否已存在用户数据
	exists, err := IsDataExists(db, &models.User{})
	if err != nil {
		return fmt.Errorf("检查用户数据失败: %w", err)
	}

	// 如果用户表中已有数据，则跳过
	if exists {
		g.LogInfo("用户数据已存在，跳过生成")
		return nil
	}

	// 获取所有网关设备（设备类型1-3）
	var devices []models.Device
	if err := db.Where("device_type IN ?", []int{1, 2, 3}).Find(&devices).Error; err != nil {
		return fmt.Errorf("获取设备列表失败: %w", err)
	}

	if len(devices) == 0 {
		return fmt.Errorf("没有找到可用的网关设备，请先生成设备数据")
	}

	g.LogInfo("开始为 %d 个网关设备生成用户数据，每个设备 %d 个用户", len(devices), count)

	// 为每个网关设备生成用户
	totalUsers := 0
	startUserID := 10000 // 用户ID从10000开始

	for _, device := range devices {
		for i := 0; i < count; i++ {
			userID := startUserID + totalUsers
			userName := fmt.Sprintf("user_%d_%d", device.DeviceID, i+1)

			// 随机生成用户状态
			var status *int
			if g.RandomBool() {
				statusValue := g.RandomInt(1, 4) // 1:在线，2:离线，3:冻结，4:注销
				status = &statusValue
			}

			// 随机生成登录和离线时间
			var lastLogin, offlineTime *time.Time
			now := time.Now()
			if g.RandomBool() {
				loginTime := now.Add(-time.Duration(g.RandomInt(1, 72)) * time.Hour)
				lastLogin = &loginTime

				if g.RandomBool() && status != nil && *status == 2 { // 如果状态是离线
					offline := loginTime.Add(time.Duration(g.RandomInt(1, 24)) * time.Hour)
					if offline.Before(now) { // 确保离线时间不在未来
						offlineTime = &offline
					}
				}
			}

			// 随机生成非法登录次数
			var illegalLogins *int
			if g.RandomBool() {
				illegalValue := g.RandomInt(0, 5)
				illegalLogins = &illegalValue
			}

			user := &models.User{
				Username:           userName,
				Password:           "pass" + g.RandomString(8),
				UserID:             userID,
				UserType:           g.RandomInt(1, 3), // 随机用户类型1-3
				GatewayDeviceID:    device.DeviceID,
				Status:             status,
				OnlineDuration:     g.RandomInt(0, 1000),
				Email:              g.RandomEmail(),
				PermissionMask:     fmt.Sprintf("%08b", g.RandomInt(0, 255)), // 随机8位二进制权限
				LastLoginTimeStamp: lastLogin,
				OffLineTimeStamp:   offlineTime,
				LoginIP:            g.RandomIP(),
				IllegalLoginTimes:  illegalLogins,
			}

			if err := db.Create(user).Error; err != nil {
				log.Printf("创建用户 %s 失败: %v", userName, err)
				continue
			}

			totalUsers++
		}
	}

	g.LogInfo("用户测试数据生成完成，成功创建 %d 个用户", totalUsers)
	return nil
}
