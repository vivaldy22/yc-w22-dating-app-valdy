package redis

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/joomcode/errorx"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"

	ierror "yc-w22-dating-app-valdy/pkg/error"
)

type (
	RateLimitRepository interface {
		Incr(ctx context.Context, rateLimitType, rateLimitIdentifier string, expire time.Duration) (int64, error)
		Get(ctx context.Context, rateLimitType, rateLimitIdentifier string) (int64, error)
	}

	rateLimitRepository struct {
		client *redis.Client
	}
)

func NewRateLimitRepository(client *redis.Client) RateLimitRepository {
	if client == nil {
		panic("redis client is nil")
	}

	return &rateLimitRepository{
		client: client,
	}
}

func (r *rateLimitRepository) Incr(ctx context.Context, rateLimitType, rateLimitIdentifier string, expire time.Duration) (int64, error) {
	key := fmt.Sprintf("%s#%s", rateLimitType, rateLimitIdentifier)
	val, err := r.Get(ctx, rateLimitType, rateLimitIdentifier)
	if errorx.IsNotFound(err) {
		err = r.client.Set(ctx, key, 1, expire).Err()
		if err != nil {
			log.Printf("failed to set rate limit %s to %s, error: %s", rateLimitType, rateLimitIdentifier, err.Error())
			return 0, ierror.ErrDatabase
		}
		return 1, nil
	}
	if err != nil {
		log.Printf("failed to get rate limit %s to %s, error: %s", rateLimitType, rateLimitIdentifier, err.Error())
		return 0, ierror.ErrDatabase
	}

	val, err = r.client.Incr(ctx, key).Result()
	if err != nil {
		log.Printf("failed to incr rate limit %s to %s, error: %s", rateLimitType, rateLimitIdentifier, err.Error())
		return 0, ierror.ErrDatabase
	}

	return val, nil
}

func (r *rateLimitRepository) Get(ctx context.Context, rateLimitType, rateLimitIdentifier string) (int64, error) {
	key := fmt.Sprintf("%s#%s", rateLimitType, rateLimitIdentifier)
	val, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		log.Printf("rate limit %s not exist", rateLimitIdentifier)
		return 0, ierror.ErrDataNotFound
	}
	if err != nil {
		log.Printf("failed to get rate limit %s to %s, error: %s", rateLimitType, rateLimitIdentifier, err.Error())
		return 0, ierror.ErrDatabase
	}

	return cast.ToInt64(val), nil
}
