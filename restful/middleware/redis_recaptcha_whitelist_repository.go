package middleware

import (
	"context"
	"errors"
	"github.com/cross-space-official-private/common/configuration"
	"github.com/cross-space-official-private/common/ratelimiter"
	"github.com/go-redis/redis/v8"
	"time"
)

type redisRecaptchaWhitelistRepository struct {
	redisClient *redis.Client
}

func (r *redisRecaptchaWhitelistRepository) AddWhitelist(c context.Context, profileID string) error {
	if r.redisClient == nil {
		return errors.New("failed to add whitelist")
	}

	err := r.redisClient.Set(c, r.generateKey(profileID), true, 24*time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *redisRecaptchaWhitelistRepository) ShouldSkipByProfileID(c context.Context, profileID string) bool {
	if r.redisClient == nil {
		return false
	}

	res := r.redisClient.Get(c, r.generateKey(profileID))
	if res == nil || res.Err() == redis.Nil {
		return false
	}

	return true
}

func (r *redisRecaptchaWhitelistRepository) generateKey(profileID string) string {
	return "recaptcha_whitelist_" + profileID
}

func NewRedisRecaptchaWhitelistRepository() RecaptchaProfileWhitelistRepository {
	redisConfig := &ratelimiter.RedisConfig{}
	configuration.Build("data.redis", redisConfig)

	redisClient := redis.NewClient(redisConfig.ToOptions())
	return &redisRecaptchaWhitelistRepository{redisClient: redisClient}
}
