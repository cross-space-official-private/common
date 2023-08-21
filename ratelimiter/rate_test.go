package ratelimiter

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redis_rate/v9"
	"testing"
	"time"
)

func Test(t *testing.T) {
	t.Skip()

	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	_ = rdb.FlushDB(ctx).Err()

	limiter := redis_rate.NewLimiter(rdb)
	for i := 1; i <= 10; i++ {
		res, err := limiter.Allow(ctx, "project:123", redis_rate.PerHour(1))
		time.Sleep(500 * time.Millisecond)
		if err != nil {
			panic(err)
		}
		fmt.Println("allowed", res.Allowed, "remaining", res.Remaining)
	}

	// Output: allowed 1 remaining 9
}
