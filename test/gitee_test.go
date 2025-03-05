package test

import (
	"encoding/json"
	"gin-server/config"
	"gin-server/configmanager/common/transfer"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGiteeTransporter(t *testing.T) {
	// 初始化配置
	config.InitConfig()
	cfg := config.GetConfig()

	if cfg.Gitee == nil {
		t.Fatal("Gitee配置为空")
	}

	// 创建测试文件
	testData := map[string]interface{}{
		"test_data": map[string]string{
			"name":        "测试文件",
			"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
			"description": "这是一个用于测试Gitee上传功能的JSON文件",
		},
	}

	// 确保test目录存在
	if err := os.MkdirAll("test", 0755); err != nil {
		t.Fatalf("创建test目录失败: %v", err)
	}

	// 创建测试文件
	localPath := filepath.Join("test", "test_data.json")
	jsonData, err := json.MarshalIndent(testData, "", "    ")
	if err != nil {
		t.Fatalf("JSON编码失败: %v", err)
	}

	if err := os.WriteFile(localPath, jsonData, 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}
	defer os.Remove(localPath) // 测试完成后清理文件

	// 创建Gitee传输器
	transporter := transfer.NewGiteeTransporter(cfg.Gitee)

	// 测试文件路径
	remotePath := "log/test/test_data.json"

	// 测试上传文件
	t.Run("上传JSON文件", func(t *testing.T) {
		t.Logf("尝试上传文件到: %s", remotePath)
		err := transporter.Upload(localPath, remotePath)
		if err != nil {
			t.Fatalf("上传文件失败: %v", err)
		}
		t.Log("文件上传成功")

		// 验证文件是否上传成功
		files, err := transporter.List("log/test")
		if err != nil {
			t.Fatalf("列出目录内容失败: %v", err)
		}

		found := false
		for _, file := range files {
			t.Logf("发现文件: %s", file.Path)
			if file.Path == remotePath {
				found = true
				break
			}
		}

		if !found {
			t.Fatal("上传的文件未在目录列表中找到")
		}
		t.Log("已在目录列表中找到上传的文件")
	})

	// 测试删除文件
	t.Run("删除JSON文件", func(t *testing.T) {
		t.Logf("尝试删除文件: %s", remotePath)
		err := transporter.Delete(remotePath)
		if err != nil {
			t.Fatalf("删除文件失败: %v", err)
		}
		t.Log("文件删除成功")

		// 验证文件是否已被删除
		files, err := transporter.List("log/test")
		if err != nil {
			t.Fatalf("列出目录内容失败: %v", err)
		}

		for _, file := range files {
			if file.Path == remotePath {
				t.Fatal("文件仍然存在于目录列表中")
			}
		}
		t.Log("已确认文件被成功删除")
	})
}
