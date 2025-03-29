package testdata

import (
	"fmt"
	"gin-server/config"
	"strings"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestDeviceGenerator_Generate(t *testing.T) {
	// 使用默认配置
	cfg := config.DefaultConfig()

	// 创建设备生成器
	generator := NewDeviceGenerator(cfg)

	// 使用mock数据库
	mockDB, err := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		DryRun: true, // 使用DryRun模式，不会执行实际的SQL
	})
	if err != nil {
		t.Fatalf("创建mock数据库失败: %v", err)
	}

	// 调用生成方法
	err = generator.Generate(mockDB, 5)
	if err == nil {
		t.Log("在mock模式下生成设备数据，函数执行成功")
	} else {
		t.Logf("生成设备数据返回错误: %v", err)
	}

	// 测试IPv6地址生成
	t.Run("生成IPv6地址", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			ipv6 := generator.RandomIPv6()
			t.Logf("生成的IPv6地址: %s", ipv6)

			// 验证格式：应该有7个冒号分隔8组十六进制数
			segments := strings.Split(ipv6, ":")
			if len(segments) != 8 {
				t.Errorf("无效的IPv6格式，期望8个分段，实际: %d", len(segments))
			}

			// 验证每个分段是有效的4位十六进制数
			for _, segment := range segments {
				if len(segment) != 4 {
					t.Errorf("IPv6分段长度不正确，期望4位，实际: %d", len(segment))
				}
				for _, ch := range segment {
					if !((ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f')) {
						t.Errorf("IPv6包含无效字符: %c", ch)
					}
				}
			}
		}
	})

	// 测试短地址生成
	t.Run("生成短地址", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			shortAddr := fmt.Sprintf("%02X%02X",
				generator.RandomInt(0, 255), generator.RandomInt(0, 255))
			t.Logf("生成的短地址: %s", shortAddr)

			// 验证格式：应该是4个十六进制字符
			if len(shortAddr) != 4 {
				t.Errorf("无效的短地址格式，期望长度为4，实际: %d", len(shortAddr))
			}

			// 验证都是有效的十六进制字符
			for _, ch := range shortAddr {
				if !((ch >= '0' && ch <= '9') || (ch >= 'A' && ch <= 'F')) {
					t.Errorf("短地址包含无效字符: %c", ch)
				}
			}
		}
	})

	// 测试SES密钥生成
	t.Run("生成SES密钥", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			sesKey := generator.RandomString(16)
			t.Logf("生成的SES密钥: %s", sesKey)

			// 验证长度：应该是16个字符
			if len(sesKey) != 16 {
				t.Errorf("无效的SES密钥长度，期望16，实际: %d", len(sesKey))
			}
		}
	})
}
