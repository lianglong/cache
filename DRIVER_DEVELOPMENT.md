# 驱动开发指南

## 概述

本文档介绍如何为 cache 包开发自定义驱动。

## 驱动结构

驱动需要实现 `cache.Cache` 接口的所有方法。

## 步骤

### 1. 创建新项目

```bash
mkdir cache-mydriver
cd cache-mydriver
go mod init github.com/yourname/cache-mydriver
```

### 2. 实现 Cache 接口

```go
package mydriver

import (
"context"
"time"

    "github.com/lianglong/cache"
)

type myCache struct {
// 你的驱动实现
}

func New(config cache.Config) (cache.Cache, error) {
// 验证配置
if err := config.Validate(); err != nil {
return nil, err
}

    // 创建客户端
    client := createClient(config)

    return &myCache{client: client}, nil
}

// 实现所有接口方法...
func (m *myCache) Get(ctx context.Context, key string) (string, error) {
// 实现
}

func (m *myCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
// 实现
}

// ... 其他 28+ 方法
```

### 3. 注册驱动

```go
package mydriver

import "github.com/lianglong/cache"

func init() {
cache.Register("mydriver", func(config cache.Config) (cache.Cache, error) {
return New(config)
})
}
```

### 4. 测试驱动

创建完整的测试套件，确保所有方法正确实现。

### 5. 发布驱动

将驱动发布为独立的 Go 模块。