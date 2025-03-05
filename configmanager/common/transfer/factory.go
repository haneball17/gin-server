package transfer

import (
	"fmt"

	"gin-server/config"
)

// NewFileTransporter 创建文件传输器
func NewFileTransporter(typ TransporterType, cfg *config.Config) (FileTransporter, error) {
	switch typ {
	case TransporterTypeGitee:
		if cfg.Gitee == nil {
			return nil, fmt.Errorf("Gitee配置为空")
		}
		return NewGiteeTransporter(cfg.Gitee), nil
	case TransporterTypeFTP:
		if cfg.FTP == nil {
			return nil, fmt.Errorf("FTP配置为空")
		}
		transporter, err := NewFTPTransporter(cfg.FTP)
		if err != nil {
			return nil, fmt.Errorf("创建FTP传输器失败: %w", err)
		}
		return transporter, nil
	default:
		return nil, fmt.Errorf("不支持的传输器类型: %s", typ)
	}
}
