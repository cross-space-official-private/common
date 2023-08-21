package ratelimiter

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisConfig struct {
	Host                 string        `mapstructure:"host"`
	Port                 int           `mapstructure:"port"`
	Username             string        `mapstructure:"username"`
	Password             string        `mapstructure:"password"`
	MaxIdleConnections   int           `mapstructure:"max-idle-connections"`
	IdleTimeout          time.Duration `mapstructure:"idle-timeout"`
	MaxActiveConnections int           `mapstructure:"max-active-connections"`
}

func (c *RedisConfig) ToOptions() *redis.Options {
	return &redis.Options{
		Addr:         fmt.Sprintf("%s:%d", c.Host, c.Port),
		Username:     c.Username,
		Password:     c.Password,
		PoolSize:     c.MaxActiveConnections,
		MinIdleConns: c.MaxIdleConnections,
		IdleTimeout:  c.IdleTimeout,
	}
}
