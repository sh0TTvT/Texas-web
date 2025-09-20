// MySQL数据库连接管理
// 作用：建立MySQL连接池，提供数据库连接管理功能

package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DB 数据库连接实例
var DB *sql.DB

// InitMySQL 初始化MySQL连接池
func InitMySQL(databaseURL string) (*sql.DB, error) {
	// 打开数据库连接
	db, err := sql.Open("mysql", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 配置连接池参数（根据2GB内存限制配置）
	db.SetMaxOpenConns(20)        // 最大打开连接数
	db.SetMaxIdleConns(10)        // 最大空闲连接数
	db.SetConnMaxLifetime(time.Hour) // 连接最大生命周期

	// 测试数据库连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	DB = db
	return db, nil
}

// GetDB 获取数据库连接实例
func GetDB() *sql.DB {
	return DB
}

// Close 关闭数据库连接
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
} 