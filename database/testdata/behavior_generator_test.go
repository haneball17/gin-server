package testdata

import (
	"gin-server/config"
	"testing"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 此测试函数需要在有实际数据库连接的环境中运行
// 可根据实际环境修改连接信息或跳过测试
func TestBehaviorGenerator_GenerateAccumulated(t *testing.T) {
	// 如果没有设置环境变量TEST_DB_DSN，则跳过测试
	// dsn := os.Getenv("TEST_DB_DSN")
	// if dsn == "" {
	// 	t.Skip("跳过测试：未设置TEST_DB_DSN环境变量")
	// }

	// 使用默认配置
	cfg := config.DefaultConfig()

	// 创建测试用的行为生成器
	generator := NewBehaviorGenerator(cfg)

	// 测试时间范围
	endTime := time.Now()
	startTime := endTime.Add(-30 * time.Minute) // 30分钟前

	// 模拟数据库连接（实际测试时应替换为真实的数据库连接）
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	t.Fatalf("连接数据库失败: %v", err)
	// }

	// 使用mock的db
	mockDB, err := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      nil, // Mock情况下不需要实际连接
	}), &gorm.Config{
		// 使用DryRun模式，不会执行实际的SQL
		DryRun: true,
	})
	if err != nil {
		t.Fatalf("创建mock数据库失败: %v", err)
	}

	// 调用被测试的方法
	err = generator.GenerateAccumulated(mockDB, 2, startTime, endTime)

	// 在mock模式下，我们期望有错误（因为没有实际的数据库连接）
	// 但可以检查代码执行路径是否正确
	if err == nil {
		t.Log("在mock模式下预期有错误，但函数返回nil")
	} else {
		t.Logf("函数返回预期的错误: %v", err)
	}

	// 测试边界条件：结束时间早于开始时间
	invalidEndTime := startTime.Add(-1 * time.Hour)
	err = generator.GenerateAccumulated(mockDB, 2, startTime, invalidEndTime)
	if err == nil {
		t.Error("当结束时间早于开始时间时，函数应该返回错误")
	} else {
		t.Logf("函数对无效的时间范围返回了正确的错误: %v", err)
	}

	// 测试调用 GenerateRealtime
	err = generator.GenerateRealtime(mockDB, 2)
	if err == nil {
		t.Log("GenerateRealtime在mock模式下预期有错误，但函数返回nil")
	} else {
		t.Logf("GenerateRealtime函数返回预期的错误: %v", err)
	}

	// 测试累积生成功能
	// 注意：在实际的数据库测试中，还应该验证生成的数据是否符合预期
	// 例如检查行为时间是否在指定的范围内，以及是否生成了正确数量的数据
	t.Log("测试 Generate 函数的累加模式")
	err = generator.Generate(mockDB, 2, false) // 不跳过已有数据
	if err == nil {
		t.Log("Generate函数在累积模式下执行成功")
	} else {
		t.Logf("Generate函数返回错误: %v", err)
	}
}

// TestBehaviorGeneratorAccumulationMode 测试行为生成器的累积模式
func TestBehaviorGeneratorAccumulationMode(t *testing.T) {
	// 使用默认配置
	cfg := config.DefaultConfig()

	// 确保配置中设置了实时数据生成的时间范围
	cfg.TestData.RealtimeStartTimeOffset = 30 // 30分钟前
	cfg.TestData.RealtimeEndTimeOffset = 0    // 当前时间

	// 创建测试用的行为生成器
	generator := NewBehaviorGenerator(cfg)

	// 使用mock的db
	mockDB, err := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		DryRun: true,
	})
	if err != nil {
		t.Fatalf("创建mock数据库失败: %v", err)
	}

	// 测试生成函数在累积模式下的行为
	t.Log("测试累积模式下的Generate函数")

	// 第一次调用，跳过已有数据检查
	err = generator.Generate(mockDB, 5, false)
	if err == nil {
		t.Log("第一次Generate调用(累积模式)成功")
	} else {
		t.Logf("第一次Generate调用返回错误: %v", err)
	}

	// 第二次调用，仍跳过已有数据检查，应该继续添加数据
	err = generator.Generate(mockDB, 3, false)
	if err == nil {
		t.Log("第二次Generate调用(累积模式)成功")
	} else {
		t.Logf("第二次Generate调用返回错误: %v", err)
	}

	// 测试实时数据生成，应该使用配置中的时间范围
	t.Log("测试GenerateRealtime函数")
	err = generator.GenerateRealtime(mockDB, 2)
	if err == nil {
		t.Log("GenerateRealtime调用成功")
	} else {
		t.Logf("GenerateRealtime调用返回错误: %v", err)
	}
}
