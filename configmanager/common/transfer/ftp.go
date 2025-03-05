package transfer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"gin-server/config"

	"github.com/jlaffaye/ftp"
)

// FTPTransporter FTP文件传输实现
type FTPTransporter struct {
	host     string
	port     int
	username string
	password string
	conn     *ftp.ServerConn
	mu       sync.Mutex // 保护连接操作
}

// NewFTPTransporter 创建FTP传输器
func NewFTPTransporter(cfg *config.FTPConfig) (*FTPTransporter, error) {
	t := &FTPTransporter{
		host:     cfg.Host,
		port:     cfg.Port,
		username: cfg.Username,
		password: cfg.Password,
	}

	if err := t.connect(); err != nil {
		return nil, NewTransferError("connect", "", err)
	}

	return t, nil
}

// connect 连接FTP服务器
func (t *FTPTransporter) connect() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	// 如果已经有连接，先关闭
	if t.conn != nil {
		t.conn.Quit()
		t.conn = nil
	}

	// 连接服务器
	conn, err := ftp.Dial(fmt.Sprintf("%s:%d", t.host, t.port), ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return fmt.Errorf("连接FTP服务器失败: %w", err)
	}

	// 登录
	if err := conn.Login(t.username, t.password); err != nil {
		conn.Quit()
		return fmt.Errorf("FTP登录失败: %w", err)
	}

	t.conn = conn
	return nil
}

// ensureConnection 确保FTP连接可用
func (t *FTPTransporter) ensureConnection() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.conn == nil {
		return t.connect()
	}

	// 测试连接是否有效
	if err := t.conn.NoOp(); err != nil {
		// 重新连接
		return t.connect()
	}

	return nil
}

// Upload 上传文件到FTP服务器
func (t *FTPTransporter) Upload(localPath, remotePath string) error {
	if err := t.ensureConnection(); err != nil {
		return NewTransferError("upload", localPath, err)
	}

	// 创建远程目录
	dir := filepath.Dir(remotePath)
	if dir != "." && dir != "/" {
		t.createRemoteDir(dir)
	}

	// 打开本地文件
	file, err := os.Open(localPath)
	if err != nil {
		return NewTransferError("upload", localPath, fmt.Errorf("打开本地文件失败: %w", err))
	}
	defer file.Close()

	// 上传文件
	err = t.conn.Stor(remotePath, file)
	if err != nil {
		return NewTransferError("upload", remotePath, fmt.Errorf("上传文件失败: %w", err))
	}

	return nil
}

// Download 从FTP服务器下载文件
func (t *FTPTransporter) Download(remotePath, localPath string) error {
	if err := t.ensureConnection(); err != nil {
		return NewTransferError("download", remotePath, err)
	}

	// 创建本地目录
	if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
		return NewTransferError("download", localPath, fmt.Errorf("创建本地目录失败: %w", err))
	}

	// 获取远程文件
	resp, err := t.conn.Retr(remotePath)
	if err != nil {
		return NewTransferError("download", remotePath, fmt.Errorf("获取远程文件失败: %w", err))
	}
	defer resp.Close()

	// 创建本地文件
	file, err := os.Create(localPath)
	if err != nil {
		return NewTransferError("download", localPath, fmt.Errorf("创建本地文件失败: %w", err))
	}
	defer file.Close()

	// 复制文件内容
	_, err = io.Copy(file, resp)
	if err != nil {
		return NewTransferError("download", remotePath, fmt.Errorf("复制文件内容失败: %w", err))
	}

	return nil
}

// List 列出远程目录文件
func (t *FTPTransporter) List(remotePath string) ([]FileInfo, error) {
	if err := t.ensureConnection(); err != nil {
		return nil, NewTransferError("list", remotePath, err)
	}

	entries, err := t.conn.List(remotePath)
	if err != nil {
		return nil, NewTransferError("list", remotePath, fmt.Errorf("列出目录失败: %w", err))
	}

	files := make([]FileInfo, 0, len(entries))
	for _, entry := range entries {
		// 跳过 "." 和 ".." 目录
		if entry.Name == "." || entry.Name == ".." {
			continue
		}

		files = append(files, FileInfo{
			Name:    entry.Name,
			Size:    int64(entry.Size),
			ModTime: entry.Time,
			IsDir:   entry.Type == ftp.EntryTypeFolder,
			Path:    filepath.Join(remotePath, entry.Name),
		})
	}

	return files, nil
}

// Delete 删除远程文件
func (t *FTPTransporter) Delete(remotePath string) error {
	if err := t.ensureConnection(); err != nil {
		return NewTransferError("delete", remotePath, err)
	}

	err := t.conn.Delete(remotePath)
	if err != nil {
		return NewTransferError("delete", remotePath, fmt.Errorf("删除文件失败: %w", err))
	}

	return nil
}

// LastModified 获取远程文件最后修改时间
func (t *FTPTransporter) LastModified(remotePath string) (time.Time, error) {
	if err := t.ensureConnection(); err != nil {
		return time.Time{}, NewTransferError("lastModified", remotePath, err)
	}

	entries, err := t.conn.List(remotePath)
	if err != nil {
		return time.Time{}, NewTransferError("lastModified", remotePath, fmt.Errorf("获取文件信息失败: %w", err))
	}

	if len(entries) == 0 {
		return time.Time{}, NewTransferError("lastModified", remotePath, fmt.Errorf("文件不存在"))
	}

	return entries[0].Time, nil
}

// Close 关闭FTP连接
func (t *FTPTransporter) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.conn != nil {
		err := t.conn.Quit()
		t.conn = nil
		if err != nil {
			return NewTransferError("close", "", fmt.Errorf("关闭FTP连接失败: %w", err))
		}
	}
	return nil
}

// createRemoteDir 创建远程目录（递归）
func (t *FTPTransporter) createRemoteDir(path string) error {
	path = filepath.ToSlash(path) // 转换为正斜杠路径
	dirs := strings.Split(path, "/")
	currentPath := ""

	for _, dir := range dirs {
		if dir == "" {
			continue
		}
		currentPath = currentPath + "/" + dir
		err := t.conn.MakeDir(currentPath)
		if err != nil {
			// 忽略目录已存在的错误
			if !strings.Contains(err.Error(), "already exists") {
				return fmt.Errorf("创建目录 %s 失败: %w", currentPath, err)
			}
		}
	}
	return nil
}
