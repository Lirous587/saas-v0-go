package adapters

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"saas/internal/captcha/domain"
	"saas/internal/common/reskit/codes"
	"saas/internal/common/uid"
	"saas/internal/common/utils"
)

type CaptchaRedisCache struct {
	client *redis.Client
}

func NewCaptchaRedisCache() domain.CaptchaCache {
	host := utils.GetEnv("REDIS_HOST")
	port := utils.GetEnv("REDIS_PORT")
	password := utils.GetEnv("REDIS_PASSWORD")
	db := utils.GetEnvAsInt("REDIS_DB")
	poolSize := utils.GetEnvAsInt("REDIS_POOL_SIZE")

	addr := host + ":" + port

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       db,
		Password: password,
		PoolSize: poolSize,
	})

	// 可选：ping 检查连接
	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	return &CaptchaRedisCache{client: client}
}

func buildKey(way domain.VerifyWay, id int64) (string, error) {
	keyPre := utils.GetRedisKey(way.GetKey())
	return fmt.Sprintf("%s:%d", keyPre, id), nil
}

func (r *CaptchaRedisCache) Save(way domain.VerifyWay, value string) (int64, error) {
	id, err := uid.Gen()
	if err != nil {
		return 0, errors.WithStack(err)
	}

	key, err := buildKey(way, id)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	if err := r.client.Set(context.Background(), key, value, way.GetExpire()).Err(); err != nil {
		return 0, errors.WithStack(err)
	}

	return id, nil
}

func (r *CaptchaRedisCache) Verify(way domain.VerifyWay, id int64, value string) error {
	key, err := buildKey(way, id)
	if err != nil {
		return errors.WithStack(err)
	}

	result, err := r.client.Get(context.Background(), key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return codes.ErrCaptchaNotFound
		}
		return err
	}
	if result != value {
		return codes.ErrCaptchaVerifyFailed
	}

	// 验证成功之后删除
	if err := r.Delete(key); err != nil {
		zap.L().Error("删除验证码id缓存失败", zap.String("key", key))
	}

	return nil
}

func (r *CaptchaRedisCache) Delete(key string) error {
	return r.client.Del(context.Background(), key).Err()
}
