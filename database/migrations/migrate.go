package migrations

import (
	"fmt"
	"gin-server/config"
	"gin-server/database/models"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate(db *gorm.DB) error {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Println("开始自动迁移数据库表结构...")
	}

	// 设置GORM全局命名策略
	db.NamingStrategy = schema.NamingStrategy{
		SingularTable: true,                  // 使用单数表名
		TablePrefix:   "",                    // 表前缀
		NameReplacer:  strings.NewReplacer(), // 名称替换器
		NoLowerCase:   false,                 // 不使用小写
	}

	// 在执行迁移前修复可能存在的数据问题
	if err := fixDuplicateUserIDs(db); err != nil {
		return fmt.Errorf("修复用户ID冲突失败: %w", err)
	}

	// 修复user_behaviors表结构
	if err := fixUserBehaviorsTable(db); err != nil {
		return fmt.Errorf("修复user_behaviors表结构失败: %w", err)
	}

	// 修复devices表结构，确保新字段存在
	if err := fixDevicesTable(db); err != nil {
		return fmt.Errorf("修复devices表结构失败: %w", err)
	}

	// 需要迁移的主数据库模型
	migrationModels := []interface{}{
		&models.User{},
		&models.Event{},
		&models.Device{},
		&models.UserBehavior{},
		&models.LogFile{},
		&models.Cert{},
	}

	// 执行主数据库迁移
	for _, model := range migrationModels {
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("自动迁移模型 %T 失败: %w", model, err)
		}

		if cfg.DebugLevel == "true" {
			log.Printf("模型 %T 迁移成功", model)
		}
	}

	// 迁移旧表（如果有）
	if err := migrateOldTables(db); err != nil {
		return fmt.Errorf("迁移旧表结构失败: %w", err)
	}

	// 确保所有注册的关键表存在
	if err := EnsureAllTablesExist(db); err != nil {
		return fmt.Errorf("确保关键表存在失败: %w", err)
	}

	if cfg.DebugLevel == "true" {
		log.Println("数据库表结构迁移完成")
	}

	return nil
}

// migrateOldTables 迁移旧的数据库表结构
func migrateOldTables(db *gorm.DB) error {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Println("开始检查旧数据库表结构...")
	}

	// 用户行为表定义
	userBehaviorsTable := TableDefinition{
		Name: "user_behaviors",
		Columns: []Column{
			{Name: "behavior_id", Type: "INT", PrimaryKey: true, Nullable: false, Comment: "AUTO_INCREMENT PRIMARY KEY"},
			{Name: "user_id", Type: "INT", Nullable: false, Comment: "用户ID"},
			{Name: "behavior_time", Type: "DATETIME(3)", Nullable: false, Comment: "行为开始时间"},
			{Name: "behavior_type", Type: "INT", Nullable: false, Comment: "行为类型，1:发送，2:接收"},
			{Name: "data_type", Type: "INT", Nullable: false, Comment: "数据类型，1:文件，2:消息"},
			{Name: "data_size", Type: "BIGINT", Nullable: false, Comment: "数据大小"},
			{Name: "created_at", Type: "TIMESTAMP(3)", Nullable: true, Comment: "创建时间"},
		},
		Indexes: []string{
			"idx_user_id (user_id)",
			"idx_behavior_time (behavior_time)",
			"idx_behavior_type (behavior_type)",
		},
		ForeignKey: &ForeignKey{
			Column:    "user_id",
			RefTable:  "users",
			RefColumn: "user_id",
		},
	}

	// 事件表定义
	eventsTable := TableDefinition{
		Name: "events",
		Columns: []Column{
			{Name: "event_id", Type: "BIGINT", PrimaryKey: true, Nullable: false, Comment: "事件ID"},
			{Name: "device_id", Type: "INT", Nullable: false, Comment: "设备ID"},
			{Name: "event_time", Type: "DATETIME", Nullable: false, Comment: "事件发生时间"},
			{Name: "event_type", Type: "INT", Nullable: false, Comment: "事件类型，1:安全事件，2:故障事件"},
			{Name: "event_code", Type: "VARCHAR(20)", Nullable: false, Comment: "事件代码"},
			{Name: "event_desc", Type: "VARCHAR(255)", Nullable: false, Comment: "事件描述"},
			{Name: "created_at", Type: "TIMESTAMP", Nullable: true, Comment: "创建时间"},
		},
		Indexes: []string{},
	}

	// 检查用户行为表
	if err := checkAndMigrateTable(db, userBehaviorsTable); err != nil {
		return fmt.Errorf("迁移用户行为表失败: %w", err)
	}

	// 检查事件表
	if err := checkAndMigrateTable(db, eventsTable); err != nil {
		return fmt.Errorf("迁移事件表失败: %w", err)
	}

	if cfg.DebugLevel == "true" {
		log.Println("旧数据库表结构检查完成！")
	}
	return nil
}

// 列定义结构体
type Column struct {
	Name       string
	Type       string
	Nullable   bool
	PrimaryKey bool
	Comment    string
}

// 表结构定义
type TableDefinition struct {
	Name       string
	Columns    []Column
	Indexes    []string
	ForeignKey *ForeignKey
}

// 外键定义
type ForeignKey struct {
	Column    string
	RefTable  string
	RefColumn string
}

// checkAndMigrateTable 检查并迁移指定表
func checkAndMigrateTable(db *gorm.DB, tableDef TableDefinition) error {
	// 检查表是否存在
	var tableExists int
	row := db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ? AND table_name = ?",
		config.GetConfig().DBName, tableDef.Name).Row()
	row.Scan(&tableExists)

	if tableExists == 0 {
		// 构建创建表的SQL语句
		sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n", tableDef.Name)

		// 添加列定义
		var columnDefs []string
		for _, col := range tableDef.Columns {
			var colDef string
			if col.PrimaryKey {
				colDef = fmt.Sprintf("  %s %s AUTO_INCREMENT PRIMARY KEY", col.Name, col.Type)
			} else {
				nullable := "NOT NULL"
				if col.Nullable {
					nullable = "NULL"
					if col.Name == "created_at" && strings.HasPrefix(col.Type, "TIMESTAMP") {
						nullable += " DEFAULT CURRENT_TIMESTAMP"
					}
				}
				colDef = fmt.Sprintf("  %s %s %s", col.Name, col.Type, nullable)
				if col.Comment != "" {
					colDef += fmt.Sprintf(" COMMENT '%s'", col.Comment)
				}
			}
			columnDefs = append(columnDefs, colDef)
		}

		// 添加索引
		for _, idx := range tableDef.Indexes {
			columnDefs = append(columnDefs, fmt.Sprintf("  INDEX %s", idx))
		}

		// 添加外键约束
		if tableDef.ForeignKey != nil {
			fk := tableDef.ForeignKey
			fkDef := fmt.Sprintf("  FOREIGN KEY (%s) REFERENCES %s(%s) ON DELETE CASCADE ON UPDATE CASCADE",
				fk.Column, fk.RefTable, fk.RefColumn)
			columnDefs = append(columnDefs, fkDef)
		}

		sql += strings.Join(columnDefs, ",\n")
		sql += "\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;"

		// 调用通用的表创建函数
		return EnsureTableExists(db, tableDef.Name, sql)
	}

	// 表存在，检查字段
	return checkAndFixColumns(db, tableDef)
}

// checkAndFixColumns 检查并修复列
func checkAndFixColumns(db *gorm.DB, tableDef TableDefinition) error {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Printf("检查表 %s 的字段结构...\n", tableDef.Name)
	}

	// 获取表中的列
	type ColumnInfo struct {
		Field   string `gorm:"column:Field"`
		Type    string `gorm:"column:Type"`
		Null    string `gorm:"column:Null"`
		Key     string `gorm:"column:Key"`
		Default string `gorm:"column:Default"`
		Extra   string `gorm:"column:Extra"`
	}

	var columns []ColumnInfo
	if err := db.Raw(fmt.Sprintf("SHOW COLUMNS FROM %s", tableDef.Name)).Scan(&columns).Error; err != nil {
		return fmt.Errorf("获取表 %s 的列信息失败: %w", tableDef.Name, err)
	}

	// 检查是否有缺失的列
	existingColumns := make(map[string]ColumnInfo)
	for _, col := range columns {
		existingColumns[col.Field] = col
	}

	var missingColumns []Column
	for _, expectedCol := range tableDef.Columns {
		if _, exists := existingColumns[expectedCol.Name]; !exists {
			missingColumns = append(missingColumns, expectedCol)
		}
	}

	// 添加缺失的列
	if len(missingColumns) > 0 {
		if cfg.DebugLevel == "true" {
			log.Printf("表 %s 中发现 %d 个缺失的列，正在添加...\n", tableDef.Name, len(missingColumns))
		}

		for _, col := range missingColumns {
			nullable := "NOT NULL"
			if col.Nullable {
				nullable = "NULL"
				if col.Name == "created_at" && strings.HasPrefix(col.Type, "TIMESTAMP") {
					nullable += " DEFAULT CURRENT_TIMESTAMP"
				}
			}

			var comment string
			if col.Comment != "" {
				comment = fmt.Sprintf(" COMMENT '%s'", col.Comment)
			}

			sql := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s %s%s",
				tableDef.Name, col.Name, col.Type, nullable, comment)

			if err := db.Exec(sql).Error; err != nil {
				return fmt.Errorf("向表 %s 添加列 %s 失败: %w", tableDef.Name, col.Name, err)
			}

			if cfg.DebugLevel == "true" {
				log.Printf("已向表 %s 添加列 %s\n", tableDef.Name, col.Name)
			}
		}
	}

	return nil
}

// SeedData 填充测试数据（可选）
func SeedData(db *gorm.DB) error {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Println("开始填充测试数据...")
	}

	// 这里可以添加填充测试数据的代码
	// 例如：创建一些初始设备、事件等

	if cfg.DebugLevel == "true" {
		log.Println("测试数据填充完成！")
	}
	return nil
}

// EnsureTableExists 确保指定表存在，如果不存在则手动创建
func EnsureTableExists(db *gorm.DB, tableName string, createTableSQL string) error {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Printf("检查%s表是否存在...", tableName)
	}

	// 检查表是否存在
	var tableExists int
	row := db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ? AND table_name = ?",
		config.GetConfig().DBName, tableName).Row()
	row.Scan(&tableExists)

	if tableExists == 0 {
		if cfg.DebugLevel == "true" {
			log.Printf("%s表不存在，手动创建中...", tableName)
		}

		if err := db.Exec(createTableSQL).Error; err != nil {
			return fmt.Errorf("手动创建%s表失败: %w", tableName, err)
		}

		if cfg.DebugLevel == "true" {
			log.Printf("手动创建%s表成功", tableName)
		}
	} else if cfg.DebugLevel == "true" {
		log.Printf("%s表已存在", tableName)
	}

	return nil
}

// EnsureLogFilesTableExists 确保log_files表存在，如果不存在则手动创建
// 保留此函数以保持向后兼容性
func EnsureLogFilesTableExists(db *gorm.DB) error {
	// 查找并使用LogFilesTableChecker
	for _, checker := range RegisteredTables {
		if checker.TableName() == "log_files" {
			// 检查表是否存在
			exists, err := checker.Check(db)
			if err != nil {
				return err
			}

			// 如果表不存在，创建它
			if !exists {
				return checker.Create(db)
			}

			// 表已存在
			return nil
		}
	}

	// 如果未找到LogFilesTableChecker，使用旧代码作为后备
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Printf("警告：未找到LogFilesTableChecker，使用旧的实现...")
	}

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

	return db.Exec(createTableSQL).Error
}

// fixDuplicateUserIDs 修复用户ID冲突问题
func fixDuplicateUserIDs(db *gorm.DB) error {
	cfg := config.GetConfig()

	// 查找user_id为0的记录数
	var count int64
	if err := db.Table("users").Where("user_id = 0").Count(&count).Error; err != nil {
		// 如果表不存在，则忽略错误（表将在后续迁移步骤中创建）
		if strings.Contains(err.Error(), "doesn't exist") {
			if cfg.DebugLevel == "true" {
				log.Println("users表不存在，将在迁移过程中创建")
			}
			return nil
		}
		return err
	}

	// 如果没有重复值，直接返回
	if count <= 1 {
		if cfg.DebugLevel == "true" && count == 1 {
			log.Println("用户表中只有一条user_id为0的记录，无需修复")
		}
		return nil
	}

	if cfg.DebugLevel == "true" {
		log.Printf("发现%d条user_id为0的记录，开始修复...", count)
	}

	// 使用事务处理更新操作
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	// 查找所有user_id为0的记录
	type UserRecord struct {
		ID     uint
		UserID int
	}
	var users []UserRecord
	if err := tx.Table("users").Select("id, user_id").Where("user_id = 0").Find(&users).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 为每条记录分配一个唯一的负数ID（保留第一条为0）
	startID := -1
	for i, user := range users {
		if i == 0 {
			// 保留第一条记录的user_id为0
			continue
		}

		// 为其他记录分配唯一的负数ID
		if err := tx.Table("users").Where("id = ?", user.ID).Update("user_id", startID).Error; err != nil {
			tx.Rollback()
			return err
		}

		if cfg.DebugLevel == "true" {
			log.Printf("更新用户ID:%d 的user_id为:%d", user.ID, startID)
		}

		startID--
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return err
	}

	if cfg.DebugLevel == "true" {
		log.Println("成功修复用户表中的重复user_id值")
	}

	return nil
}

// fixUserBehaviorsTable 修复user_behaviors表的behavior_id字段问题
func fixUserBehaviorsTable(db *gorm.DB) error {
	cfg := config.GetConfig()

	// 始终输出日志，不依赖于DebugLevel
	log.Println("[数据库修复] 开始检查user_behaviors表结构...")

	// 检查user_behaviors表是否存在
	var tableExists int
	err := db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ? AND table_name = ?",
		cfg.DBName, "user_behaviors").Scan(&tableExists).Error
	if err != nil {
		log.Printf("[数据库修复] 错误：检查user_behaviors表是否存在失败: %v", err)
		return fmt.Errorf("检查user_behaviors表是否存在时出错: %w", err)
	}

	if tableExists == 0 {
		log.Println("[数据库修复] user_behaviors表不存在，将通过AutoMigrate创建")
		return nil
	}

	log.Println("[数据库修复] user_behaviors表存在，继续检查字段")

	// 首先检查behavior_id列是否正确设置为自增主键
	var behaviorIDInfo struct {
		Field   string `gorm:"column:Field"`
		Type    string `gorm:"column:Type"`
		Null    string `gorm:"column:Null"`
		Key     string `gorm:"column:Key"`
		Default string `gorm:"column:Default"`
		Extra   string `gorm:"column:Extra"`
	}

	err = db.Raw("SHOW COLUMNS FROM user_behaviors WHERE Field = 'behavior_id'").Scan(&behaviorIDInfo).Error
	if err != nil {
		log.Printf("[数据库修复] 错误：获取behavior_id列信息失败: %v", err)
	} else {
		log.Printf("[数据库修复] 检查behavior_id列: 类型=%s, 可空=%s, 键类型=%s, 默认值=%s, 额外=%s",
			behaviorIDInfo.Type, behaviorIDInfo.Null, behaviorIDInfo.Key,
			behaviorIDInfo.Default, behaviorIDInfo.Extra)

		// 如果behavior_id已正确设置为自增主键，则不需要进一步修复
		if behaviorIDInfo.Key == "PRI" && strings.Contains(behaviorIDInfo.Extra, "auto_increment") {
			log.Println("[数据库修复] 验证成功: behavior_id已正确设置为自增主键，无需修复")
			return nil
		}
	}

	// 只有当表结构确实存在问题时，才执行备份和重建操作
	log.Println("[数据库修复] 检测到behavior_id列配置不正确，需要修复...")

	// 检查表结构
	type TableInfo struct {
		TableSchema string `gorm:"column:TABLE_SCHEMA"`
		TableName   string `gorm:"column:TABLE_NAME"`
		Engine      string `gorm:"column:ENGINE"`
		TableRows   int    `gorm:"column:TABLE_ROWS"`
	}

	var tableInfo TableInfo
	err = db.Raw(`SELECT TABLE_SCHEMA, TABLE_NAME, ENGINE, TABLE_ROWS 
		FROM information_schema.tables 
		WHERE table_schema = ? AND table_name = ?`,
		cfg.DBName, "user_behaviors").Scan(&tableInfo).Error

	if err != nil {
		log.Printf("[数据库修复] 错误：获取user_behaviors表信息失败: %v", err)
	} else {
		log.Printf("[数据库修复] 表信息: 数据库=%s, 表名=%s, 引擎=%s, 行数=%d",
			tableInfo.TableSchema, tableInfo.TableName, tableInfo.Engine, tableInfo.TableRows)
	}

	// 备份并重建表，这是一个激进的方法，仅在结构确实有问题时使用
	log.Println("[数据库修复] 将删除并重建user_behaviors表...")

	// 1. 备份表内容（如果有必要）
	if tableInfo.TableRows > 0 {
		log.Printf("[数据库修复] 表中有 %d 行数据，将先备份", tableInfo.TableRows)

		// 创建备份表
		backupTableName := "user_behaviors_backup_" + time.Now().Format("20060102150405")
		createBackupSQL := fmt.Sprintf("CREATE TABLE %s LIKE user_behaviors", backupTableName)

		err = db.Exec(createBackupSQL).Error
		if err != nil {
			log.Printf("[数据库修复] 警告：创建备份表失败: %v", err)
		} else {
			// 复制数据
			copyDataSQL := fmt.Sprintf("INSERT INTO %s SELECT * FROM user_behaviors", backupTableName)
			err = db.Exec(copyDataSQL).Error
			if err != nil {
				log.Printf("[数据库修复] 警告：复制数据到备份表失败: %v", err)
			} else {
				log.Printf("[数据库修复] 成功创建备份表 %s", backupTableName)
			}
		}
	}

	// 2. 删除现有表
	dropTableSQL := "DROP TABLE IF EXISTS user_behaviors"
	err = db.Exec(dropTableSQL).Error
	if err != nil {
		log.Printf("[数据库修复] 错误：删除user_behaviors表失败: %v", err)
		return fmt.Errorf("删除user_behaviors表失败: %w", err)
	}
	log.Println("[数据库修复] 成功删除user_behaviors表")

	// 3. 通过GORM的AutoMigrate重新创建表
	err = db.AutoMigrate(&models.UserBehavior{})
	if err != nil {
		log.Printf("[数据库修复] 错误：通过AutoMigrate重建user_behaviors表失败: %v", err)
		return fmt.Errorf("重建user_behaviors表失败: %w", err)
	}
	log.Println("[数据库修复] 成功通过AutoMigrate重建user_behaviors表")

	// 4. 验证表结构
	var newColumns []struct {
		Field   string `gorm:"column:Field"`
		Type    string `gorm:"column:Type"`
		Null    string `gorm:"column:Null"`
		Key     string `gorm:"column:Key"`
		Default string `gorm:"column:Default"`
		Extra   string `gorm:"column:Extra"`
	}

	if err := db.Raw("SHOW COLUMNS FROM user_behaviors").Scan(&newColumns).Error; err != nil {
		log.Printf("[数据库修复] 错误：验证时获取表列信息失败: %v", err)
	} else {
		log.Println("[数据库修复] 重建后的表列信息:")
		for _, col := range newColumns {
			log.Printf("[数据库修复] 列名=%s, 类型=%s, 可空=%s, 键类型=%s, 默认值=%s, 额外=%s",
				col.Field, col.Type, col.Null, col.Key, col.Default, col.Extra)
		}
	}

	// 5. 特别检查behavior_id字段
	err = db.Raw("SHOW COLUMNS FROM user_behaviors WHERE Field = 'behavior_id'").Scan(&behaviorIDInfo).Error
	if err != nil {
		log.Printf("[数据库修复] 错误：验证时获取behavior_id列信息失败: %v", err)
	} else {
		log.Printf("[数据库修复] 验证behavior_id列: 类型=%s, 可空=%s, 键类型=%s, 默认值=%s, 额外=%s",
			behaviorIDInfo.Type, behaviorIDInfo.Null, behaviorIDInfo.Key,
			behaviorIDInfo.Default, behaviorIDInfo.Extra)

		if behaviorIDInfo.Key == "PRI" && strings.Contains(behaviorIDInfo.Extra, "auto_increment") {
			log.Println("[数据库修复] 验证成功: behavior_id已正确设置为自增主键")
		} else {
			log.Println("[数据库修复] 警告：验证失败，behavior_id列设置可能不正确")
		}
	}

	log.Println("[数据库修复] user_behaviors表结构修复完成")
	return nil
}

// fixDevicesTable 确保devices表中存在新增字段
func fixDevicesTable(db *gorm.DB) error {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Println("检查设备表结构，确保新字段存在...")
	}

	// 检查表是否存在
	var tableExists int
	row := db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ? AND table_name = ?",
		cfg.DBName, "devices").Row()
	row.Scan(&tableExists)

	if tableExists == 0 {
		if cfg.DebugLevel == "true" {
			log.Println("设备表不存在，将由自动迁移创建")
		}
		return nil
	}

	// 检查long_address字段是否存在
	var longAddressExists int
	row = db.Raw("SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = ? AND table_name = ? AND column_name = ?",
		cfg.DBName, "devices", "long_address").Row()
	row.Scan(&longAddressExists)

	// 如果long_address字段不存在，添加它
	if longAddressExists == 0 {
		if err := db.Exec("ALTER TABLE devices ADD COLUMN long_address VARCHAR(255)").Error; err != nil {
			return fmt.Errorf("添加long_address字段失败: %w", err)
		}
		if cfg.DebugLevel == "true" {
			log.Println("成功添加long_address字段到devices表")
		}
	}

	// 检查short_address字段是否存在
	var shortAddressExists int
	row = db.Raw("SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = ? AND table_name = ? AND column_name = ?",
		cfg.DBName, "devices", "short_address").Row()
	row.Scan(&shortAddressExists)

	// 如果short_address字段不存在，添加它
	if shortAddressExists == 0 {
		if err := db.Exec("ALTER TABLE devices ADD COLUMN short_address VARCHAR(255)").Error; err != nil {
			return fmt.Errorf("添加short_address字段失败: %w", err)
		}
		if cfg.DebugLevel == "true" {
			log.Println("成功添加short_address字段到devices表")
		}
	}

	// 检查ses_key字段是否存在
	var sesKeyExists int
	row = db.Raw("SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = ? AND table_name = ? AND column_name = ?",
		cfg.DBName, "devices", "ses_key").Row()
	row.Scan(&sesKeyExists)

	// 如果ses_key字段不存在，添加它
	if sesKeyExists == 0 {
		if err := db.Exec("ALTER TABLE devices ADD COLUMN ses_key VARCHAR(255)").Error; err != nil {
			return fmt.Errorf("添加ses_key字段失败: %w", err)
		}
		if cfg.DebugLevel == "true" {
			log.Println("成功添加ses_key字段到devices表")
		}
	}

	return nil
}
