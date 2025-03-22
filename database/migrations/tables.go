package migrations

import (
	"fmt"
	"gin-server/config"
	"log"

	"gorm.io/gorm"
)

// TableChecker 表检查器接口，用于定义表检查和创建的行为
type TableChecker interface {
	// TableName 返回表名
	TableName() string

	// Check 检查表是否存在且结构是否正确
	Check(db *gorm.DB) (bool, error)

	// Create 创建表
	Create(db *gorm.DB) error
}

// RegisteredTables 保存所有注册的表检查器
var RegisteredTables []TableChecker

// RegisterTable 注册一个表检查器
func RegisterTable(checker TableChecker) {
	RegisteredTables = append(RegisteredTables, checker)
}

// LogFilesTableChecker log_files表检查器
type LogFilesTableChecker struct{}

// TableName 返回表名
func (l *LogFilesTableChecker) TableName() string {
	return "log_files"
}

// Check 检查log_files表是否存在且结构是否正确
func (l *LogFilesTableChecker) Check(db *gorm.DB) (bool, error) {
	cfg := config.GetConfig()

	// 检查表是否存在
	var tableExists int
	err := db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ? AND table_name = ?",
		cfg.DBName, l.TableName()).Scan(&tableExists).Error
	if err != nil {
		return false, fmt.Errorf("检查%s表是否存在时出错: %w", l.TableName(), err)
	}

	return tableExists > 0, nil
}

// Create 创建log_files表
func (l *LogFilesTableChecker) Create(db *gorm.DB) error {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Printf("%s表不存在，开始创建...", l.TableName())
	}

	// 表不存在，创建表
	createTableSQL := `CREATE TABLE IF NOT EXISTS log_files (
		id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
		created_at DATETIME(3) NULL,
		updated_at DATETIME(3) NULL,
		deleted_at DATETIME(3) NULL,
		file_name VARCHAR(255) NOT NULL,
		file_path VARCHAR(255) NOT NULL,
		file_size BIGINT NOT NULL,
		start_time DATETIME(3) NOT NULL,
		end_time DATETIME(3) NOT NULL,
		is_encrypted BOOLEAN NOT NULL DEFAULT FALSE,
		is_uploaded BOOLEAN NOT NULL DEFAULT FALSE,
		remote_path VARCHAR(255) NOT NULL DEFAULT '',
		uploaded_time DATETIME(3) NULL,
		INDEX idx_log_files_deleted_at (deleted_at),
		UNIQUE INDEX idx_log_files_file_name (file_name),
		INDEX idx_log_files_start_time (start_time),
		INDEX idx_log_files_end_time (end_time)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`

	// 执行创建表操作
	if err := db.Exec(createTableSQL).Error; err != nil {
		return fmt.Errorf("创建%s表失败: %w", l.TableName(), err)
	}

	// 检查表是否创建成功
	exists, err := l.Check(db)
	if err != nil {
		return fmt.Errorf("验证%s表创建是否成功时出错: %w", l.TableName(), err)
	}

	if !exists {
		return fmt.Errorf("%s表创建失败，没有找到新创建的表", l.TableName())
	}

	if cfg.DebugLevel == "true" {
		log.Printf("%s表创建成功", l.TableName())
	}

	return nil
}

// RadPostAuthTableChecker radpostauth表检查器
type RadPostAuthTableChecker struct{}

// TableName 返回表名
func (r *RadPostAuthTableChecker) TableName() string {
	return "radpostauth"
}

// Check 检查radpostauth表是否存在且结构是否正确
func (r *RadPostAuthTableChecker) Check(db *gorm.DB) (bool, error) {
	cfg := config.GetConfig()

	// 检查表是否存在
	var tableExists int
	err := db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ? AND table_name = ?",
		cfg.RadiusDBName, r.TableName()).Scan(&tableExists).Error
	if err != nil {
		return false, fmt.Errorf("检查%s表是否存在时出错: %w", r.TableName(), err)
	}

	return tableExists > 0, nil
}

// Create 创建radpostauth表
func (r *RadPostAuthTableChecker) Create(db *gorm.DB) error {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Printf("%s表不存在，开始创建...", r.TableName())
	}

	// 表不存在，创建表
	createTableSQL := `CREATE TABLE IF NOT EXISTS radpostauth (
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(64) NOT NULL,
		pass VARCHAR(64) NOT NULL,
		reply VARCHAR(32) NOT NULL,
		authdate TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
		class VARCHAR(64) NULL DEFAULT NULL,
		INDEX idx_username (username),
		INDEX idx_class (class)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`

	// 执行创建表操作
	if err := db.Exec(createTableSQL).Error; err != nil {
		return fmt.Errorf("创建%s表失败: %w", r.TableName(), err)
	}

	// 检查表是否创建成功
	exists, err := r.Check(db)
	if err != nil {
		return fmt.Errorf("验证%s表创建是否成功时出错: %w", r.TableName(), err)
	}

	if !exists {
		return fmt.Errorf("%s表创建失败，没有找到新创建的表", r.TableName())
	}

	if cfg.DebugLevel == "true" {
		log.Printf("%s表创建成功", r.TableName())
	}

	return nil
}

// 初始化代码，注册所有表检查器
func init() {
	// 注册log_files表检查器
	RegisterTable(&LogFilesTableChecker{})

	// 注册radpostauth表检查器
	RegisterTable(&RadPostAuthTableChecker{})

	// 在这里可以注册其他表检查器
	// 例如: RegisterTable(&XXXTableChecker{})
}

// EnsureAllTablesExist 确保所有注册的表都存在
func EnsureAllTablesExist(db *gorm.DB) error {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Printf("开始检查所有关键表是否存在...")
	}

	for _, checker := range RegisteredTables {
		// 检查表是否存在
		exists, err := checker.Check(db)
		if err != nil {
			return fmt.Errorf("检查表%s时出错: %w", checker.TableName(), err)
		}

		// 如果表不存在，创建它
		if !exists {
			if err := checker.Create(db); err != nil {
				return fmt.Errorf("创建表%s时出错: %w", checker.TableName(), err)
			}
		} else if cfg.DebugLevel == "true" {
			log.Printf("表%s已存在，无需创建", checker.TableName())
		}
	}

	if cfg.DebugLevel == "true" {
		log.Printf("所有关键表检查完成")
	}

	return nil
}
