package adapters

import (
	"context"
	"fmt"
	"os"
	"saas/internal/common/reskit/codes"
	"saas/internal/common/utils"
	"saas/internal/tenant/domain"
	"strconv"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/redis/go-redis/v9"
)

type TenantRedisCache struct {
	client *redis.Client
}

func NewTenantRedisCache() domain.TenantCache {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")
	dbStr := os.Getenv("REDIS_DB")
	poolSizeStr := os.Getenv("REDIS_POOL_SIZE")

	db, _ := strconv.Atoi(dbStr)
	poolSize, _ := strconv.Atoi(poolSizeStr)

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

	return &TenantRedisCache{client: client}
}

const publicInvitedMap = "invite:public"

func (cache *TenantRedisCache) buildPublicInviteKey(tenantID int64) string {
	return fmt.Sprintf("%s:%d", utils.GetRedisKey(publicInvitedMap), tenantID)
}

func (cache *TenantRedisCache) GenPublicInviteToken(tenantID int64, expireSecond time.Duration) (string, error) {
	token, err := utils.GenRandomHexToken()
	if err != nil {
		return "", errors.WithStack(err)
	}

	key := cache.buildPublicInviteKey(tenantID)

	// 将令牌加入Set
	if err := cache.client.HSet(context.Background(), key, token, "1").Err(); err != nil {
		return "", errors.WithStack(err)
	}

	// 设置列表过期时间
	if err := cache.client.HExpire(context.Background(), key, expireSecond*time.Second, token).Err(); err != nil {
		return "", errors.WithStack(err)
	}

	return token, nil
}

func (cache *TenantRedisCache) ValidatePublicInviteToken(tenantID int64, value string) error {
	key := cache.buildPublicInviteKey(tenantID)

	result, err := cache.client.HGet(context.Background(), key, value).Result()
	if err != nil {
		return err
	}

	if result != "1" {
		return codes.ErrTenantInviteTokenInvalid
	}

	return nil
}

const secretInvitedMap = "invite:secret"

func (cache *TenantRedisCache) buildSecretInviteKey(tenantID int64) string {
	return utils.GetRedisKey(fmt.Sprintf("%s:%d", secretInvitedMap, tenantID))
}

func (cache *TenantRedisCache) GenSecretInviteToken(tenantID int64, expireSecond time.Duration, email string) (string, error) {
	key := cache.buildSecretInviteKey(tenantID)

	token, err := utils.GenRandomHexToken()
	if err != nil {
		return "", errors.WithStack(err)
	}

	// 使用email作为field，token作为value
	fileds := []interface{}{email, token}

	if err := cache.client.HSet(context.Background(), key, fileds...).Err(); err != nil {
		return "", errors.WithStack(err)
	}

	// 为指定字段设置过期时间
	if err := cache.client.HExpire(context.Background(), key, expireSecond*time.Second, email).Err(); err != nil {
		return "", errors.WithStack(err)
	}

	return token, nil
}

func (cache *TenantRedisCache) ValidateSecretInviteToken(tenantID int64, email, value string) error {
	key := cache.buildSecretInviteKey(tenantID)

	result, err := cache.client.HGet(context.Background(), key, email).Result()
	if err != nil {
		return errors.WithStack(err)
	}

	if result != value {
		return codes.ErrTenantInviteTokenInvalid
	}

	return nil
}

func (cache *TenantRedisCache) DeleteSecretInviteToken(tenantID int64, email string) error {
	key := cache.buildSecretInviteKey(tenantID)

	return cache.client.HDel(context.Background(), key, email).Err()
}
