package cache

import "time"

// Config 通用缓存配置
type Config struct {
	// 连接地址
	Addr string

	// 认证密码
	Password string

	// 数据库编号（Redis 专用）
	DB int

	// 连接超时
	DialTimeout time.Duration

	// 读超时
	ReadTimeout time.Duration

	// 写超时
	WriteTimeout time.Duration

	// 连接池配置
	Pool PoolConfig

	// 驱动特定配置（可选）
	Extra map[string]interface{}
}

// PoolConfig 连接池配置
type PoolConfig struct {
	// 最大空闲连接数
	MaxIdleConns int

	// 最大活跃连接数
	MaxActiveConns int

	// 连接最大空闲时间
	IdleTimeout time.Duration

	// 连接最大生命周期
	MaxConnLifetime time.Duration
}

// DefaultConfig 返回默认配置
func DefaultConfig() Config {
	return Config{
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		Pool: PoolConfig{
			MaxIdleConns:    10,
			MaxActiveConns:  100,
			IdleTimeout:     5 * time.Minute,
			MaxConnLifetime: 1 * time.Hour,
		},
	}
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.Addr == "" {
		return ErrInvalidConfig("addr is required")
	}
	return nil
}

// ErrInvalidConfig 配置错误
type ErrInvalidConfig string

func (e ErrInvalidConfig) Error() string {
	return "invalid config: " + string(e)
}
