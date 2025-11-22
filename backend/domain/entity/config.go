package entity

import "time"

type AppConfig struct {
	Server     ServerConfig     `mapstructure:"server"`
	Database   DatabaseConfig   `mapstructure:"database"`
	Redis      RedisConfig      `mapstructure:"redis"`
	Sync       SyncConfig       `mapstructure:"sync"`
	Cache      CacheConfig      `mapstructure:"cache"`
	Pagination PaginationConfig `mapstructure:"pagination"`
	RateLimit  RateLimitConfig  `mapstructure:"rate_limit"`
	CircuitBreaker CircuitBreakerConfig `mapstructure:"circuit_breaker"`
}

type RateLimitConfig struct {
	RPS   int `mapstructure:"rps"`
	Burst int `mapstructure:"burst"`
}

type CircuitBreakerConfig struct {
	Name          string `mapstructure:"name"`
	MaxRequests   uint32 `mapstructure:"max_requests"`
	Interval      int    `mapstructure:"interval"`
	Timeout       int    `mapstructure:"timeout"`
	ReadyToTrip   bool   `mapstructure:"ready_to_trip"` // Simplified for config, logic in code
}
type PaginationConfig struct {
	DefaultPage     int `mapstructure:"default_page"`
	DefaultPageSize int `mapstructure:"default_page_size"`
	MaxPageSize     int `mapstructure:"max_page_size"`
}

type ServerConfig struct {
	GRPCPort int `mapstructure:"grpc_port"`
	HTTPPort int `mapstructure:"http_port"`
}

type DatabaseConfig struct {
	DSN string `mapstructure:"dsn"`
}

type RedisConfig struct {
	Addr string `mapstructure:"addr"`
}

type SyncConfig struct {
	IntervalSeconds int `mapstructure:"interval_seconds"`
}

func (c SyncConfig) GetInterval() time.Duration {
	if c.IntervalSeconds <= 0 {
		return 60 * time.Second
	}
	return time.Duration(c.IntervalSeconds) * time.Second
}

type CacheConfig struct {
	TTLSeconds int `mapstructure:"ttl_seconds"`
}

func (c CacheConfig) GetTTL() time.Duration {
	if c.TTLSeconds <= 0 {
		return 3600 * time.Second
	}
	return time.Duration(c.TTLSeconds) * time.Second
}
