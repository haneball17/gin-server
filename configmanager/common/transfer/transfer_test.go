package transfer

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"gin-server/config"
)

func TestMain(m *testing.M) {
	// 创建测试目录
	err := os.MkdirAll("test_data", 0755)
	if err != nil {
		panic(err)
	}

	// 运行测试
	code := m.Run()

	// 清理测试目录
	os.RemoveAll("test_data")

	os.Exit(code)
}

func TestNewFileTransporter(t *testing.T) {
	tests := []struct {
		name    string
		typ     TransporterType
		cfg     *config.Config
		wantErr bool
		skip    bool
	}{
		{
			name: "创建Gitee传输器-成功",
			typ:  TransporterTypeGitee,
			cfg: &config.Config{
				Gitee: &config.GiteeConfig{
					AccessToken: "test_token",
					Owner:       "test_owner",
					Repo:        "test_repo",
					Branch:      "master",
				},
			},
			wantErr: false,
		},
		{
			name: "创建FTP传输器-成功",
			typ:  TransporterTypeFTP,
			cfg: &config.Config{
				FTP: &config.FTPConfig{
					Host:     os.Getenv("FTP_TEST_HOST"),
					Port:     21,
					Username: os.Getenv("FTP_TEST_USERNAME"),
					Password: os.Getenv("FTP_TEST_PASSWORD"),
				},
			},
			wantErr: false,
			skip:    os.Getenv("FTP_TEST_HOST") == "" || os.Getenv("FTP_TEST_USERNAME") == "" || os.Getenv("FTP_TEST_PASSWORD") == "",
		},
		{
			name:    "创建Gitee传输器-配置为空",
			typ:     TransporterTypeGitee,
			cfg:     &config.Config{},
			wantErr: true,
		},
		{
			name:    "创建FTP传输器-配置为空",
			typ:     TransporterTypeFTP,
			cfg:     &config.Config{},
			wantErr: true,
		},
		{
			name:    "不支持的传输器类型",
			typ:     "unknown",
			cfg:     &config.Config{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip("未设置测试环境变量")
			}
			_, err := NewFileTransporter(tt.typ, tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFileTransporter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGiteeTransporter(t *testing.T) {
	cfg := &config.GiteeConfig{
		AccessToken: os.Getenv("GITEE_TEST_TOKEN"),
		Owner:       os.Getenv("GITEE_TEST_OWNER"),
		Repo:        os.Getenv("GITEE_TEST_REPO"),
		Branch:      "master",
	}

	// 如果没有设置测试环境变量，跳过测试
	if cfg.AccessToken == "" || cfg.Owner == "" || cfg.Repo == "" {
		t.Skip("未设置Gitee测试环境变量")
	}

	transporter := NewGiteeTransporter(cfg)

	// 准备测试文件
	testFile := filepath.Join("test_data", "test.txt")
	testContent := []byte("test content")
	err := os.WriteFile(testFile, testContent, 0644)
	if err != nil {
		t.Fatal(err)
	}

	// 测试上传文件
	t.Run("上传文件", func(t *testing.T) {
		err := transporter.Upload(testFile, "test.txt")
		if err != nil {
			t.Fatal(err)
		}
	})

	// 等待Gitee API更新
	time.Sleep(time.Second)

	// 测试获取文件修改时间
	t.Run("获取文件修改时间", func(t *testing.T) {
		_, err := transporter.LastModified("test.txt")
		if err != nil {
			t.Fatal(err)
		}
	})

	// 测试列出目录内容
	t.Run("列出目录内容", func(t *testing.T) {
		files, err := transporter.List("")
		if err != nil {
			t.Fatal(err)
		}
		if len(files) == 0 {
			t.Error("目录为空")
		}
	})

	// 测试下载文件
	t.Run("下载文件", func(t *testing.T) {
		downloadFile := filepath.Join("test_data", "download.txt")
		err := transporter.Download("test.txt", downloadFile)
		if err != nil {
			t.Fatal(err)
		}

		content, err := os.ReadFile(downloadFile)
		if err != nil {
			t.Fatal(err)
		}

		if string(content) != string(testContent) {
			t.Errorf("文件内容不匹配: got = %v, want = %v", string(content), string(testContent))
		}
	})

	// 测试删除文件
	t.Run("删除文件", func(t *testing.T) {
		err := transporter.Delete("test.txt")
		if err != nil {
			t.Fatal(err)
		}
	})

	// 测试关闭传输器
	t.Run("关闭传输器", func(t *testing.T) {
		err := transporter.Close()
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestFTPTransporter(t *testing.T) {
	cfg := &config.FTPConfig{
		Host:     os.Getenv("FTP_TEST_HOST"),
		Port:     21,
		Username: os.Getenv("FTP_TEST_USERNAME"),
		Password: os.Getenv("FTP_TEST_PASSWORD"),
	}

	// 如果没有设置测试环境变量，跳过测试
	if cfg.Host == "" || cfg.Username == "" || cfg.Password == "" {
		t.Skip("未设置FTP测试环境变量")
	}

	transporter, err := NewFTPTransporter(cfg)
	if err != nil {
		t.Fatal(err)
	}

	// 准备测试文件
	testFile := filepath.Join("test_data", "test.txt")
	testContent := []byte("test content")
	err = os.WriteFile(testFile, testContent, 0644)
	if err != nil {
		t.Fatal(err)
	}

	// 测试上传文件
	t.Run("上传文件", func(t *testing.T) {
		err := transporter.Upload(testFile, "test.txt")
		if err != nil {
			t.Fatal(err)
		}
	})

	// 测试获取文件修改时间
	t.Run("获取文件修改时间", func(t *testing.T) {
		_, err := transporter.LastModified("test.txt")
		if err != nil {
			t.Fatal(err)
		}
	})

	// 测试列出目录内容
	t.Run("列出目录内容", func(t *testing.T) {
		files, err := transporter.List("")
		if err != nil {
			t.Fatal(err)
		}
		if len(files) == 0 {
			t.Error("目录为空")
		}
	})

	// 测试下载文件
	t.Run("下载文件", func(t *testing.T) {
		downloadFile := filepath.Join("test_data", "download.txt")
		err := transporter.Download("test.txt", downloadFile)
		if err != nil {
			t.Fatal(err)
		}

		content, err := os.ReadFile(downloadFile)
		if err != nil {
			t.Fatal(err)
		}

		if string(content) != string(testContent) {
			t.Errorf("文件内容不匹配: got = %v, want = %v", string(content), string(testContent))
		}
	})

	// 测试删除文件
	t.Run("删除文件", func(t *testing.T) {
		err := transporter.Delete("test.txt")
		if err != nil {
			t.Fatal(err)
		}
	})

	// 测试关闭传输器
	t.Run("关闭传输器", func(t *testing.T) {
		err := transporter.Close()
		if err != nil {
			t.Fatal(err)
		}
	})
}
