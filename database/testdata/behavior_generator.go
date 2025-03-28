package testdata

import (
	"fmt"
	"gin-server/config"
	"gin-server/database/models"
	"time"

	"gorm.io/gorm"
)

// BehaviorGenerator 用户行为数据生成器
type BehaviorGenerator struct {
	*BaseGenerator
}

// NewBehaviorGenerator 创建用户行为数据生成器
func NewBehaviorGenerator(cfg *config.Config) *BehaviorGenerator {
	return &BehaviorGenerator{
		BaseGenerator: NewBaseGenerator(cfg),
	}
}

// Generate 生成用户行为测试数据
// 参数:
//   - db: 数据库连接
//   - count: 每个用户需要生成的行为数据数量
//   - skipIfExists: 是否在数据已存在时跳过生成，默认为true
func (g *BehaviorGenerator) Generate(db *gorm.DB, count int, skipIfExists ...bool) error {
	// 默认在数据存在时跳过
	skip := true
	if len(skipIfExists) > 0 {
		skip = skipIfExists[0]
	}

	// 检查是否已存在用户行为数据
	if skip {
		exists, err := IsDataExists(db, &models.UserBehavior{})
		if err != nil {
			return fmt.Errorf("检查用户行为数据失败: %w", err)
		}

		// 如果用户行为表中已有数据，且设置了跳过，则不生成新数据
		if exists {
			g.LogInfo("用户行为数据已存在，跳过生成")
			return nil
		}
	}

	// 获取所有用户
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		return fmt.Errorf("获取用户列表失败: %w", err)
	}

	if len(users) == 0 {
		return fmt.Errorf("没有找到可用的用户，请先生成用户数据")
	}

	// 确定要生成的行为数据，如果设置了跳过且数据已存在，count为0
	if skip {
		var existingCount int64
		if err := db.Model(&models.UserBehavior{}).Count(&existingCount).Error; err != nil {
			g.LogInfo("统计现有行为数据失败: %v", err)
		} else if existingCount > 0 {
			g.LogInfo("用户行为表中已有 %d 条数据", existingCount)
		}
	}

	g.LogInfo("开始为 %d 个用户生成行为数据，每个用户 %d 条行为记录", len(users), count)

	// 为每个用户生成行为数据
	totalBehaviors := 0
	currentTime := time.Now()

	for _, user := range users {
		// 为每个用户生成指定数量的行为数据
		for i := 0; i < count; i++ {
			// 随机生成行为时间，范围为过去30天内
			behaviorTime := currentTime.Add(-time.Duration(g.RandomInt(1, 30*24)) * time.Hour)
			// 再增加随机分钟数
			behaviorTime = behaviorTime.Add(time.Duration(g.RandomInt(0, 60)) * time.Minute)

			// 随机生成行为类型：1-发送，2-接收
			behaviorType := g.RandomInt(1, 2)

			// 随机生成数据类型：1-文件，2-消息
			dataType := g.RandomInt(1, 2)

			// 随机生成数据大小
			var dataSize int64
			if dataType == 1 { // 文件通常较大
				dataSize = g.RandomInt64(1024, 10*1024*1024) // 1KB - 10MB
			} else { // 消息通常较小
				dataSize = g.RandomInt64(10, 1024) // 10B - 1KB
			}

			behavior := &models.UserBehavior{
				BehaviorID:   0, // 设置为0，MySQL会自动处理自增字段
				UserID:       user.UserID,
				BehaviorTime: behaviorTime,
				BehaviorType: behaviorType,
				DataType:     dataType,
				DataSize:     dataSize,
			}

			if err := db.Create(behavior).Error; err != nil {
				g.LogInfo("创建用户 %d 的行为数据失败: %v", user.UserID, err)
				continue
			}

			totalBehaviors++
		}
	}

	g.LogInfo("用户行为测试数据生成完成，成功创建 %d 条行为记录", totalBehaviors)
	return nil
}

// GenerateRealtime 生成实时用户行为数据
// 为所有用户生成最近时间的行为数据
func (g *BehaviorGenerator) GenerateRealtime(db *gorm.DB, count int) error {
	// 获取当前时间
	currentTime := time.Now()

	// 从配置中获取时间偏移量（分钟）
	startTimeOffset := time.Duration(-g.Cfg.TestData.RealtimeStartTimeOffset) * time.Minute
	endTimeOffset := time.Duration(-g.Cfg.TestData.RealtimeEndTimeOffset) * time.Minute

	// 计算起始时间和结束时间
	startTime := currentTime.Add(startTimeOffset)
	endTime := currentTime.Add(endTimeOffset)

	// 调用 GenerateAccumulated 方法，使用计算得到的时间范围
	return g.GenerateAccumulated(db, count, startTime, endTime)
}

// GenerateAccumulated 生成累积的用户行为数据
// 该方法不检查是否已存在数据，可用于累积生成数据
// 参数:
//   - db: 数据库连接
//   - count: 每个用户生成的数据条数
//   - startTime: 行为时间的起始时间范围
//   - endTime: 行为时间的结束时间范围
func (g *BehaviorGenerator) GenerateAccumulated(db *gorm.DB, count int, startTime, endTime time.Time) error {
	// 获取所有用户
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		return fmt.Errorf("获取用户列表失败: %w", err)
	}

	if len(users) == 0 {
		return fmt.Errorf("没有找到可用的用户，无法生成行为数据")
	}

	g.LogInfo("开始为 %d 个用户生成累积行为数据，每个用户 %d 条记录", len(users), count)
	g.LogInfo("行为时间范围: %s 至 %s", startTime.Format("2006-01-02 15:04:05"), endTime.Format("2006-01-02 15:04:05"))

	// 验证时间范围
	if endTime.Before(startTime) {
		return fmt.Errorf("结束时间不能早于起始时间")
	}

	// 时间范围（以秒为单位）
	timeRange := int(endTime.Sub(startTime).Seconds())
	if timeRange <= 0 {
		timeRange = 1 // 如果时间范围太小，至少设置为1秒
	}

	// 统计生成的数据总量
	totalBehaviors := 0

	// 为每个用户生成行为数据
	for _, user := range users {
		// 为用户生成指定数量的行为数据
		for i := 0; i < count; i++ {
			// 在起始时间和结束时间之间随机选择一个时间点
			randomSeconds := g.RandomInt(0, timeRange)
			behaviorTime := startTime.Add(time.Duration(randomSeconds) * time.Second)

			// 随机生成行为类型：1-发送，2-接收
			behaviorType := g.RandomInt(1, 2)

			// 随机生成数据类型：1-文件，2-消息
			dataType := g.RandomInt(1, 2)

			// 随机生成数据大小
			var dataSize int64
			if dataType == 1 { // 文件
				dataSize = g.RandomInt64(1024, 5*1024*1024) // 1KB - 5MB
			} else { // 消息
				dataSize = g.RandomInt64(10, 1024) // 10B - 1KB
			}

			behavior := &models.UserBehavior{
				BehaviorID:   0, // 设置为0，MySQL会自动处理自增字段
				UserID:       user.UserID,
				BehaviorTime: behaviorTime,
				BehaviorType: behaviorType,
				DataType:     dataType,
				DataSize:     dataSize,
			}

			if err := db.Create(behavior).Error; err != nil {
				g.LogInfo("创建用户 %d 的累积行为数据失败: %v", user.UserID, err)
				continue
			}

			totalBehaviors++
		}
	}

	g.LogInfo("累积用户行为数据生成完成，成功创建 %d 条行为记录", totalBehaviors)
	return nil
}
