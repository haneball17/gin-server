package fileutil

import (
	"bytes"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// 创建测试目录
	testDir := "test_data"
	if err := os.MkdirAll(testDir, 0755); err != nil {
		panic(err)
	}

	// 运行测试
	code := m.Run()

	// 清理测试目录
	os.RemoveAll(testDir)

	os.Exit(code)
}

func TestEnsureDir(t *testing.T) {
	testCases := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "创建新目录",
			path:    "test_data/new_dir",
			wantErr: false,
		},
		{
			name:    "创建嵌套目录",
			path:    "test_data/nested/dir",
			wantErr: false,
		},
		{
			name:    "创建已存在的目录",
			path:    "test_data/existing_dir",
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.name == "创建已存在的目录" {
				if err := os.MkdirAll(tc.path, 0755); err != nil {
					t.Fatal(err)
				}
			}

			err := EnsureDir(tc.path)
			if (err != nil) != tc.wantErr {
				t.Errorf("EnsureDir() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr {
				if _, err := os.Stat(tc.path); os.IsNotExist(err) {
					t.Errorf("目录未创建: %s", tc.path)
				}
			}
		})
	}
}

func TestCopyFile(t *testing.T) {
	// 创建测试文件
	srcPath := "test_data/source.txt"
	content := []byte("测试文件内容")
	if err := os.WriteFile(srcPath, content, 0644); err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name    string
		src     string
		dst     string
		wantErr bool
	}{
		{
			name:    "复制文件到新位置",
			src:     srcPath,
			dst:     "test_data/dest.txt",
			wantErr: false,
		},
		{
			name:    "源文件不存在",
			src:     "test_data/nonexistent.txt",
			dst:     "test_data/dest2.txt",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := CopyFile(tc.src, tc.dst)
			if (err != nil) != tc.wantErr {
				t.Errorf("CopyFile() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr {
				// 验证文件内容
				destContent, err := os.ReadFile(tc.dst)
				if err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(content, destContent) {
					t.Error("复制的文件内容不匹配")
				}

				// 验证文件权限
				srcInfo, _ := os.Stat(tc.src)
				dstInfo, _ := os.Stat(tc.dst)
				if srcInfo.Mode() != dstInfo.Mode() {
					t.Error("复制的文件权限不匹配")
				}
			}
		})
	}
}

func TestBackupFile(t *testing.T) {
	// 创建测试文件
	srcPath := "test_data/original.txt"
	content := []byte("需要备份的文件内容")
	if err := os.WriteFile(srcPath, content, 0644); err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name    string
		src     string
		wantErr bool
	}{
		{
			name:    "备份存在的文件",
			src:     srcPath,
			wantErr: false,
		},
		{
			name:    "备份不存在的文件",
			src:     "test_data/nonexistent.txt",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			backupPath, err := BackupFile(tc.src)
			if (err != nil) != tc.wantErr {
				t.Errorf("BackupFile() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr {
				// 验证备份文件存在
				if _, err := os.Stat(backupPath); os.IsNotExist(err) {
					t.Error("备份文件未创建")
				}

				// 验证备份文件内容
				backupContent, err := os.ReadFile(backupPath)
				if err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(content, backupContent) {
					t.Error("备份文件内容不匹配")
				}
			}
		})
	}
}

func TestGetFileInfo(t *testing.T) {
	// 创建测试文件
	filePath := "test_data/info_test.txt"
	content := []byte("测试文件信息")
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		t.Fatal(err)
	}

	// 创建测试目录
	dirPath := "test_data/info_dir"
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name    string
		path    string
		wantDir bool
		wantErr bool
	}{
		{
			name:    "获取文件信息",
			path:    filePath,
			wantDir: false,
			wantErr: false,
		},
		{
			name:    "获取目录信息",
			path:    dirPath,
			wantDir: true,
			wantErr: false,
		},
		{
			name:    "获取不存在的文件信息",
			path:    "test_data/nonexistent.txt",
			wantDir: false,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			info, err := GetFileInfo(tc.path)
			if (err != nil) != tc.wantErr {
				t.Errorf("GetFileInfo() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr {
				if info.IsDir != tc.wantDir {
					t.Errorf("GetFileInfo() IsDir = %v, want %v", info.IsDir, tc.wantDir)
				}

				if info.Path != tc.path {
					t.Errorf("GetFileInfo() Path = %v, want %v", info.Path, tc.path)
				}

				if !tc.wantDir && info.Size != int64(len(content)) {
					t.Errorf("GetFileInfo() Size = %v, want %v", info.Size, len(content))
				}
			}
		})
	}
}

func TestIsFileExists(t *testing.T) {
	// 创建测试文件
	filePath := "test_data/exists_test.txt"
	if err := os.WriteFile(filePath, []byte("测试文件"), 0644); err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "文件存在",
			path: filePath,
			want: true,
		},
		{
			name: "文件不存在",
			path: "test_data/nonexistent.txt",
			want: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsFileExists(tc.path); got != tc.want {
				t.Errorf("IsFileExists() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestWriteAndReadFile(t *testing.T) {
	filePath := "test_data/write_read_test.txt"
	content := []byte("测试文件读写")

	// 测试写入文件
	t.Run("写入文件", func(t *testing.T) {
		err := WriteFile(filePath, content, 0644)
		if err != nil {
			t.Errorf("WriteFile() error = %v", err)
			return
		}

		// 验证文件存在
		if !IsFileExists(filePath) {
			t.Error("文件未创建")
		}
	})

	// 测试读取文件
	t.Run("读取文件", func(t *testing.T) {
		data, err := ReadFile(filePath)
		if err != nil {
			t.Errorf("ReadFile() error = %v", err)
			return
		}

		if !bytes.Equal(data, content) {
			t.Error("读取的文件内容不匹配")
		}
	})

	// 测试读取不存在的文件
	t.Run("读取不存在的文件", func(t *testing.T) {
		_, err := ReadFile("test_data/nonexistent.txt")
		if err == nil {
			t.Error("期望出现错误，但没有")
		}
	})
}
