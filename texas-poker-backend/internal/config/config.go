// 配置管理
// 作用：加载和管理应用程序配置，包括数据库连接、Redis连接、JWT密钥等

package config

import (
	"os"
)

// Config 应用程序配置结构
type Config struct {
	Port        string // 服务器端口
	GinMode     string // Gin运行模式
	DatabaseURL string // MySQL数据库连接URL
	RedisURL    string // Redis连接URL
	JWTSecret   string // JWT签名密钥
	AdminSecret string // 管理员特殊密钥
}

// Load 加载配置
func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		GinMode:     getEnv("GIN_MODE", "debug"),
		DatabaseURL: getEnv("DATABASE_URL", "root:password@tcp(localhost:3306)/texas_poker?charset=utf8mb4&parseTime=True&loc=Local"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379/0"),
		JWTSecret:   getEnv("JWT_SECRET", "texas-poker-jwt-secret-key-2024"),
		AdminSecret: getEnv("ADMIN_SECRET", "texas-poker-admin-secret-2024"),
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 