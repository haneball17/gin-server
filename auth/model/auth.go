package model

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"gin-server/regist/model"
)

// AuthRecord 认证记录结构体
type AuthRecord struct {
	ID               int       `json:"id"`
	Username         string    `json:"username"`
	Pass             string    `json:"pass"`
	Reply            string    `json:"reply"`
	AuthDate         time.Time `json:"authdate"`
	Class            string    `json:"class"`
	CalledStationID  string    `json:"calledstationid"`
	CallingStationID string    `json:"callingstationid"`
}

// AuthRecordQuery 查询条件结构体
type AuthRecordQuery struct {
	Username         string `form:"username"`
	Reply            string `form:"reply"`
	StartDate        string `form:"start_date"`
	EndDate          string `form:"end_date"`
	Page             int    `form:"page,default=1"`
	PageSize         int    `form:"page_size,default=10"`
	CalledStationID  string `form:"calledstationid"`
	CallingStationID string `form:"callingstationid"`
}

// 全局变量，存储radius数据库连接
var radiusDB *sql.DB

// InitRadiusDB 初始化Radius数据库连接
func InitRadiusDB() error {
	config := model.LoadConfig()
	if config.DebugLevel == "true" {
		log.Println("开始初始化Radius数据库连接")
	}

	// 构建DSN，添加parseTime=true参数
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/radius?parseTime=true",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort)

	var err error
	radiusDB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("连接Radius数据库失败: %v\n", err)
		return fmt.Errorf("连接Radius数据库失败: %w", err)
	}

	// 测试连接
	if err = radiusDB.Ping(); err != nil {
		log.Printf("Radius数据库Ping失败: %v\n", err)
		return fmt.Errorf("Radius数据库Ping失败: %w", err)
	}

	if config.DebugLevel == "true" {
		log.Println("Radius数据库连接成功")
	}
	return nil
}

// GetRadiusDB 获取Radius数据库连接
func GetRadiusDB() *sql.DB {
	return radiusDB
}

// GetAuthRecords 获取认证记录
func GetAuthRecords(query AuthRecordQuery) ([]AuthRecord, int, error) {
	config := model.LoadConfig()
	if config.DebugLevel == "true" {
		log.Printf("开始查询认证记录，查询条件: %+v\n", query)
	}

	// 构建查询条件
	whereClause := "WHERE 1=1"
	args := []interface{}{}

	if query.Username != "" {
		whereClause += " AND username LIKE ?"
		args = append(args, "%"+query.Username+"%")
	}

	if query.Reply != "" {
		whereClause += " AND reply = ?"
		args = append(args, query.Reply)
	}

	if query.StartDate != "" {
		whereClause += " AND authdate >= ?"
		args = append(args, query.StartDate)
	}

	if query.EndDate != "" {
		whereClause += " AND authdate <= ?"
		args = append(args, query.EndDate+" 23:59:59")
	}

	if query.CalledStationID != "" {
		whereClause += " AND calledstationid LIKE ?"
		args = append(args, "%"+query.CalledStationID+"%")
	}

	if query.CallingStationID != "" {
		whereClause += " AND callingstationid LIKE ?"
		args = append(args, "%"+query.CallingStationID+"%")
	}

	// 查询总记录数
	countSQL := fmt.Sprintf("SELECT COUNT(*) FROM radpostauth %s", whereClause)
	var total int
	err := radiusDB.QueryRow(countSQL, args...).Scan(&total)
	if err != nil {
		log.Printf("查询认证记录总数失败: %v\n", err)
		return nil, 0, fmt.Errorf("查询认证记录总数失败: %w", err)
	}

	// 计算分页
	offset := (query.Page - 1) * query.PageSize

	// 查询记录
	querySQL := fmt.Sprintf("SELECT id, username, pass, reply, authdate, class, calledstationid, callingstationid FROM radpostauth %s ORDER BY authdate DESC LIMIT ? OFFSET ?", whereClause)
	args = append(args, query.PageSize, offset)

	if config.DebugLevel == "true" {
		log.Printf("执行SQL: %s, 参数: %v\n", querySQL, args)
	}

	rows, err := radiusDB.Query(querySQL, args...)
	if err != nil {
		log.Printf("查询认证记录失败: %v\n", err)
		return nil, 0, fmt.Errorf("查询认证记录失败: %w", err)
	}
	defer rows.Close()

	var records []AuthRecord
	for rows.Next() {
		var record AuthRecord
		var class sql.NullString
		err := rows.Scan(
			&record.ID,
			&record.Username,
			&record.Pass,
			&record.Reply,
			&record.AuthDate,
			&class,
			&record.CalledStationID,
			&record.CallingStationID,
		)
		if err != nil {
			log.Printf("扫描认证记录失败: %v\n", err)
			return nil, 0, fmt.Errorf("扫描认证记录失败: %w", err)
		}

		if class.Valid {
			record.Class = class.String
		}

		records = append(records, record)
	}

	if config.DebugLevel == "true" {
		log.Printf("成功获取 %d 条认证记录，总记录数: %d\n", len(records), total)
	}

	return records, total, nil
}
