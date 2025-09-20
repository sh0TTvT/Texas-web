// Redis连接管理
// 作用：建立Redis连接池，提供缓存服务

package database

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisClient Redis客户端实例
var RedisClient *redis.Client

// InitRedis 初始化Redis连接
func InitRedis(redisURL string) *redis.Client {
	// 解析Redis URL并创建客户端
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		// 如果解析失败，使用默认配置
		opt = &redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		}
	}

	// 配置连接池参数（根据2GB内存限制配置）
	opt.PoolSize = 50                // 最大连接数
	opt.MinIdleConns = 10            // 最小空闲连接数
	opt.MaxConnAge = time.Hour       // 连接最大生命周期
	opt.PoolTimeout = 30 * time.Second // 连接池超时
	opt.IdleTimeout = 5 * time.Minute  // 空闲连接超时

	// 创建Redis客户端
	client := redis.NewClient(opt)

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("Failed to connect to Redis: %v\n", err)
	}

	RedisClient = client
	return client
}

// GetRedis 获取Redis客户端实例
func GetRedis() *redis.Client {
	return RedisClient
}

// CloseRedis 关闭Redis连接
func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
} 