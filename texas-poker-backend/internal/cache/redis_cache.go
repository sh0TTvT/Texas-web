// Redis缓存层优化
// 作用：提供高性能的缓存服务，优化数据访问速度

package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

// 缓存键前缀
const (
	UserPrefix     = "user:"
	RoomPrefix     = "room:"
	GamePrefix     = "game:"
	SessionPrefix  = "session:"
	OnlinePrefix   = "online:"
	StatsPrefix    = "stats:"
	RankingPrefix  = "ranking:"
)

// 缓存过期时间
const (
	UserCacheTTL     = 30 * time.Minute
	RoomCacheTTL     = 10 * time.Minute
	GameCacheTTL     = 5 * time.Minute
	SessionCacheTTL  = 24 * time.Hour
	OnlineCacheTTL   = 5 * time.Minute
	StatsCacheTTL    = 5 * time.Minute
	RankingCacheTTL  = 15 * time.Minute
)

// NewRedisCache 创建Redis缓存实例
func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{
		client: client,
		ctx:    context.Background(),
	}
}

// Set 设置缓存值
func (r *RedisCache) Set(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal value: %w", err)
	}

	err = r.client.Set(r.ctx, key, data, ttl).Err()
	if err != nil {
		return fmt.Errorf("set cache: %w", err)
	}

	return nil
}

// Get 获取缓存值
func (r *RedisCache) Get(key string, dest interface{}) error {
	data, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("cache miss")
		}
		return fmt.Errorf("get cache: %w", err)
	}

	err = json.Unmarshal([]byte(data), dest)
	if err != nil {
		return fmt.Errorf("unmarshal value: %w", err)
	}

	return nil
}

// Del 删除缓存
func (r *RedisCache) Del(keys ...string) error {
	err := r.client.Del(r.ctx, keys...).Err()
	if err != nil {
		return fmt.Errorf("delete cache: %w", err)
	}
	return nil
}

// Exists 检查缓存是否存在
func (r *RedisCache) Exists(key string) (bool, error) {
	count, err := r.client.Exists(r.ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("check cache existence: %w", err)
	}
	return count > 0, nil
}

// SetEx 设置带过期时间的缓存
func (r *RedisCache) SetEx(key string, value interface{}, seconds int) error {
	return r.Set(key, value, time.Duration(seconds)*time.Second)
}

// Incr 递增计数器
func (r *RedisCache) Incr(key string) (int64, error) {
	val, err := r.client.Incr(r.ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("increment counter: %w", err)
	}
	return val, nil
}

// IncrBy 按指定值递增
func (r *RedisCache) IncrBy(key string, value int64) (int64, error) {
	val, err := r.client.IncrBy(r.ctx, key, value).Result()
	if err != nil {
		return 0, fmt.Errorf("increment by value: %w", err)
	}
	return val, nil
}

// HSet 设置哈希字段
func (r *RedisCache) HSet(key, field string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal hash value: %w", err)
	}

	err = r.client.HSet(r.ctx, key, field, data).Err()
	if err != nil {
		return fmt.Errorf("set hash field: %w", err)
	}

	return nil
}

// HGet 获取哈希字段
func (r *RedisCache) HGet(key, field string, dest interface{}) error {
	data, err := r.client.HGet(r.ctx, key, field).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("hash field not found")
		}
		return fmt.Errorf("get hash field: %w", err)
	}

	err = json.Unmarshal([]byte(data), dest)
	if err != nil {
		return fmt.Errorf("unmarshal hash value: %w", err)
	}

	return nil
}

// HDel 删除哈希字段
func (r *RedisCache) HDel(key string, fields ...string) error {
	err := r.client.HDel(r.ctx, key, fields...).Err()
	if err != nil {
		return fmt.Errorf("delete hash fields: %w", err)
	}
	return nil
}

// HGetAll 获取所有哈希字段
func (r *RedisCache) HGetAll(key string) (map[string]string, error) {
	data, err := r.client.HGetAll(r.ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("get all hash fields: %w", err)
	}
	return data, nil
}

// SAdd 添加集合成员
func (r *RedisCache) SAdd(key string, members ...interface{}) error {
	err := r.client.SAdd(r.ctx, key, members...).Err()
	if err != nil {
		return fmt.Errorf("add set members: %w", err)
	}
	return nil
}

// SRem 移除集合成员
func (r *RedisCache) SRem(key string, members ...interface{}) error {
	err := r.client.SRem(r.ctx, key, members...).Err()
	if err != nil {
		return fmt.Errorf("remove set members: %w", err)
	}
	return nil
}

// SMembers 获取集合所有成员
func (r *RedisCache) SMembers(key string) ([]string, error) {
	members, err := r.client.SMembers(r.ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("get set members: %w", err)
	}
	return members, nil
}

// SIsMember 检查是否是集合成员
func (r *RedisCache) SIsMember(key string, member interface{}) (bool, error) {
	isMember, err := r.client.SIsMember(r.ctx, key, member).Result()
	if err != nil {
		return false, fmt.Errorf("check set membership: %w", err)
	}
	return isMember, nil
}

// ZAdd 添加有序集合成员
func (r *RedisCache) ZAdd(key string, score float64, member interface{}) error {
	z := &redis.Z{Score: score, Member: member}
	err := r.client.ZAdd(r.ctx, key, z).Err()
	if err != nil {
		return fmt.Errorf("add sorted set member: %w", err)
	}
	return nil
}

// ZRem 移除有序集合成员
func (r *RedisCache) ZRem(key string, members ...interface{}) error {
	err := r.client.ZRem(r.ctx, key, members...).Err()
	if err != nil {
		return fmt.Errorf("remove sorted set members: %w", err)
	}
	return nil
}

// ZRange 获取有序集合范围内的成员
func (r *RedisCache) ZRange(key string, start, stop int64) ([]string, error) {
	members, err := r.client.ZRange(r.ctx, key, start, stop).Result()
	if err != nil {
		return nil, fmt.Errorf("get sorted set range: %w", err)
	}
	return members, nil
}

// ZRevRange 获取有序集合范围内的成员（倒序）
func (r *RedisCache) ZRevRange(key string, start, stop int64) ([]string, error) {
	members, err := r.client.ZRevRange(r.ctx, key, start, stop).Result()
	if err != nil {
		return nil, fmt.Errorf("get sorted set reverse range: %w", err)
	}
	return members, nil
}

// ZRangeWithScores 获取有序集合范围内的成员和分数
func (r *RedisCache) ZRangeWithScores(key string, start, stop int64) ([]redis.Z, error) {
	members, err := r.client.ZRangeWithScores(r.ctx, key, start, stop).Result()
	if err != nil {
		return nil, fmt.Errorf("get sorted set range with scores: %w", err)
	}
	return members, nil
}

// ExpireAt 设置过期时间
func (r *RedisCache) ExpireAt(key string, tm time.Time) error {
	err := r.client.ExpireAt(r.ctx, key, tm).Err()
	if err != nil {
		return fmt.Errorf("set expiration: %w", err)
	}
	return nil
}

// TTL 获取剩余生存时间
func (r *RedisCache) TTL(key string) (time.Duration, error) {
	ttl, err := r.client.TTL(r.ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("get TTL: %w", err)
	}
	return ttl, nil
}

// Keys 获取匹配模式的键
func (r *RedisCache) Keys(pattern string) ([]string, error) {
	keys, err := r.client.Keys(r.ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("get keys: %w", err)
	}
	return keys, nil
}

// Pipeline 创建管道
func (r *RedisCache) Pipeline() redis.Pipeliner {
	return r.client.Pipeline()
}

// TxPipeline 创建事务管道
func (r *RedisCache) TxPipeline() redis.Pipeliner {
	return r.client.TxPipeline()
}

// 用户相关缓存方法
func (r *RedisCache) SetUser(userID string, user interface{}) error {
	key := UserPrefix + userID
	return r.Set(key, user, UserCacheTTL)
}

func (r *RedisCache) GetUser(userID string, dest interface{}) error {
	key := UserPrefix + userID
	return r.Get(key, dest)
}

func (r *RedisCache) DelUser(userID string) error {
	key := UserPrefix + userID
	return r.Del(key)
}

// 房间相关缓存方法
func (r *RedisCache) SetRoom(roomID string, room interface{}) error {
	key := RoomPrefix + roomID
	return r.Set(key, room, RoomCacheTTL)
}

func (r *RedisCache) GetRoom(roomID string, dest interface{}) error {
	key := RoomPrefix + roomID
	return r.Get(key, dest)
}

func (r *RedisCache) DelRoom(roomID string) error {
	key := RoomPrefix + roomID
	return r.Del(key)
}

// 游戏相关缓存方法
func (r *RedisCache) SetGame(gameID string, game interface{}) error {
	key := GamePrefix + gameID
	return r.Set(key, game, GameCacheTTL)
}

func (r *RedisCache) GetGame(gameID string, dest interface{}) error {
	key := GamePrefix + gameID
	return r.Get(key, dest)
}

func (r *RedisCache) DelGame(gameID string) error {
	key := GamePrefix + gameID
	return r.Del(key)
}

// 会话相关缓存方法
func (r *RedisCache) SetSession(sessionID string, session interface{}) error {
	key := SessionPrefix + sessionID
	return r.Set(key, session, SessionCacheTTL)
}

func (r *RedisCache) GetSession(sessionID string, dest interface{}) error {
	key := SessionPrefix + sessionID
	return r.Get(key, dest)
}

func (r *RedisCache) DelSession(sessionID string) error {
	key := SessionPrefix + sessionID
	return r.Del(key)
}

// 在线用户管理
func (r *RedisCache) AddOnlineUser(userID string) error {
	key := OnlinePrefix + "users"
	return r.SAdd(key, userID)
}

func (r *RedisCache) RemoveOnlineUser(userID string) error {
	key := OnlinePrefix + "users"
	return r.SRem(key, userID)
}

func (r *RedisCache) GetOnlineUsers() ([]string, error) {
	key := OnlinePrefix + "users"
	return r.SMembers(key)
}

func (r *RedisCache) IsUserOnline(userID string) (bool, error) {
	key := OnlinePrefix + "users"
	return r.SIsMember(key, userID)
}

// 排行榜管理
func (r *RedisCache) UpdateRanking(userID string, score float64) error {
	key := RankingPrefix + "chips"
	return r.ZAdd(key, score, userID)
}

func (r *RedisCache) GetTopRanking(limit int64) ([]redis.Z, error) {
	key := RankingPrefix + "chips"
	return r.ZRangeWithScores(key, 0, limit-1)
}

func (r *RedisCache) GetUserRank(userID string) (int64, error) {
	key := RankingPrefix + "chips"
	rank, err := r.client.ZRevRank(r.ctx, key, userID).Result()
	if err != nil {
		return 0, fmt.Errorf("get user rank: %w", err)
	}
	return rank + 1, nil // 排名从1开始
}

// 统计数据缓存
func (r *RedisCache) SetStats(statsType string, stats interface{}) error {
	key := StatsPrefix + statsType
	return r.Set(key, stats, StatsCacheTTL)
}

func (r *RedisCache) GetStats(statsType string, dest interface{}) error {
	key := StatsPrefix + statsType
	return r.Get(key, dest)
}

// 批量操作
func (r *RedisCache) MSet(pairs map[string]interface{}) error {
	pipe := r.Pipeline()
	
	for key, value := range pairs {
		data, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("marshal value for key %s: %w", key, err)
		}
		pipe.Set(r.ctx, key, data, 0)
	}
	
	_, err := pipe.Exec(r.ctx)
	if err != nil {
		return fmt.Errorf("execute pipeline: %w", err)
	}
	
	return nil
}

func (r *RedisCache) MGet(keys []string) ([]interface{}, error) {
	values, err := r.client.MGet(r.ctx, keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("multi get: %w", err)
	}
	return values, nil
}

// 清理过期缓存
func (r *RedisCache) CleanupExpired() error {
	// Redis会自动清理过期键，这里主要是日志记录
	return nil
}

// 获取缓存统计信息
func (r *RedisCache) GetCacheStats() (map[string]interface{}, error) {
	info, err := r.client.Info(r.ctx, "memory").Result()
	if err != nil {
		return nil, fmt.Errorf("get cache stats: %w", err)
	}
	
	stats := map[string]interface{}{
		"memory_info": info,
	}
	
	return stats, nil
} 