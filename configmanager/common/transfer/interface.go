package transfer

import "time"

// TransporterType 传输器类型
type TransporterType string

const (
	// TransporterTypeGitee Gitee传输器
	TransporterTypeGitee TransporterType = "gitee"
	// TransporterTypeFTP FTP传输器
	TransporterTypeFTP TransporterType = "ftp"
)

// FileInfo 文件信息
type FileInfo struct {
	Name    string    // 文件名
	Size    int64     // 文件大小
	ModTime time.Time // 修改时间
	IsDir   bool      // 是否是目录
	Path    string    // 相对路径
	Hash    string    // 文件哈希（可选）
}

// FileTransporter 文件传输器接口
type FileTransporter interface {
	// Upload 上传文件
	Upload(localPath, remotePath string) error

	// Download 下载文件
	Download(remotePath, localPath string) error

	// List 列出目录内容
	List(remotePath string) ([]FileInfo, error)

	// Delete 删除文件
	Delete(remotePath string) error

	// LastModified 获取文件最后修改时间
	LastModified(remotePath string) (time.Time, error)

	// Close 关闭传输器
	Close() error
}
