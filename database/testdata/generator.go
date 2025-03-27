package testdata

import (
	"fmt"
	"gin-server/config"
	"log"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

// Generator 是测试数据生成器接口
type Generator interface {
	// Generate 生成测试数据
	// db 是数据库连接
	// count 是要生成的数据量
	Generate(db *gorm.DB, count int) error
}

// BaseGenerator 是所有生成器的基础实现
type BaseGenerator struct {
	Cfg    *config.Config
	random *rand.Rand
}

// NewBaseGenerator 创建一个基础生成器
func NewBaseGenerator(cfg *config.Config) *BaseGenerator {
	source := rand.NewSource(time.Now().UnixNano())
	return &BaseGenerator{
		Cfg:    cfg,
		random: rand.New(source),
	}
}

// RandomString 生成指定长度的随机字符串
func (g *BaseGenerator) RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[g.random.Intn(len(charset))]
	}
	return string(b)
}

// RandomIP 生成随机IP地址
func (g *BaseGenerator) RandomIP() string {
	return fmt.Sprintf("%d.%d.%d.%d",
		g.random.Intn(256),
		g.random.Intn(256),
		g.random.Intn(256),
		g.random.Intn(256))
}

// RandomPort 生成随机端口号
func (g *BaseGenerator) RandomPort() int {
	return g.random.Intn(1000) + 8000 // 生成8000-8999之间的端口
}

// RandomEmail 生成随机邮箱
func (g *BaseGenerator) RandomEmail() string {
	domains := []string{"example.com", "test.org", "mail.cn", "company.net"}
	return g.RandomString(8) + "@" + domains[g.random.Intn(len(domains))]
}

// RandomInt 生成范围内的随机整数
func (g *BaseGenerator) RandomInt(min, max int) int {
	return g.random.Intn(max-min+1) + min
}

// RandomInt64 生成范围内的随机64位整数
func (g *BaseGenerator) RandomInt64(min, max int64) int64 {
	return min + g.random.Int63n(max-min+1)
}

// RandomBool 生成随机布尔值
func (g *BaseGenerator) RandomBool() bool {
	return g.random.Intn(2) == 1
}

// LogInfo 记录信息日志
func (g *BaseGenerator) LogInfo(format string, args ...interface{}) {
	if g.Cfg.DebugLevel == "true" {
		log.Printf(format, args...)
	}
}

// IsDataExists 检查指定表中是否已存在数据
func IsDataExists(db *gorm.DB, model interface{}) (bool, error) {
	var count int64
	err := db.Model(model).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
