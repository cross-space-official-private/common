package ratelimiter

import (
	"context"
	"github.com/cross-space-official/common/configuration"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redis_rate/v9"
	"sync"
	"time"
)

var limiter RateLimiter
var lock sync.Mutex

type (
	Result struct {
		// Allowed whether the rate limiter accept the new request
		Allowed bool

		// RetryAfter is the time until the next request will be permitted.
		// It should be -1 unless the rate limit has been exceeded.
		RetryAfter time.Duration
	}

	Limit struct {
		// Burst refers to the allowed quota at beginning, if not set, will start from 0
		Burst int

		// Rate and Period refers to the quota regaining ratio
		Rate   int
		Period time.Duration
	}

	RateLimiter interface {
		AllowByKey(c context.Context, key string, limit Limit) (*Result, error)
		ResetKey(c context.Context, key string) error
	}

	redisRateLimiter struct {
		RLimiter *redis_rate.Limiter
	}
)

func GetRateLimiter() RateLimiter {
	if limiter != nil {
		return limiter
	}

	lock.Lock()
	defer lock.Unlock()

	if limiter != nil {
		return limiter
	}

	redisConfig := &RedisConfig{}
	configuration.Build("data.redis", redisConfig)
	rdb := redis.NewClient(redisConfig.ToOptions())

	limiter = &redisRateLimiter{RLimiter: redis_rate.NewLimiter(rdb)}
	return limiter
}

func (rl *redisRateLimiter) AllowByKey(c context.Context, key string, limit Limit) (*Result, error) {
	res, err := rl.RLimiter.Allow(c, key, redis_rate.Limit{
		Rate:   limit.Rate,
		Burst:  limit.Burst,
		Period: limit.Period,
	})

	if err != nil {
		return nil, err
	}

	return &Result{
		Allowed:    res.Allowed == 1,
		RetryAfter: res.RetryAfter,
	}, nil
}

func (rl *redisRateLimiter) ResetKey(c context.Context, key string) error {
	return rl.RLimiter.Reset(c, key)
}
