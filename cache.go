package cache

import (
	"context"
	"time"
)

// Cache 缓存接口
type Cache interface {
	// 基础操作
	Get(ctx context.Context, key string) (string, error)
	GetBytes(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) // 新增：仅当不存在时设置
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, keys ...string) (int64, error) // 改进：支持多个 key

	// 批量操作（性能优化）
	MGet(ctx context.Context, keys ...string) ([]interface{}, error)
	MSet(ctx context.Context, pairs map[string]interface{}) error
	MDelete(ctx context.Context, keys ...string) error

	// 数值操作
	Incr(ctx context.Context, key string) (int64, error)
	IncrBy(ctx context.Context, key string, value int64) (int64, error) // 新增
	Decr(ctx context.Context, key string) (int64, error)                // 新增
	DecrBy(ctx context.Context, key string, value int64) (int64, error) // 新增

	// TTL 操作
	TTL(ctx context.Context, key string) (time.Duration, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error // 新增：更新过期时间
	Persist(ctx context.Context, key string) error                          // 新增：移除过期时间

	// Hash 操作（常用场景）
	HGet(ctx context.Context, key, field string) (string, error)
	HSet(ctx context.Context, key, field string, value interface{}) error
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HDel(ctx context.Context, key string, fields ...string) error

	// 列表操作（可选，看需求）
	LPush(ctx context.Context, key string, values ...interface{}) error
	RPush(ctx context.Context, key string, values ...interface{}) error
	LPop(ctx context.Context, key string) (string, error)
	RPop(ctx context.Context, key string) (string, error)
	LLen(ctx context.Context, key string) (int64, error)

	// 集合操作（可选）
	SAdd(ctx context.Context, key string, members ...interface{}) error
	SMembers(ctx context.Context, key string) ([]string, error)
	SRem(ctx context.Context, key string, members ...interface{}) error

	// 管理操作
	Ping(ctx context.Context) error                             // 新增：健康检查
	FlushDB(ctx context.Context) error                          // 改名：原 Clear，更明确
	Keys(ctx context.Context, pattern string) ([]string, error) // 新增：查找键（谨慎使用）
	Close() error
}

// StringCache 简化版接口（如果只需要字符串缓存）
type StringCache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Close() error
}
