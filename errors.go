package cache

import "errors"

// 标准错误
var (
	ErrNotFound       = errors.New("cache: key not found")
	ErrInvalidValue   = errors.New("cache: invalid value")
	ErrKeyExpired     = errors.New("cache: key expired")
	ErrConnectionLost = errors.New("cache: connection lost")
	ErrTimeout        = errors.New("cache: operation timeout")
	ErrCacheFull      = errors.New("cache: cache is full")
)

// IsNotFound 判断是否为未找到错误
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsTimeout 判断是否为超时错误
func IsTimeout(err error) bool {
	return errors.Is(err, ErrTimeout)
}
