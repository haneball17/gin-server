package test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"gin-server/config"
	"gin-server/configmanager/log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestLogGeneration(t *testing.T) {
	// 初始化数据库连接
	dsn := "root:1234@tcp(127.0.0.1:3306)/gin_server?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("连接数据库失败: %v", err)
	}

	// 创建配置
	cfg := &config.Config{
		DebugLevel: "true",
		ConfigManager: config.ConfigManagerConfig{
			LogManager: config.LogManagerConfig{
				GenerateInterval: 5,
				EnableEncryption: false,
			},
		},
	}

	// 创建日志管理器
	logManager := log.NewLogManager(cfg, db)

	// 生成日志
	err = logManager.GenerateLog()
	if err != nil {
		t.Fatalf("生成日志失败: %v", err)
	}

	// 等待日志生成完成
	time.Sleep(2 * time.Second)

	// 列出logs目录下的所有子目录
	entries, err := os.ReadDir("logs")
	if err != nil {
		t.Fatalf("读取logs目录失败: %v", err)
	}

	// 找到最新的日志目录
	var latestDir string
	var latestTime time.Time
	for _, entry := range entries {
		if entry.IsDir() {
			// 解析目录名中的时间戳
			dirTime, err := time.ParseInLocation("20060102150405", entry.Name(), time.Local)
			if err != nil {
				continue
			}
			if latestDir == "" || dirTime.After(latestTime) {
				latestDir = entry.Name()
				latestTime = dirTime
			}
		}
	}

	if latestDir == "" {
		t.Fatal("未找到日志目录")
	}

	// 检查日志文件是否存在
	logFile := filepath.Join("logs", latestDir, "log.json")
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Errorf("日志文件未生成: %s", logFile)
	} else {
		// 读取并打印日志文件内容
		data, err := os.ReadFile(logFile)
		if err != nil {
			t.Errorf("读取日志文件失败: %v", err)
		} else {
			t.Logf("日志文件内容:\n%s", string(data))
		}
	}

	t.Log("日志生成成功")
}
