package service

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gin-server/config"
	"gin-server/configmanager/common/alert"
	"gin-server/configmanager/common/transfer"
)

// File 表示远程文件信息
type File struct {
	Path string
	Size int64
}

// UploadContext 上传上下文
type UploadContext struct {
	// 文件路径
	LogPath   string    // 日志文件路径
	KeyPath   string    // 密钥文件路径(可选)
	Timestamp time.Time // 上传时间戳

	// 处理结果
	CompressedPath string // 压缩后的文件路径
	RemotePath     string // 远程存储路径
	TempDir        string // 临时目录路径

	// 配置和服务
	Config  *config.Config // 配置信息
	Alerter alert.Alerter  // 告警服务
}

// UploadStep 上传步骤接口
type UploadStep interface {
	// Execute 执行上传步骤
	Execute(ctx *UploadContext) error
}

// UploadError 上传错误
type UploadError struct {
	Step    string // 步骤名称
	Message string // 错误信息
	Err     error  // 原始错误
}

// Error 实现error接口
func (e *UploadError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s - %v", e.Step, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Step, e.Message)
}

// Unwrap 返回原始错误
func (e *UploadError) Unwrap() error {
	return e.Err
}

// NewUploadError 创建上传错误
func NewUploadError(step, message string, err error) error {
	return &UploadError{
		Step:    step,
		Message: message,
		Err:     err,
	}
}

// CompressStep 压缩步骤
type CompressStep struct{}

// NewCompressStep 创建压缩步骤
func NewCompressStep() *CompressStep {
	return &CompressStep{}
}

// Execute 执行压缩步骤
func (s *CompressStep) Execute(ctx *UploadContext) error {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "upload_*")
	if err != nil {
		return NewUploadError("compress", "创建临时目录失败", err)
	}
	ctx.TempDir = tempDir // 保存临时目录路径，供后续清理

	// 创建压缩文件
	timestamp := ctx.Timestamp.Format("20060102150405")
	compressedPath := filepath.Join(tempDir, fmt.Sprintf("%s.tar.gz", timestamp))

	// 创建tar.gz文件
	file, err := os.Create(compressedPath)
	if err != nil {
		os.RemoveAll(tempDir)
		return NewUploadError("compress", "创建压缩文件失败", err)
	}
	defer file.Close()

	// 创建gzip写入器
	gw := gzip.NewWriter(file)
	defer gw.Close()

	// 创建tar写入器
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// 添加日志文件
	if err := addFileToTar(tw, ctx.LogPath, filepath.Base(ctx.LogPath)); err != nil {
		os.RemoveAll(tempDir)
		return NewUploadError("compress", "添加日志文件失败", err)
	}

	// 如果有密钥文件，也添加到压缩包
	if ctx.KeyPath != "" {
		if err := addFileToTar(tw, ctx.KeyPath, filepath.Base(ctx.KeyPath)); err != nil {
			os.RemoveAll(tempDir)
			return NewUploadError("compress", "添加密钥文件失败", err)
		}
	}

	// 设置压缩后的文件路径
	ctx.CompressedPath = compressedPath
	return nil
}

// UploadTransferStep 上传传输步骤
type UploadTransferStep struct {
	transporter transfer.FileTransporter
}

// NewUploadTransferStep 创建上传传输步骤
func NewUploadTransferStep(cfg *config.Config) (*UploadTransferStep, error) {
	transporter, err := transfer.NewFileTransporter(transfer.TransporterTypeGitee, cfg)
	if err != nil {
		return nil, NewUploadError("transfer", "创建传输器失败", err)
	}

	return &UploadTransferStep{
		transporter: transporter,
	}, nil
}

// Execute 执行上传传输步骤
func (s *UploadTransferStep) Execute(ctx *UploadContext) error {
	// 规范化远程路径
	uploadDir := strings.ReplaceAll(ctx.Config.ConfigManager.LogManager.UploadDir, "\\", "/")
	if !strings.HasSuffix(uploadDir, "/") {
		uploadDir += "/"
	}

	// 构建远程路径
	remotePath := uploadDir + filepath.Base(ctx.CompressedPath)
	ctx.RemotePath = remotePath

	// 上传文件
	if err := s.transporter.Upload(ctx.CompressedPath, remotePath); err != nil {
		return NewUploadError("transfer", "上传文件失败", err)
	}

	return nil
}

// Close 关闭传输器
func (s *UploadTransferStep) Close() error {
	if s.transporter != nil {
		return s.transporter.Close()
	}
	return nil
}

// UploadManager 上传管理器
type UploadManager struct {
	config  *config.Config
	alerter alert.Alerter
	steps   []UploadStep
}

// NewUploadManager 创建上传管理器
func NewUploadManager(cfg *config.Config, alerter alert.Alerter) (*UploadManager, error) {
	// 创建上传步骤
	compressStep := NewCompressStep()
	uploadStep, err := NewUploadTransferStep(cfg)
	if err != nil {
		return nil, err
	}

	return &UploadManager{
		config:  cfg,
		alerter: alerter,
		steps:   []UploadStep{compressStep, uploadStep},
	}, nil
}

// Upload 上传文件
func (m *UploadManager) Upload(logPath, keyPath string) error {
	// 创建上传上下文
	ctx := &UploadContext{
		LogPath:   logPath,
		KeyPath:   keyPath,
		Timestamp: time.Now(),
		Config:    m.config,
		Alerter:   m.alerter,
	}

	// 执行每个步骤
	for _, step := range m.steps {
		if err := step.Execute(ctx); err != nil {
			if ctx.TempDir != "" {
				os.RemoveAll(ctx.TempDir)
			}
			return err
		}
	}

	// 清理临时目录
	if ctx.TempDir != "" {
		os.RemoveAll(ctx.TempDir)
	}

	return nil
}

// ListFiles 列出远程仓库中的文件
func (m *UploadManager) ListFiles(dir string) ([]File, error) {
	transporter, err := transfer.NewFileTransporter(transfer.TransporterTypeGitee, m.config)
	if err != nil {
		return nil, fmt.Errorf("创建传输器失败: %v", err)
	}
	defer transporter.Close()

	files, err := transporter.List(dir)
	if err != nil {
		return nil, fmt.Errorf("列出远程文件失败: %v", err)
	}

	var result []File
	for _, f := range files {
		result = append(result, File{
			Path: f.Path,
			Size: f.Size,
		})
	}
	return result, nil
}

// DownloadFile 从远程仓库下载文件
func (m *UploadManager) DownloadFile(remotePath string) ([]byte, error) {
	transporter, err := transfer.NewFileTransporter(transfer.TransporterTypeGitee, m.config)
	if err != nil {
		return nil, fmt.Errorf("创建传输器失败: %v", err)
	}
	defer transporter.Close()

	// 创建临时文件
	tempFile := filepath.Join(os.TempDir(), filepath.Base(remotePath))
	if err := transporter.Download(remotePath, tempFile); err != nil {
		return nil, fmt.Errorf("下载文件失败: %v", err)
	}
	defer os.Remove(tempFile)

	// 读取文件内容
	data, err := os.ReadFile(tempFile)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	return data, nil
}

// Close 关闭上传管理器
func (m *UploadManager) Close() error {
	// 关闭需要清理的步骤
	for _, step := range m.steps {
		if closer, ok := step.(interface{ Close() error }); ok {
			if err := closer.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

// addFileToTar 添加文件到tar
func addFileToTar(tarWriter *tar.Writer, src, relPath string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	header := &tar.Header{
		Name:    relPath,
		Size:    stat.Size(),
		Mode:    int64(stat.Mode()),
		ModTime: stat.ModTime(),
	}

	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}

	_, err = io.Copy(tarWriter, file)
	return err
}

// UploadFile 直接上传文件数据到远程仓库
func (m *UploadManager) UploadFile(remotePath string, data []byte) error {
	transporter, err := transfer.NewFileTransporter(transfer.TransporterTypeGitee, m.config)
	if err != nil {
		return fmt.Errorf("创建传输器失败: %v", err)
	}
	defer transporter.Close()

	// 创建临时文件
	tempFile := filepath.Join(os.TempDir(), filepath.Base(remotePath))
	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		return fmt.Errorf("创建临时文件失败: %v", err)
	}
	defer os.Remove(tempFile)

	// 上传文件
	if err := transporter.Upload(tempFile, remotePath); err != nil {
		return fmt.Errorf("上传文件失败: %v", err)
	}

	return nil
}
