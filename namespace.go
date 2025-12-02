package cache

import (
	"context"
	"strings"
	"time"
)

// Namespace 缓存命名空间包装器
// 为缓存键添加统一的前缀，避免不同业务模块的键冲突
type Namespace struct {
	cache  Cache
	prefix string
	sep    string // 分隔符，默认为 ":"
}

// NamespaceOption 命名空间配置选项
type NamespaceOption func(*Namespace)

// WithSeparator 设置键分隔符
func WithSeparator(sep string) NamespaceOption {
	return func(n *Namespace) {
		n.sep = sep
	}
}

// NewNamespace 创建带命名空间的缓存
// prefix: 命名空间前缀，如 "user", "session", "product"
// opts: 可选配置
func NewNamespace(cache Cache, prefix string, opts ...NamespaceOption) *Namespace {
	ns := &Namespace{
		cache:  cache,
		prefix: strings.TrimSpace(prefix),
		sep:    ":", // 默认分隔符
	}

	// 应用选项
	for _, opt := range opts {
		opt(ns)
	}

	return ns
}

// key 为原始键添加命名空间前缀
func (n *Namespace) key(key string) string {
	if n.prefix == "" {
		return key
	}
	return n.prefix + n.sep + key
}

// keys 批量添加命名空间前缀
func (n *Namespace) keys(keys ...string) []string {
	if n.prefix == "" {
		return keys
	}

	result := make([]string, len(keys))
	for i, key := range keys {
		result[i] = n.key(key)
	}
	return result
}

// stripPrefix 从带前缀的键中移除命名空间前缀
func (n *Namespace) stripPrefix(key string) string {
	if n.prefix == "" {
		return key
	}

	prefix := n.prefix + n.sep
	if strings.HasPrefix(key, prefix) {
		return strings.TrimPrefix(key, prefix)
	}
	return key
}

// ============================================
// 基础操作
// ============================================

func (n *Namespace) Get(ctx context.Context, key string) (string, error) {
	return n.cache.Get(ctx, n.key(key))
}

func (n *Namespace) GetBytes(ctx context.Context, key string) ([]byte, error) {
	return n.cache.GetBytes(ctx, n.key(key))
}

func (n *Namespace) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return n.cache.Set(ctx, n.key(key), value, expiration)
}

func (n *Namespace) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return n.cache.SetNX(ctx, n.key(key), value, expiration)
}

func (n *Namespace) Delete(ctx context.Context, key string) error {
	return n.cache.Delete(ctx, n.key(key))
}

func (n *Namespace) Exists(ctx context.Context, keys ...string) (int64, error) {
	return n.cache.Exists(ctx, n.keys(keys...)...)
}

// ============================================
// 批量操作
// ============================================

func (n *Namespace) MGet(ctx context.Context, keys ...string) ([]interface{}, error) {
	return n.cache.MGet(ctx, n.keys(keys...)...)
}

func (n *Namespace) MSet(ctx context.Context, pairs map[string]interface{}) error {
	// 为所有键添加命名空间前缀
	namespacedPairs := make(map[string]interface{}, len(pairs))
	for k, v := range pairs {
		namespacedPairs[n.key(k)] = v
	}
	return n.cache.MSet(ctx, namespacedPairs)
}

func (n *Namespace) MDelete(ctx context.Context, keys ...string) error {
	return n.cache.MDelete(ctx, n.keys(keys...)...)
}

// ============================================
// 数值操作
// ============================================

func (n *Namespace) Incr(ctx context.Context, key string) (int64, error) {
	return n.cache.Incr(ctx, n.key(key))
}

func (n *Namespace) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return n.cache.IncrBy(ctx, n.key(key), value)
}

func (n *Namespace) Decr(ctx context.Context, key string) (int64, error) {
	return n.cache.Decr(ctx, n.key(key))
}

func (n *Namespace) DecrBy(ctx context.Context, key string, value int64) (int64, error) {
	return n.cache.DecrBy(ctx, n.key(key), value)
}

// ============================================
// TTL 操作
// ============================================

func (n *Namespace) TTL(ctx context.Context, key string) (time.Duration, error) {
	return n.cache.TTL(ctx, n.key(key))
}

func (n *Namespace) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return n.cache.Expire(ctx, n.key(key), expiration)
}

func (n *Namespace) Persist(ctx context.Context, key string) error {
	return n.cache.Persist(ctx, n.key(key))
}

// ============================================
// Hash 操作
// ============================================

func (n *Namespace) HGet(ctx context.Context, key, field string) (string, error) {
	return n.cache.HGet(ctx, n.key(key), field)
}

func (n *Namespace) HSet(ctx context.Context, key, field string, value interface{}) error {
	return n.cache.HSet(ctx, n.key(key), field, value)
}

func (n *Namespace) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return n.cache.HGetAll(ctx, n.key(key))
}

func (n *Namespace) HDel(ctx context.Context, key string, fields ...string) error {
	return n.cache.HDel(ctx, n.key(key), fields...)
}

// ============================================
// 列表操作
// ============================================

func (n *Namespace) LPush(ctx context.Context, key string, values ...interface{}) error {
	return n.cache.LPush(ctx, n.key(key), values...)
}

func (n *Namespace) RPush(ctx context.Context, key string, values ...interface{}) error {
	return n.cache.RPush(ctx, n.key(key), values...)
}

func (n *Namespace) LPop(ctx context.Context, key string) (string, error) {
	return n.cache.LPop(ctx, n.key(key))
}

func (n *Namespace) RPop(ctx context.Context, key string) (string, error) {
	return n.cache.RPop(ctx, n.key(key))
}

func (n *Namespace) LLen(ctx context.Context, key string) (int64, error) {
	return n.cache.LLen(ctx, n.key(key))
}

// ============================================
// 集合操作
// ============================================

func (n *Namespace) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return n.cache.SAdd(ctx, n.key(key), members...)
}

func (n *Namespace) SMembers(ctx context.Context, key string) ([]string, error) {
	return n.cache.SMembers(ctx, n.key(key))
}

func (n *Namespace) SRem(ctx context.Context, key string, members ...interface{}) error {
	return n.cache.SRem(ctx, n.key(key), members...)
}

// ============================================
// Pub/Sub 操作
// ============================================

func (n *Namespace) Publish(ctx context.Context, channel string, message string) error {
	return n.cache.Publish(ctx, n.key(channel), message)
}

func (n *Namespace) Subscribe(ctx context.Context, channels ...string) (PubSub, error) {
	return n.cache.Subscribe(ctx, n.keys(channels...)...)
}

// ============================================
// 管理操作
// ============================================

func (n *Namespace) Ping(ctx context.Context) error {
	return n.cache.Ping(ctx)
}

// FlushDB 清空整个数据库（危险操作，不建议在命名空间中使用）
// 注意：这会清空整个 Redis DB，不仅仅是当前命名空间
func (n *Namespace) FlushDB(ctx context.Context) error {
	return n.cache.FlushDB(ctx)
}

// FlushNamespace 仅清空当前命名空间的所有键
// 这是一个相对安全的清空操作，只影响当前命名空间
func (n *Namespace) FlushNamespace(ctx context.Context) error {
	if n.prefix == "" {
		// 如果没有前缀，等同于 FlushDB
		return n.cache.FlushDB(ctx)
	}

	// 查找所有匹配的键
	pattern := n.prefix + n.sep + "*"
	keys, err := n.cache.Keys(ctx, pattern)
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	// 批量删除
	return n.cache.MDelete(ctx, keys...)
}

// Keys 查找当前命名空间下匹配模式的键
// pattern: 键的模式，如 "user*"
// 返回的键已经移除了命名空间前缀
func (n *Namespace) Keys(ctx context.Context, pattern string) ([]string, error) {
	// 为模式添加命名空间前缀
	namespacedPattern := n.key(pattern)

	keys, err := n.cache.Keys(ctx, namespacedPattern)
	if err != nil {
		return nil, err
	}

	// 移除命名空间前缀
	result := make([]string, len(keys))
	for i, key := range keys {
		result[i] = n.stripPrefix(key)
	}

	return result, nil
}

func (n *Namespace) Close() error {
	return n.cache.Close()
}

// ============================================
// 辅助方法
// ============================================

// Prefix 返回当前命名空间前缀
func (n *Namespace) Prefix() string {
	return n.prefix
}

// Separator 返回当前使用的分隔符
func (n *Namespace) Separator() string {
	return n.sep
}

// Unwrap 返回底层的 Cache 实例
func (n *Namespace) Unwrap() Cache {
	return n.cache
}

// SubNamespace 创建子命名空间
// 例如：parent namespace 是 "app"，创建子命名空间 "user" 后，前缀变为 "app:user"
func (n *Namespace) SubNamespace(subPrefix string, opts ...NamespaceOption) *Namespace {
	newPrefix := n.prefix
	if newPrefix != "" {
		newPrefix += n.sep + subPrefix
	} else {
		newPrefix = subPrefix
	}

	sub := &Namespace{
		cache:  n.cache,
		prefix: newPrefix,
		sep:    n.sep,
	}

	// 应用选项
	for _, opt := range opts {
		opt(sub)
	}

	return sub
}
