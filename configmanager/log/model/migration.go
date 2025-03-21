package model

import (
	"fmt"
	"log"
	"strings"

	"gin-server/config"

	"gorm.io/gorm"
)

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

// 用户行为表定义
var userBehaviorsTable = TableDefinition{
	Name: "user_behaviors",
	Columns: []Column{
		{Name: "behaviorID", Type: "INT", PrimaryKey: true, Nullable: false, Comment: "AUTO_INCREMENT PRIMARY KEY"},
		{Name: "userID", Type: "INT", Nullable: false, Comment: "用户ID"},
		{Name: "behaviorTime", Type: "DATETIME(3)", Nullable: false, Comment: "行为开始时间"},
		{Name: "behaviorType", Type: "INT", Nullable: false, Comment: "行为类型，1:发送，2:接收"},
		{Name: "dataType", Type: "INT", Nullable: false, Comment: "数据类型，1:文件，2:消息"},
		{Name: "dataSize", Type: "BIGINT", Nullable: false, Comment: "数据大小"},
		{Name: "created_at", Type: "TIMESTAMP(3)", Nullable: true, Comment: "创建时间"},
	},
	Indexes: []string{
		"idx_userid (userID)",
		"idx_behaviortime (behaviorTime)",
		"idx_behaviortype (behaviorType)",
	},
	ForeignKey: &ForeignKey{
		Column:    "userID",
		RefTable:  "users",
		RefColumn: "userID",
	},
}

// 事件表定义
var eventsTable = TableDefinition{
	Name: "events",
	Columns: []Column{
		{Name: "eventId", Type: "BIGINT", PrimaryKey: true, Nullable: false, Comment: "事件ID"},
		{Name: "deviceId", Type: "VARCHAR(12)", Nullable: false, Comment: "设备ID"},
		{Name: "eventTime", Type: "DATETIME", Nullable: false, Comment: "事件发生时间"},
		{Name: "eventType", Type: "INT", Nullable: false, Comment: "事件类型，1:安全事件，2:故障事件"},
		{Name: "eventCode", Type: "VARCHAR(20)", Nullable: false, Comment: "事件代码"},
		{Name: "eventDesc", Type: "VARCHAR(255)", Nullable: false, Comment: "事件描述"},
		{Name: "createdAt", Type: "TIMESTAMP", Nullable: true, Comment: "创建时间"},
	},
	Indexes: []string{},
}

// MigrateDatabase 检查并迁移数据库表
func MigrateDatabase(db *gorm.DB) error {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Println("开始检查数据库表结构...")
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
		log.Println("数据库表结构检查完成！")
	}
	return nil
}

// checkAndMigrateTable 检查并迁移指定表
func checkAndMigrateTable(db *gorm.DB, tableDef TableDefinition) error {
	// 检查表是否存在
	if !tableExists(db, tableDef.Name) {
		return createTable(db, tableDef)
	}

	// 表存在，检查字段
	return checkAndFixColumns(db, tableDef)
}

// tableExists 检查表是否存在
func tableExists(db *gorm.DB, tableName string) bool {
	var count int64
	db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?", tableName).Count(&count)
	return count > 0
}

// createTable 创建表
func createTable(db *gorm.DB, tableDef TableDefinition) error {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Printf("表 %s 不存在，正在创建...\n", tableDef.Name)
	}

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

	// 执行SQL
	if err := db.Exec(sql).Error; err != nil {
		return fmt.Errorf("创建表 %s 失败: %w", tableDef.Name, err)
	}

	if cfg.DebugLevel == "true" {
		log.Printf("表 %s 创建成功\n", tableDef.Name)
	}
	return nil
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

	// 检查是否缺少索引
	var tableIndexes []struct {
		KeyName    string `gorm:"column:Key_name"`
		ColumnName string `gorm:"column:Column_name"`
	}
	if err := db.Raw("SHOW INDEX FROM " + tableDef.Name).Scan(&tableIndexes).Error; err != nil {
		return fmt.Errorf("获取表 %s 的索引信息失败: %w", tableDef.Name, err)
	}

	existingIndexes := make(map[string]bool)
	for _, idx := range tableIndexes {
		existingIndexes[idx.KeyName] = true
	}

	// 添加缺失的索引
	for _, idxStr := range tableDef.Indexes {
		parts := strings.SplitN(idxStr, " ", 2)
		if len(parts) != 2 {
			continue
		}

		idxName := parts[0]
		idxCols := parts[1]

		if !existingIndexes[idxName] {
			sql := fmt.Sprintf("ALTER TABLE %s ADD INDEX %s %s", tableDef.Name, idxName, idxCols)
			if err := db.Exec(sql).Error; err != nil {
				return fmt.Errorf("向表 %s 添加索引 %s 失败: %w", tableDef.Name, idxName, err)
			}

			if cfg.DebugLevel == "true" {
				log.Printf("已向表 %s 添加索引 %s\n", tableDef.Name, idxName)
			}
		}
	}

	if cfg.DebugLevel == "true" {
		log.Printf("表 %s 字段结构检查完成\n", tableDef.Name)
	}
	return nil
}
