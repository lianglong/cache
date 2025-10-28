# Cache Namespace 使用指南

## 简介

Namespace（命名空间）为缓存键提供前缀隔离功能，避免不同业务模块的键冲突。

## 基础使用

### 1. 创建命名空间

```go
package main

import (
	"context"
	"time"

	"github.com/lianglong/cache"
	_ "github.com/lianglong/cache-redis"
)

func main() {
	ctx := context.Background()

	// 创建基础缓存客户端
	cacheClient, err := cache.New("redis", cache.Config{
		Addr:     "127.0.0.1:6379",
		Password: "123456",
		DB:       0,
	})
	if err != nil {
		panic(err)
	}
	defer cacheClient.Close()

	// 创建不同业务模块的命名空间
	userCache := cache.NewNamespace(cacheClient, "user")
	sessionCache := cache.NewNamespace(cacheClient, "session")
	productCache := cache.NewNamespace(cacheClient, "product")
}
```

### 2. 基础操作

```go
// 设置用户数据
userCache.Set(ctx, "123", "Alice", time.Hour)
// 实际存储的键: "user:123"

// 获取用户数据
value, err := userCache.Get(ctx, "123")
// 查询的键: "user:123"

// 删除用户数据
userCache.Delete(ctx, "123")
```

### 3. 键隔离演示

```go
// 不同命名空间的相同键名是隔离的
userCache.Set(ctx, "123", "User Alice", time.Hour)       // 实际键: user:123
sessionCache.Set(ctx, "123", "Session Data", time.Hour)  // 实际键: session:123

// 获取到的是不同的值
userData, _ := userCache.Get(ctx, "123")       // "User Alice"
sessionData, _ := sessionCache.Get(ctx, "123") // "Session Data"
```

## 高级功能

### 1. 批量操作

```go
// 批量设置
productCache.MSet(ctx, map[string]interface{}{
	"001": "Product A",  // 实际键: product:001
	"002": "Product B",  // 实际键: product:002
	"003": "Product C",  // 实际键: product:003
})

// 批量获取
values, _ := productCache.MGet(ctx, "001", "002", "003")

// 批量删除
productCache.MDelete(ctx, "001", "002", "003")
```

### 2. 计数器操作

```go
statsCache := cache.NewNamespace(cacheClient, "stats")

// 递增
count, _ := statsCache.Incr(ctx, "page_views")           // +1
count, _ = statsCache.IncrBy(ctx, "page_views", 10)      // +10

// 递减
count, _ = statsCache.Decr(ctx, "page_views")            // -1
count, _ = statsCache.DecrBy(ctx, "page_views", 5)       // -5
```

### 3. Hash 操作

```go
userCache := cache.NewNamespace(cacheClient, "user")

// 设置用户资料的多个字段
userCache.HSet(ctx, "1:profile", "name", "Alice")
userCache.HSet(ctx, "1:profile", "age", 30)
userCache.HSet(ctx, "1:profile", "email", "alice@example.com")
// 实际键: user:1:profile

// 获取单个字段
name, _ := userCache.HGet(ctx, "1:profile", "name")

// 获取所有字段
profile, _ := userCache.HGetAll(ctx, "1:profile")
// profile = {"name": "Alice", "age": "30", "email": "alice@example.com"}

// 删除字段
userCache.HDel(ctx, "1:profile", "email")
```

### 4. 列表操作

```go
queueCache := cache.NewNamespace(cacheClient, "queue")

// 从左侧推入
queueCache.LPush(ctx, "tasks", "task1", "task2", "task3")

// 从右侧弹出（FIFO）
task, _ := queueCache.RPop(ctx, "tasks")  // "task1"

// 获取列表长度
length, _ := queueCache.LLen(ctx, "tasks")
```

### 5. 集合操作

```go
tagCache := cache.NewNamespace(cacheClient, "tag")

// 添加标签
tagCache.SAdd(ctx, "article:1", "golang", "redis", "cache")

// 获取所有标签
tags, _ := tagCache.SMembers(ctx, "article:1")
// tags = ["golang", "redis", "cache"]

// 删除标签
tagCache.SRem(ctx, "article:1", "cache")
```

## 子命名空间

### 创建层级命名空间

```go
// 创建应用命名空间
appCache := cache.NewNamespace(cacheClient, "app")

// 创建用户子命名空间
userCache := appCache.SubNamespace("user")        // 前缀: app:user

// 创建用户资料子命名空间
profileCache := userCache.SubNamespace("profile") // 前缀: app:user:profile

// 设置数据
profileCache.Set(ctx, "123", "ProfileData", time.Hour)
// 实际键: app:user:profile:123
```

### 使用场景

```go
// 多租户应用
tenantCache := cache.NewNamespace(cacheClient, "tenant")
tenant1Cache := tenantCache.SubNamespace("1001")  // tenant:1001
tenant2Cache := tenantCache.SubNamespace("1002")  // tenant:1002

// 每个租户的用户缓存
tenant1UserCache := tenant1Cache.SubNamespace("user")  // tenant:1001:user
tenant2UserCache := tenant2Cache.SubNamespace("user")  // tenant:1002:user
```

## 自定义分隔符

```go
// 使用 "-" 作为分隔符
customCache := cache.NewNamespace(cacheClient, "custom", cache.WithSeparator("-"))

customCache.Set(ctx, "key", "value", time.Hour)
// 实际键: custom-key （而不是 custom:key）
```

## 命名空间管理

### 清空命名空间

```go
testCache := cache.NewNamespace(cacheClient, "test")

// 设置一些数据
testCache.Set(ctx, "key1", "value1", time.Hour)
testCache.Set(ctx, "key2", "value2", time.Hour)
testCache.Set(ctx, "key3", "value3", time.Hour)

// 清空整个命名空间（只删除 test:* 的键）
testCache.FlushNamespace(ctx)
```

### 查找键

```go
userCache := cache.NewNamespace(cacheClient, "user")

// 设置一些用户数据
userCache.Set(ctx, "user:1", "Alice", time.Hour)
userCache.Set(ctx, "user:2", "Bob", time.Hour)
userCache.Set(ctx, "admin:1", "Admin", time.Hour)

// 查找所有用户键（返回的键已移除命名空间前缀）
keys, _ := userCache.Keys(ctx, "user:*")
// keys = ["user:1", "user:2"]
```

### 获取命名空间信息

```go
ns := cache.NewNamespace(cacheClient, "myapp")

// 获取前缀
prefix := ns.Prefix()  // "myapp"

// 获取分隔符
sep := ns.Separator()  // ":"

// 获取底层缓存实例
baseCache := ns.Unwrap()
```

## 实际应用场景

### 场景 1：用户会话管理

```go
sessionCache := cache.NewNamespace(cacheClient, "session")

// 创建会话
sessionID := "abc123"
sessionData := map[string]interface{}{
	"user_id": 123,
	"username": "alice",
	"role": "admin",
}
sessionCache.HSet(ctx, sessionID, "user_id", sessionData["user_id"])
sessionCache.HSet(ctx, sessionID, "username", sessionData["username"])
sessionCache.HSet(ctx, sessionID, "role", sessionData["role"])
sessionCache.Expire(ctx, sessionID, 30*time.Minute)

// 验证会话
userData, _ := sessionCache.HGetAll(ctx, sessionID)

// 延长会话
sessionCache.Expire(ctx, sessionID, 30*time.Minute)

// 销毁会话
sessionCache.Delete(ctx, sessionID)
```

### 场景 2：商品库存计数

```go
inventoryCache := cache.NewNamespace(cacheClient, "inventory")

// 初始化库存
productID := "P001"
inventoryCache.Set(ctx, productID, 100, 0) // 不过期

// 减少库存
remaining, _ := inventoryCache.DecrBy(ctx, productID, 1)
if remaining < 0 {
	// 库存不足
	inventoryCache.IncrBy(ctx, productID, 1) // 回滚
	return errors.New("insufficient inventory")
}

// 增加库存（补货）
inventoryCache.IncrBy(ctx, productID, 50)
```

### 场景 3：访问限流

```go
rateLimitCache := cache.NewNamespace(cacheClient, "ratelimit")

func checkRateLimit(userID string) bool {
	key := fmt.Sprintf("user:%s", userID)

	// 获取当前计数
	count, err := rateLimitCache.Incr(ctx, key)
	if err != nil {
		return false
	}

	// 首次访问，设置过期时间
	if count == 1 {
		rateLimitCache.Expire(ctx, key, time.Minute)
	}

	// 检查是否超过限制（例如：每分钟100次）
	return count <= 100
}
```

### 场景 4：缓存穿透防护（布隆过滤器模拟）

```go
bloomCache := cache.NewNamespace(cacheClient, "bloom")

// 记录已存在的 ID
func markIDExists(id string) {
	bloomCache.SAdd(ctx, "valid_ids", id)
}

// 检查 ID 是否可能存在
func mightIDExist(id string) bool {
	// 注意：这是简化版，真正的布隆过滤器需要更复杂的实现
	members, _ := bloomCache.SMembers(ctx, "valid_ids")
	for _, member := range members {
		if member == id {
			return true
		}
	}
	return false
}
```

### 场景 5：分布式锁

```go
lockCache := cache.NewNamespace(cacheClient, "lock")

func acquireLock(resourceID string, ttl time.Duration) bool {
	lockKey := fmt.Sprintf("resource:%s", resourceID)

	// 使用 SetNX 获取锁
	acquired, _ := lockCache.SetNX(ctx, lockKey, "locked", ttl)
	return acquired
}

func releaseLock(resourceID string) {
	lockKey := fmt.Sprintf("resource:%s", resourceID)
	lockCache.Delete(ctx, lockKey)
}

// 使用示例
if acquireLock("order:123", 10*time.Second) {
	defer releaseLock("order:123")

	// 处理订单
	processOrder("123")
}
```

## 性能优化建议

### 1. 合理使用批量操作

```go
// ❌ 不推荐：循环单次操作
for id := range userIDs {
	userCache.Get(ctx, id)
}

// ✅ 推荐：使用批量操作
values, _ := userCache.MGet(ctx, userIDs...)
```

### 2. 避免过长的键名

```go
// ❌ 不推荐：键名过长
userCache.Set(ctx, "user:profile:personal:information:detailed:123", data, time.Hour)

// ✅ 推荐：简洁的键名
userCache.Set(ctx, "123:profile", data, time.Hour)
```

### 3. 合理设置过期时间

```go
// 热点数据：短过期时间
hotDataCache := cache.NewNamespace(cacheClient, "hot")
hotDataCache.Set(ctx, "trending:1", data, 5*time.Minute)

// 冷数据：长过期时间
coldDataCache := cache.NewNamespace(cacheClient, "cold")
coldDataCache.Set(ctx, "archive:1", data, 24*time.Hour)
```

## 注意事项

1. **命名空间前缀建议**
    - 使用有意义的英文单词
    - 避免使用特殊字符
    - 保持简短（建议不超过 20 个字符）

2. **键命名规范**
    - 使用清晰的命名规则
    - 建议格式：`类型:ID` 或 `ID:属性`
    - 例如：`user:123`, `123:profile`, `order:456:items`

3. **性能考虑**
    - `Keys()` 和 `FlushNamespace()` 是 O(N) 操作，生产环境谨慎使用
    - 优先使用批量操作而非循环单次操作
    - 合理设置过期时间，避免内存浪费

4. **安全性**
    - 不要在键中存储敏感信息
    - 使用 `FlushNamespace()` 而非 `FlushDB()` 避免误删其他数据

## 总结

Namespace 提供了：
- ✅ 键命名隔离，避免冲突
- ✅ 完整的 Cache 接口实现
- ✅ 灵活的层级命名空间
- ✅ 自定义分隔符支持
- ✅ 安全的命名空间清空操作

适用于：
- 多业务模块的应用
- 多租户系统
- 微服务架构
- 需要键隔离的任何场景
