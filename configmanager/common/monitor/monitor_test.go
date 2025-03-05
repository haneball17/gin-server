package monitor

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"gin-server/configmanager/common/fileutil"
	"gin-server/configmanager/common/transfer"
)

// MockTransporter 模拟文件传输器
type MockTransporter struct {
	t            *testing.T
	files        map[string][]byte
	modTimes     map[string]time.Time
	lastModError map[string]error
}

func NewMockTransporter() *MockTransporter {
	return &MockTransporter{
		files:        make(map[string][]byte),
		modTimes:     make(map[string]time.Time),
		lastModError: make(map[string]error),
	}
}

func (m *MockTransporter) SetT(t *testing.T) {
	m.t = t
}

func (m *MockTransporter) Upload(localPath, remotePath string) error {
	content, err := os.ReadFile(localPath)
	if err != nil {
		return err
	}

	m.files[remotePath] = content
	return nil
}

func (m *MockTransporter) Download(remotePath, localPath string) error {
	m.t.Logf("下载文件: %s -> %s\n", remotePath, localPath)
	content, ok := m.files[remotePath]
	if !ok {
		return fmt.Errorf("file not found")
	}

	if err := fileutil.EnsureDir(filepath.Dir(localPath)); err != nil {
		return err
	}

	if err := os.WriteFile(localPath, content, 0644); err != nil {
		return err
	}

	m.t.Logf("文件下载成功: %s\n", localPath)
	return nil
}

func (m *MockTransporter) List(dir string) ([]transfer.FileInfo, error) {
	m.t.Logf("列出目录内容: %s\n", dir)
	var result []transfer.FileInfo

	// 统一使用正斜杠作为路径分隔符
	dir = filepath.ToSlash(dir)

	// 遍历所有文件，找出指定目录下的文件
	for path, content := range m.files {
		// 统一使用正斜杠作为路径分隔符
		path = filepath.ToSlash(path)
		if strings.HasPrefix(path, dir) {
			modTime := m.modTimes[path]
			if modTime.IsZero() {
				modTime = time.Now()
			}
			result = append(result, transfer.FileInfo{
				Path:    path,
				Size:    int64(len(content)),
				ModTime: modTime,
			})
		}
	}

	m.t.Logf("找到 %d 个文件\n", len(result))
	return result, nil
}

func (m *MockTransporter) Delete(remotePath string) error {
	return nil
}

func (m *MockTransporter) LastModified(path string) (time.Time, error) {
	m.t.Logf("获取文件修改时间: %s\n", path)
	if err, ok := m.lastModError[path]; ok && err != nil {
		m.t.Logf("获取文件修改时间失败: %s - %v\n", path, err)
		return time.Time{}, err
	}

	if t, ok := m.modTimes[path]; ok {
		m.t.Logf("返回文件修改时间: %s - %v\n", path, t)
		return t, nil
	}

	return time.Time{}, fmt.Errorf("file not found")
}

func (m *MockTransporter) Close() error {
	return nil
}

func (m *MockTransporter) SetLastModified(path string, t time.Time) {
	m.t.Logf("设置文件修改时间: %s - %v\n", path, t)
	m.modTimes[path] = t
}

func (m *MockTransporter) SetLastModifiedError(path string, err error) {
	m.t.Logf("设置获取修改时间错误: %v\n", err)
	m.lastModError[path] = err
}

func (m *MockTransporter) SetFileContent(path string, content []byte) {
	m.t.Logf("设置文件内容: %s - %d bytes\n", path, len(content))
	m.files[path] = content
}

func TestRemoteMonitor(t *testing.T) {
	// 创建模拟传输器
	mock := NewMockTransporter()
	mock.SetT(t)

	// 设置初始文件内容
	mock.SetFileContent("/test/config.yaml", []byte("initial content"))

	// 创建监控器
	monitor := NewRemoteMonitor(mock, time.Millisecond*100)

	// 测试启动监控器
	t.Run("启动监控器", func(t *testing.T) {
		mock.SetT(t)
		if err := monitor.Start(); err != nil {
			t.Errorf("启动监控器失败: %v", err)
		}
		// 等待监控协程启动
		time.Sleep(time.Millisecond * 200)
	})

	// 测试添加监控路径
	path := "/test/config.yaml"
	t.Run("添加监控路径", func(t *testing.T) {
		mock.SetT(t)
		// 设置初始修改时间
		initialTime := time.Now()
		mock.SetLastModified(path, initialTime)
		if err := monitor.AddWatch(path); err != nil {
			t.Errorf("添加监控路径失败: %v", err)
		}
		// 等待路径被添加
		time.Sleep(time.Millisecond * 100)
	})

	// 测试文件修改事件
	t.Run("文件修改事件", func(t *testing.T) {
		mock.SetT(t)
		// 设置新的文件内容和修改时间
		mock.SetFileContent(path, []byte("updated content"))
		newTime := time.Now().Add(time.Hour)
		t.Logf("设置新的修改时间: %v", newTime)
		mock.SetLastModified(path, newTime)

		// 等待事件
		t.Log("等待文件修改事件...")
		select {
		case event := <-monitor.Events():
			t.Logf("收到事件: 类型=%s, 路径=%s, 时间=%v", event.Type, event.Path, event.FileInfo.ModTime)
			if event.Type != EventModify {
				t.Errorf("期望事件类型为 %s，实际为 %s", EventModify, event.Type)
			}
			if event.Path != path {
				t.Errorf("期望路径为 %s，实际为 %s", path, event.Path)
			}
			if !event.FileInfo.ModTime.Equal(newTime) {
				t.Errorf("期望修改时间为 %v，实际为 %v", newTime, event.FileInfo.ModTime)
			}
		case <-time.After(time.Second * 2):
			t.Error("等待事件超时")
		}
	})

	// 测试文件删除事件
	t.Run("文件删除事件", func(t *testing.T) {
		mock.SetT(t)
		// 等待上一个事件处理完成
		time.Sleep(time.Millisecond * 200)
		// 设置错误以模拟文件不存在
		mock.SetLastModifiedError(path, errors.New("file not found"))

		// 等待事件
		select {
		case event := <-monitor.Events():
			if event.Type != EventDelete {
				t.Errorf("期望事件类型为 %s，实际为 %s", EventDelete, event.Type)
			}
			if event.Path != path {
				t.Errorf("期望路径为 %s，实际为 %s", path, event.Path)
			}
		case <-time.After(time.Second * 2):
			t.Error("等待事件超时")
		}
	})

	// 测试移除监控路径
	t.Run("移除监控路径", func(t *testing.T) {
		mock.SetT(t)
		if err := monitor.RemoveWatch(path); err != nil {
			t.Errorf("移除监控路径失败: %v", err)
		}
	})

	// 测试停止监控器
	t.Run("停止监控器", func(t *testing.T) {
		mock.SetT(t)
		if err := monitor.Stop(); err != nil {
			t.Errorf("停止监控器失败: %v", err)
		}
	})
}

func TestRemoteMonitorErrors(t *testing.T) {
	// 创建模拟传输器
	mock := NewMockTransporter()
	mock.SetT(t)

	// 设置初始文件内容和路径
	path := "/test/config.yaml"
	mock.SetFileContent(path, []byte("initial content"))
	initialTime := time.Now()
	mock.SetLastModified(path, initialTime)

	// 创建监控器
	monitor := NewRemoteMonitor(mock, time.Millisecond*100)

	// 测试下载错误处理
	t.Run("下载错误处理", func(t *testing.T) {
		mock.SetT(t)
		// 启动监控器
		if err := monitor.Start(); err != nil {
			t.Errorf("启动监控器失败: %v", err)
		}

		// 添加监控路径
		if err := monitor.AddWatch(path); err != nil {
			t.Errorf("添加监控路径失败: %v", err)
		}

		// 设置下载错误
		mock.SetLastModifiedError(path, errors.New("download failed"))
		// 触发文件修改
		mock.SetLastModified(path, time.Now().Add(time.Hour))

		// 等待事件
		select {
		case event := <-monitor.Events():
			if event.Type != EventDelete {
				t.Errorf("期望事件类型为 %s，实际为 %s", EventDelete, event.Type)
			}
			if event.Error == nil {
				t.Error("期望有错误，实际为 nil")
			}
		case <-time.After(time.Second * 2):
			t.Error("等待事件超时")
		}

		// 停止监控器
		if err := monitor.Stop(); err != nil {
			t.Errorf("停止监控器失败: %v", err)
		}
	})
}

func TestExtractTimeFromPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		wantTime time.Time
		wantErr  bool
	}{
		{
			name:     "正常路径 - 策略文件",
			path:     "/policies/20240305150000/policy.json",
			wantTime: time.Date(2024, 3, 5, 15, 0, 0, 0, time.Local),
			wantErr:  false,
		},
		{
			name:     "正常路径 - 日志文件",
			path:     "/logs/20240305150000/app.log",
			wantTime: time.Date(2024, 3, 5, 15, 0, 0, 0, time.Local),
			wantErr:  false,
		},
		{
			name:     "多层路径",
			path:     "/data/backup/20240305150000/config/settings.json",
			wantTime: time.Date(2024, 3, 5, 15, 0, 0, 0, time.Local),
			wantErr:  false,
		},
		{
			name:    "无时间信息",
			path:    "/config/settings.json",
			wantErr: true,
		},
		{
			name:    "错误的时间格式",
			path:    "/logs/2024-03-05/app.log",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTime, err := ExtractTimeFromPath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractTimeFromPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !gotTime.Equal(tt.wantTime) {
				t.Errorf("ExtractTimeFromPath() = %v, want %v", gotTime, tt.wantTime)
			}
		})
	}
}

func TestRemoteMonitorWithPathTime(t *testing.T) {
	// 创建模拟传输器
	mock := NewMockTransporter()
	mock.SetT(t)

	// 设置初始文件内容和路径
	path := "/policies/20240305150000/policy.json"
	mock.SetFileContent(path, []byte("initial content"))
	initialTime := time.Now()
	mock.SetLastModified(path, initialTime)

	// 创建监控器
	monitor := NewRemoteMonitor(mock, time.Millisecond*100)

	// 测试启动监控器
	t.Run("启动监控器", func(t *testing.T) {
		mock.SetT(t)
		if err := monitor.Start(); err != nil {
			t.Errorf("启动监控器失败: %v", err)
		}
		// 等待监控协程启动
		time.Sleep(time.Millisecond * 200)
	})

	// 测试添加监控路径
	t.Run("添加监控路径", func(t *testing.T) {
		mock.SetT(t)
		if err := monitor.AddWatch(path); err != nil {
			t.Errorf("添加监控路径失败: %v", err)
		}
		// 等待路径被添加
		time.Sleep(time.Millisecond * 100)
	})

	// 测试文件修改事件（通过更新路径时间）
	t.Run("文件修改事件 - 路径时间更新", func(t *testing.T) {
		mock.SetT(t)
		// 设置新的文件内容和路径（时间增加1小时）
		newPath := "/policies/20240305160000/policy.json"
		mock.SetFileContent(newPath, []byte("updated content"))
		mock.SetLastModified(newPath, initialTime) // 使用相同的修改时间，确保只通过路径时间检测变化

		// 移除旧路径，添加新路径
		monitor.RemoveWatch(path)
		if err := monitor.AddWatch(newPath); err != nil {
			t.Errorf("添加新监控路径失败: %v", err)
		}

		// 等待监控器检测到新路径
		time.Sleep(time.Millisecond * 300)

		// 等待事件
		t.Log("等待文件修改事件...")
		select {
		case event := <-monitor.Events():
			t.Logf("收到事件: 类型=%s, 路径=%s, 时间=%v, 路径时间=%v",
				event.Type, event.Path, event.FileInfo.ModTime, event.FileInfo.PathTime)
			if event.Type != EventModify {
				t.Errorf("期望事件类型为 %s，实际为 %s", EventModify, event.Type)
			}
			if event.Path != newPath {
				t.Errorf("期望路径为 %s，实际为 %s", newPath, event.Path)
			}
			// 验证路径时间
			pathTime, err := ExtractTimeFromPath(event.Path)
			if err != nil {
				t.Errorf("从路径提取时间失败: %v", err)
			}
			expectedTime := time.Date(2024, 3, 5, 16, 0, 0, 0, time.Local)
			if !pathTime.Equal(expectedTime) {
				t.Errorf("期望路径时间为 %v，实际为 %v", expectedTime, pathTime)
			}
		case <-time.After(time.Second * 2):
			t.Error("等待事件超时")
		}
	})

	// 测试停止监控器
	t.Run("停止监控器", func(t *testing.T) {
		mock.SetT(t)
		if err := monitor.Stop(); err != nil {
			t.Errorf("停止监控器失败: %v", err)
		}
	})
}
