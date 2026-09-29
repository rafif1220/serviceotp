package redisrepo

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rafif1220/serviceotp/internal/domain"
)

type cacheRepo struct {
	client *redis.Client
}

// NewCacheRepository adalah constructor buat nyetak object cacheRepo
func NewCacheRepository(client *redis.Client) domain.CacheRepository {
	return &cacheRepo{client: client}
}

func (c *cacheRepo) SetOTP(recipient string, code string, ttl time.Duration) error {
	ctx := context.Background()
	key := "nukar:otp:" + recipient
	return c.client.Set(ctx, key, code, ttl).Err()
}

func (c *cacheRepo) GetOTP(recipient string) (string, error) {
	ctx := context.Background()
	key := "nukar:otp:" + recipient
	return c.client.Get(ctx, key).Result() // Bakal return error redis.Nil kalau OTP nggak ada/expired
}

func (c *cacheRepo) DeleteOTP(recipient string) error {
	ctx := context.Background()
	key := "nukar:otp:" + recipient
	return c.client.Del(ctx, key).Err()
}

func (c *cacheRepo) IncrementRateLimit(key string, ttl time.Duration) (int, error) {
	ctx := context.Background()
	
	// Pake Pipeline biar INCR dan EXPIRE dieksekusi dalam satu tarikan napas (atomic)
	pipe := c.client.TxPipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, ttl)
	
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}
	
	return int(incr.Val()), nil
}
