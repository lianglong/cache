# Cache - 通用缓存接口库

[![Go Version](https://img.shields.io/badge/go-1.19+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## 简介

Cache 是一个通用的 Go 缓存接口库，提供统一的 API 来操作不同的缓存后端（Redis、Memcached 等）。

## 特性

- 🔌 **插件化架构** - 驱动注册机制，支持多种缓存后端
- 🚀 **完整的接口** - 30+ 方法，覆盖所有常用缓存操作
- 🏷️ **命名空间** - 自动键前缀隔离，避免键冲突
- ⚡ **批量操作** - MGet、MSet、MDelete 高性能批处理
- 🔢 **原子计数** - Incr、Decr 原子递增递减
- 📦 **数据结构** - Hash、List、Set 支持
- ⏰ **TTL 管理** - 灵活的过期时间控制
- 🔧 **连接池** - 完整的连接池配置支持

## 安装

```bash
go get github.com/lianglong/cache
```

## 快速开始

### 基础使用

```go
package main

import (
"context"
"log"
"time"

    "github.com/lianglong/cache"
    _ "github.com/lianglong/cache-redis"  // 导入 Redis 驱动
)

func main() {
ctx := context.Background()

    // 创建缓存实例
    c, err := cache.New("redis", cache.Config{
        Addr:     "127.0.0.1:6379",
        Password: "your-password",
        DB:       0,
    })
    if err != nil {
        log.Fatal(err)
    }
    defer c.Close()

    // 设置值
    err = c.Set(ctx, "key", "value", time.Hour)

    // 获取值
    value, err := c.Get(ctx, "key")
    log.Println(value)  // "value"
}
```

### 命名空间使用

```go
// 创建命名空间
userCache := cache.NewNamespace(c, "user")
sessionCache := cache.NewNamespace(c, "session")

// 不同命名空间的相同键是隔离的
userCache.Set(ctx, "123", "Alice", time.Hour)       // 实际键: user:123
sessionCache.Set(ctx, "123", "SessionData", time.Hour)  // 实际键: session:123
```


## 文档

- [完整文档](https://pkg.go.dev/github.com/lianglong/cache)
- [命名空间使用指南](NAMESPACE_GUIDE.md)
- [驱动开发指南](DRIVER_DEVELOPMENT.md)

## 支持的驱动

- ✅ **Redis** - [cache-redis](https://github.com/lianglong/cache-redis)
- 🚧 **Memcached** - 开发中
- 🚧 **内存缓存** - 计划中

## 开发自定义驱动

查看 [驱动开发指南](DRIVER_DEVELOPMENT.md) 了解如何开发自定义驱动。

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License - 查看 [LICENSE](LICENSE) 文件