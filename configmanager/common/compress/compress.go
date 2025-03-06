package compress

import (
	"fmt"
	"io"
	"log"

	"gin-server/config"
)

// CompressFormat 压缩格式类型
type CompressFormat string

const (
	// FormatTarGz tar.gz格式
	FormatTarGz CompressFormat = "tar.gz"
)

// ProgressCallback 进度回调函数
type ProgressCallback func(current, total int64)

// Options 压缩选项
type Options struct {
	// CompressionLevel 压缩级别 (1-9)
	CompressionLevel int
	// BufferSize 缓冲区大小
	BufferSize int
	// IncludePatterns 包含的文件模式
	IncludePatterns []string
	// ExcludePatterns 排除的文件模式
	ExcludePatterns []string
	// ProgressCallback 进度回调
	ProgressCallback ProgressCallback
}

// Option 选项设置函数
type Option func(*Options)

// DefaultOptions 默认选项
var DefaultOptions = Options{
	CompressionLevel: 6,
	BufferSize:       32 * 1024, // 32KB
}

// Compressor 压缩器接口
type Compressor interface {
	// Compress 压缩文件或目录
	// src: 源文件或目录路径
	// dest: 目标文件路径
	// options: 压缩选项
	Compress(src string, dest string, options ...Option) error

	// Decompress 解压文件
	// src: 源文件路径
	// dest: 目标目录路径
	// options: 解压选项
	Decompress(src string, dest string, options ...Option) error
}

// CompressError 压缩错误
type CompressError struct {
	Operation string // 操作名称
	Path      string // 文件路径
	Err       error  // 原始错误
}

// Error 实现error接口
func (e *CompressError) Error() string {
	if e.Path == "" {
		return fmt.Sprintf("%s: %v", e.Operation, e.Err)
	}
	return fmt.Sprintf("%s %s: %v", e.Operation, e.Path, e.Err)
}

// Unwrap 返回原始错误
func (e *CompressError) Unwrap() error {
	return e.Err
}

// NewCompressError 创建压缩错误
func NewCompressError(op, path string, err error) error {
	return &CompressError{
		Operation: op,
		Path:      path,
		Err:       err,
	}
}

// WithCompressionLevel 设置压缩级别
func WithCompressionLevel(level int) Option {
	return func(o *Options) {
		if level >= 1 && level <= 9 {
			o.CompressionLevel = level
		}
	}
}

// WithBufferSize 设置缓冲区大小
func WithBufferSize(size int) Option {
	return func(o *Options) {
		if size > 0 {
			o.BufferSize = size
		}
	}
}

// WithIncludePatterns 设置包含的文件模式
func WithIncludePatterns(patterns []string) Option {
	return func(o *Options) {
		o.IncludePatterns = patterns
	}
}

// WithExcludePatterns 设置排除的文件模式
func WithExcludePatterns(patterns []string) Option {
	return func(o *Options) {
		o.ExcludePatterns = patterns
	}
}

// WithProgressCallback 设置进度回调
func WithProgressCallback(callback ProgressCallback) Option {
	return func(o *Options) {
		o.ProgressCallback = callback
	}
}

// processOptions 处理选项
func processOptions(opts ...Option) Options {
	options := DefaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	return options
}

// copyWithProgress 带进度的复制
func copyWithProgress(dst io.Writer, src io.Reader, size int64, callback ProgressCallback, bufSize int) error {
	if callback == nil {
		_, err := io.CopyBuffer(dst, src, make([]byte, bufSize))
		return err
	}

	buf := make([]byte, bufSize)
	var written int64

	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = fmt.Errorf("invalid write result")
				}
			}
			written += int64(nw)
			if ew != nil {
				return ew
			}
			if nr != nw {
				return io.ErrShortWrite
			}
			callback(written, size)
		}
		if er != nil {
			if er != io.EOF {
				return er
			}
			break
		}
	}
	return nil
}

// logDebug 输出调试日志
func logDebug(format string, v ...interface{}) {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Printf(format, v...)
	}
}
