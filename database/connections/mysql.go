package connections

import (
	"database/sql"
	"fmt"
	"gin-server/config"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// MySQLConnection MySQL数据库连接
type MySQLConnection struct {
	db     *gorm.DB
	sqlDB  *sql.DB
	name   string
	config *config.Config
}

// NewMySQLConnection 创建MySQL连接实例
func NewMySQLConnection(name string, cfg *config.Config) (Connection, error) {
	var dsn string

	// 根据连接名称选择不同的数据库配置
	switch name {
	case "default":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	case "radius":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.RadiusDBUser, cfg.RadiusDBPassword, cfg.RadiusDBHost, cfg.RadiusDBPort, cfg.RadiusDBName)
	default:
		return nil, fmt.Errorf("不支持的数据库连接名称: %s", name)
	}

	// 配置GORM
	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("连接MySQL数据库失败: %w", err)
	}

	// 获取并配置原生sql.DB连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取底层SQL数据库实例失败: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生命周期

	return &MySQLConnection{
		db:     db,
		sqlDB:  sqlDB,
		name:   name,
		config: cfg,
	}, nil
}

// GetDB 获取GORM数据库实例
func (c *MySQLConnection) GetDB() *gorm.DB {
	return c.db
}

// Close 关闭数据库连接
func (c *MySQLConnection) Close() error {
	return c.sqlDB.Close()
}
