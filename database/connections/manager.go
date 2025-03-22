package connections

import (
	"fmt"
	"gin-server/config"
	"sync"
)

// ConnectionManager 数据库连接管理器
type ConnectionManager struct {
	connections map[string]Connection
	config      *config.Config
	mu          sync.RWMutex
}

// NewConnectionManager 创建连接管理器实例
func NewConnectionManager(cfg *config.Config) *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]Connection),
		config:      cfg,
	}
}

// RegisterConnection 注册数据库连接
func (m *ConnectionManager) RegisterConnection(name string, conn Connection) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 如果已存在同名连接，先关闭旧连接
	if oldConn, exists := m.connections[name]; exists {
		_ = oldConn.Close()
	}

	m.connections[name] = conn
}

// GetConnection 获取指定名称的数据库连接
func (m *ConnectionManager) GetConnection(name string) (Connection, error) {
	m.mu.RLock()
	conn, exists := m.connections[name]
	m.mu.RUnlock()

	if !exists {
		// 如果连接不存在，尝试创建并注册
		var err error
		conn, err = m.createConnection(name)
		if err != nil {
			return nil, err
		}

		// 注册新创建的连接
		m.RegisterConnection(name, conn)
	}

	return conn, nil
}

// createConnection 创建数据库连接
func (m *ConnectionManager) createConnection(name string) (Connection, error) {
	switch name {
	case "default", "radius":
		return NewMySQLConnection(name, m.config)
	default:
		return nil, fmt.Errorf("不支持的数据库连接类型: %s", name)
	}
}

// Default 获取默认数据库连接
func (m *ConnectionManager) Default() (Connection, error) {
	return m.GetConnection("default")
}

// Radius 获取Radius数据库连接
func (m *ConnectionManager) Radius() (Connection, error) {
	return m.GetConnection("radius")
}

// CloseAll 关闭所有数据库连接
func (m *ConnectionManager) CloseAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errs []error
	for name, conn := range m.connections {
		if err := conn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("关闭连接 %s 失败: %w", name, err))
		}
	}

	// 清空连接映射
	m.connections = make(map[string]Connection)

	// 如果有错误，返回第一个错误
	if len(errs) > 0 {
		return errs[0]
	}

	return nil
}
