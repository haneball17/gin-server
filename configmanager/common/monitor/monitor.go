package monitor

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"gin-server/config"
	"gin-server/configmanager/common/fileutil"
	"gin-server/configmanager/common/transfer"
)

// EventType 文件事件类型
type EventType string

const (
	// EventCreate 文件创建事件
	EventCreate EventType = "create"
	// EventModify 文件修改事件
	EventModify EventType = "modify"
	// EventDelete 文件删除事件
	EventDelete EventType = "delete"

	// 时间格式常量
	TimeFormat = "20060102150405" // YYYYMMDDHHmmss
)

// Event 文件事件
type Event struct {
	Type      EventType // 事件类型
	Path      string    // 文件路径
	FileInfo  *FileInfo // 文件信息
	Timestamp time.Time // 事件时间戳
	Error     error     // 事件相关错误
}

// FileInfo 文件信息
type FileInfo struct {
	Path         string    // 文件路径
	Hash         string    // 文件哈希
	Size         int64     // 文件大小
	ModTime      time.Time // 修改时间
	LastChecksum string    // 上次校验和
	PathTime     time.Time // 从路径中提取的时间信息
}

// Monitor 文件监控器接口
type Monitor interface {
	// Start 启动监控
	Start() error

	// Stop 停止监控
	Stop() error

	// AddWatch 添加监控路径
	AddWatch(path string) error

	// RemoveWatch 移除监控路径
	RemoveWatch(path string) error

	// Events 返回事件通道
	Events() <-chan Event
}

// RemoteMonitor 远程文件监控器
type RemoteMonitor struct {
	transporter transfer.FileTransporter // 文件传输器
	interval    time.Duration            // 轮询间隔
	paths       map[string]*FileInfo     // 监控的文件路径
	events      chan Event               // 事件通道
	stopChan    chan struct{}            // 停止信号通道
	wg          sync.WaitGroup           // 等待组
	mu          sync.RWMutex             // 读写锁
}

// normalizePath 统一路径分隔符
func normalizePath(path string) string {
	// 统一使用正斜杠作为路径分隔符
	return filepath.ToSlash(path)
}

// ExtractTimeFromPath 从文件路径中提取时间信息
func ExtractTimeFromPath(path string) (time.Time, error) {
	// 统一路径分隔符
	path = normalizePath(path)

	// 使用正则表达式匹配路径中的时间信息
	// 匹配格式：YYYYMMDDHHmmss
	re := regexp.MustCompile(`\d{14}`)
	match := re.FindString(path)
	if match == "" {
		return time.Time{}, fmt.Errorf("路径中未找到时间信息: %s", path)
	}

	// 解析时间字符串
	t, err := time.ParseInLocation(TimeFormat, match, time.Local)
	if err != nil {
		return time.Time{}, fmt.Errorf("解析时间信息失败: %w", err)
	}

	return t, nil
}

// NewRemoteMonitor 创建远程文件监控器
func NewRemoteMonitor(transporter transfer.FileTransporter, interval time.Duration) *RemoteMonitor {
	return &RemoteMonitor{
		transporter: transporter,
		interval:    interval,
		paths:       make(map[string]*FileInfo),
		events:      make(chan Event, 100), // 缓冲区大小为100
		stopChan:    make(chan struct{}),
	}
}

// Start 启动监控
func (m *RemoteMonitor) Start() error {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Println("启动远程文件监控器")
	}

	m.wg.Add(1)
	go m.monitor()
	return nil
}

// Stop 停止监控
func (m *RemoteMonitor) Stop() error {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Println("停止远程文件监控器")
	}

	close(m.stopChan)
	m.wg.Wait()
	close(m.events)
	return nil
}

// AddWatch 添加监控路径
func (m *RemoteMonitor) AddWatch(path string) error {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Printf("添加监控路径: %s\n", path)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// 获取文件信息
	modTime, err := m.transporter.LastModified(path)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("获取文件修改时间失败: %s - %v\n", path, err)
		}
		return err
	}

	// 从路径中提取时间信息
	pathTime, err := ExtractTimeFromPath(path)
	if err != nil && cfg.DebugLevel == "true" {
		log.Printf("从路径提取时间信息失败: %s - %v\n", path, err)
	}

	// 下载文件到临时目录以获取完整信息
	tmpDir := filepath.Join(os.TempDir(), "configmanager")
	if err := fileutil.EnsureDir(tmpDir); err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("创建临时目录失败: %v\n", err)
		}
		return err
	}

	tmpFile := filepath.Join(tmpDir, filepath.Base(path))
	if err := m.transporter.Download(path, tmpFile); err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("下载文件失败: %s - %v\n", path, err)
		}
		return err
	}

	// 获取文件信息
	fileInfo, err := fileutil.GetFileInfo(tmpFile)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("获取文件信息失败: %s - %v\n", tmpFile, err)
		}
		os.Remove(tmpFile)
		return err
	}

	// 清理临时文件
	os.Remove(tmpFile)

	// 保存文件信息
	m.paths[path] = &FileInfo{
		Path:     path,
		Size:     fileInfo.Size,
		ModTime:  modTime,
		PathTime: pathTime,
	}

	if cfg.DebugLevel == "true" {
		log.Printf("成功添加监控路径: %s, 修改时间: %v, 路径时间: %v, 大小: %d\n",
			path, modTime, pathTime, fileInfo.Size)
	}

	return nil
}

// RemoveWatch 移除监控路径
func (m *RemoteMonitor) RemoveWatch(path string) error {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Printf("移除监控路径: %s\n", path)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.paths, path)
	return nil
}

// Events 返回事件通道
func (m *RemoteMonitor) Events() <-chan Event {
	return m.events
}

// monitor 监控协程
func (m *RemoteMonitor) monitor() {
	defer m.wg.Done()

	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for {
		select {
		case <-m.stopChan:
			return
		case <-ticker.C:
			m.checkFiles()
		}
	}
}

// checkFiles 检查文件变化
func (m *RemoteMonitor) checkFiles() {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Println("开始检查文件变化...")
	}

	m.mu.RLock()
	paths := make([]string, 0, len(m.paths))
	fileInfos := make(map[string]*FileInfo)
	for path, info := range m.paths {
		paths = append(paths, path)
		fileInfos[path] = &FileInfo{
			Path:     info.Path,
			ModTime:  info.ModTime,
			Size:     info.Size,
			PathTime: info.PathTime,
		}
		if cfg.DebugLevel == "true" {
			log.Printf("当前监控的文件: %s, 修改时间: %v, 路径时间: %v\n",
				path, info.ModTime, info.PathTime)
		}
	}
	m.mu.RUnlock()

	for _, path := range paths {
		// 获取文件最后修改时间
		modTime, err := m.transporter.LastModified(path)
		if err != nil {
			if cfg.DebugLevel == "true" {
				log.Printf("获取文件修改时间失败: %s - %v\n", path, err)
			}
			m.events <- Event{
				Type:      EventDelete,
				Path:      path,
				Timestamp: time.Now(),
				Error:     err,
			}
			continue
		}

		if cfg.DebugLevel == "true" {
			log.Printf("文件 %s 的最新修改时间: %v\n", path, modTime)
		}

		info := fileInfos[path]
		if info == nil {
			if cfg.DebugLevel == "true" {
				log.Printf("文件 %s 的信息不存在\n", path)
			}
			continue
		}

		// 检查文件是否被修改（优先使用路径中的时间信息）
		isModified := false
		var newPathTime time.Time
		var pathTimeErr error

		// 尝试从路径中提取时间信息
		newPathTime, pathTimeErr = ExtractTimeFromPath(path)
		if pathTimeErr == nil {
			if info.PathTime.IsZero() || newPathTime.After(info.PathTime) {
				isModified = true
				if cfg.DebugLevel == "true" {
					log.Printf("检测到路径时间变化: %s, 原时间: %v, 新时间: %v\n",
						path, info.PathTime, newPathTime)
				}
			}
		} else if modTime.After(info.ModTime) {
			// 如果无法从路径获取时间，则使用修改时间
			isModified = true
			if cfg.DebugLevel == "true" {
				log.Printf("检测到文件修改时间变化: %s, 原时间: %v, 新时间: %v\n",
					path, info.ModTime, modTime)
			}
		}

		// 检查是否有新的时间目录
		if !isModified && !info.PathTime.IsZero() {
			// 获取当前路径的基本信息
			dir := filepath.ToSlash(filepath.Dir(filepath.Dir(path))) // 获取父目录的父目录，并统一使用正斜杠
			base := filepath.Base(path)                               // 获取文件名

			// 检查是否有更新的时间目录
			files, err := m.transporter.List(dir)
			if err == nil {
				for _, file := range files {
					// 统一使用正斜杠作为路径分隔符
					filePath := filepath.ToSlash(file.Path)
					if filepath.Base(filePath) == base {
						filePathTime, err := ExtractTimeFromPath(filePath)
						if err == nil && filePathTime.After(info.PathTime) {
							// 发现更新的时间目录
							isModified = true
							newPathTime = filePathTime
							path = filePath // 更新为新路径
							// 更新监控路径
							m.mu.Lock()
							delete(m.paths, info.Path) // 删除旧路径
							m.paths[path] = &FileInfo{ // 添加新路径
								Path:     path,
								ModTime:  modTime,
								PathTime: newPathTime,
							}
							m.mu.Unlock()
							if cfg.DebugLevel == "true" {
								log.Printf("检测到更新的时间目录: %s, 原时间: %v, 新时间: %v\n",
									filePath, info.PathTime, filePathTime)
							}
							break
						}
					}
				}
			}
		}

		if isModified {
			// 下载文件到临时目录
			tmpDir := filepath.Join(os.TempDir(), "configmanager")
			if err := fileutil.EnsureDir(tmpDir); err != nil {
				if cfg.DebugLevel == "true" {
					log.Printf("创建临时目录失败: %v\n", err)
				}
				continue
			}

			tmpFile := filepath.Join(tmpDir, filepath.Base(path))
			if err := m.transporter.Download(path, tmpFile); err != nil {
				if cfg.DebugLevel == "true" {
					log.Printf("下载文件失败: %s - %v\n", path, err)
				}
				continue
			}

			// 获取文件信息
			fileInfo, err := fileutil.GetFileInfo(tmpFile)
			if err != nil {
				if cfg.DebugLevel == "true" {
					log.Printf("获取文件信息失败: %s - %v\n", tmpFile, err)
				}
				continue
			}

			// 更新文件信息
			m.mu.Lock()
			if pathInfo, exists := m.paths[path]; exists {
				pathInfo.ModTime = modTime
				pathInfo.Size = fileInfo.Size
				if pathTimeErr == nil {
					pathInfo.PathTime = newPathTime
				}
				if cfg.DebugLevel == "true" {
					log.Printf("更新文件信息: %s, 修改时间: %v, 路径时间: %v, 大小: %d\n",
						path, modTime, pathInfo.PathTime, fileInfo.Size)
				}
			}
			m.mu.Unlock()

			// 发送修改事件
			event := Event{
				Type: EventModify,
				Path: path,
				FileInfo: &FileInfo{
					Path:     path,
					Size:     fileInfo.Size,
					ModTime:  modTime,
					PathTime: newPathTime,
				},
				Timestamp: time.Now(),
			}
			if cfg.DebugLevel == "true" {
				log.Printf("发送修改事件: %+v\n", event)
			}
			m.events <- event

			// 清理临时文件
			os.Remove(tmpFile)
		} else {
			if cfg.DebugLevel == "true" {
				log.Printf("文件未修改: %s, 当前修改时间: %v, 当前路径时间: %v\n",
					path, modTime, newPathTime)
			}
		}
	}
}
