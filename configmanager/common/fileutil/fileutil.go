package fileutil

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"gin-server/config"
)

// FileInfo 文件信息结构体
type FileInfo struct {
	Path         string    `json:"path"`           // 文件路径
	Name         string    `json:"name"`           // 文件名
	Size         int64     `json:"size"`           // 文件大小
	ModTime      time.Time `json:"mod_time"`       // 修改时间
	IsDir        bool      `json:"is_dir"`         // 是否是目录
	Permissions  string    `json:"permissions"`    // 权限字符串
	BackupPath   string    `json:"backup_path"`    // 备份路径
	LastBackupAt time.Time `json:"last_backup_at"` // 最后备份时间
}

// EnsureDir 确保目录存在，如果不存在则创建
func EnsureDir(path string) error {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Printf("确保目录存在: %s\n", path)
	}

	err := os.MkdirAll(path, 0755)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("创建目录失败: %v\n", err)
		}
		return fmt.Errorf("创建目录失败: %w", err)
	}

	if cfg.DebugLevel == "true" {
		log.Printf("目录已就绪: %s\n", path)
	}
	return nil
}

// CopyFile 复制文件
func CopyFile(src, dst string) error {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Printf("复制文件: %s -> %s\n", src, dst)
	}

	// 打开源文件
	sourceFile, err := os.Open(src)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("打开源文件失败: %v\n", err)
		}
		return fmt.Errorf("打开源文件失败: %w", err)
	}
	defer sourceFile.Close()

	// 创建目标文件
	destFile, err := os.Create(dst)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("创建目标文件失败: %v\n", err)
		}
		return fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer destFile.Close()

	// 复制文件内容
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("复制文件内容失败: %v\n", err)
		}
		return fmt.Errorf("复制文件内容失败: %w", err)
	}

	// 获取源文件权限
	sourceInfo, err := os.Stat(src)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("获取源文件信息失败: %v\n", err)
		}
		return fmt.Errorf("获取源文件信息失败: %w", err)
	}

	// 设置目标文件权限
	err = os.Chmod(dst, sourceInfo.Mode())
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("设置目标文件权限失败: %v\n", err)
		}
		return fmt.Errorf("设置目标文件权限失败: %w", err)
	}

	if cfg.DebugLevel == "true" {
		log.Printf("文件复制完成: %s -> %s\n", src, dst)
	}
	return nil
}

// BackupFile 备份文件
func BackupFile(src string) (string, error) {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Printf("开始备份文件: %s\n", src)
	}

	// 检查源文件是否存在
	if _, err := os.Stat(src); os.IsNotExist(err) {
		if cfg.DebugLevel == "true" {
			log.Printf("源文件不存在: %s\n", src)
		}
		return "", fmt.Errorf("源文件不存在: %w", err)
	}

	// 创建备份目录
	backupDir := filepath.Join(filepath.Dir(src), "backup")
	if err := EnsureDir(backupDir); err != nil {
		return "", err
	}

	// 生成备份文件名
	timestamp := time.Now().Format("20060102150405")
	ext := filepath.Ext(src)
	baseName := filepath.Base(src[:len(src)-len(ext)])
	backupPath := filepath.Join(backupDir, fmt.Sprintf("%s_%s%s", baseName, timestamp, ext))

	// 复制文件到备份位置
	if err := CopyFile(src, backupPath); err != nil {
		return "", err
	}

	if cfg.DebugLevel == "true" {
		log.Printf("文件备份完成: %s -> %s\n", src, backupPath)
	}
	return backupPath, nil
}

// GetFileInfo 获取文件信息
func GetFileInfo(path string) (*FileInfo, error) {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Printf("获取文件信息: %s\n", path)
	}

	info, err := os.Stat(path)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("获取文件信息失败: %v\n", err)
		}
		return nil, fmt.Errorf("获取文件信息失败: %w", err)
	}

	fileInfo := &FileInfo{
		Path:        path,
		Name:        info.Name(),
		Size:        info.Size(),
		ModTime:     info.ModTime(),
		IsDir:       info.IsDir(),
		Permissions: info.Mode().String(),
	}

	if cfg.DebugLevel == "true" {
		log.Printf("文件信息: %+v\n", fileInfo)
	}
	return fileInfo, nil
}

// IsFileExists 检查文件是否存在
func IsFileExists(path string) bool {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Printf("检查文件是否存在: %s\n", path)
	}

	_, err := os.Stat(path)
	exists := !os.IsNotExist(err)

	if cfg.DebugLevel == "true" {
		log.Printf("文件 %s 存在: %v\n", path, exists)
	}
	return exists
}

// WriteFile 写入文件
func WriteFile(path string, data []byte, perm os.FileMode) error {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Printf("写入文件: %s, 数据长度: %d\n", path, len(data))
	}

	// 确保目录存在
	dir := filepath.Dir(path)
	if err := EnsureDir(dir); err != nil {
		return err
	}

	// 写入文件
	err := os.WriteFile(path, data, perm)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("写入文件失败: %v\n", err)
		}
		return fmt.Errorf("写入文件失败: %w", err)
	}

	if cfg.DebugLevel == "true" {
		log.Printf("文件写入完成: %s\n", path)
	}
	return nil
}

// ReadFile 读取文件
func ReadFile(path string) ([]byte, error) {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Printf("读取文件: %s\n", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if cfg.DebugLevel == "true" {
			log.Printf("读取文件失败: %v\n", err)
		}
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	if cfg.DebugLevel == "true" {
		log.Printf("文件读取完成: %s, 数据长度: %d\n", path, len(data))
	}
	return data, nil
}
