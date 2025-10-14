package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"saas/internal/comment/domain"
	"saas/internal/common/reskit/codes"
	"saas/internal/common/utils"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type CommentRedisCache struct {
	client *redis.Client
}

func NewCommentRedisCache() domain.CommentCache {
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

	return &CommentRedisCache{client: client}
}

const CommentTenantConfigKey = "comment:tenant:config"
const CommentTenantConfigExpired = 1 * time.Hour
const commentPlateConfigKey = "comment:plate:config"
const commentPlateConfigExpired = 1 * time.Hour

func (cache *CommentRedisCache) SetTenantConfig(config *domain.TenantConfig) error {
	key := utils.GetRedisKey(CommentTenantConfigKey)
	field := fmt.Sprintf("%d", config.TenantID)
	data, err := json.Marshal(config)
	if err != nil {
		return errors.WithStack(err)
	}

	pipeline := cache.client.Pipeline()
	pipeline.HSet(context.Background(), key, field, data)
	pipeline.HExpire(context.Background(), key, CommentTenantConfigExpired, field).Err()
	_, err = pipeline.Exec(context.Background())
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (cache *CommentRedisCache) GetTenantConfig(tenantID domain.TenantID) (*domain.TenantConfig, error) {
	key := utils.GetRedisKey(CommentTenantConfigKey)
	field := fmt.Sprintf("%d", tenantID)

	result, err := cache.client.HGet(context.Background(), key, field).Result()
	if err != nil {
		return nil, errors.WithStack(codes.ErrCommentTenantConfigCacheMissing)
	}

	config := new(domain.TenantConfig)

	if err := json.Unmarshal([]byte(result), config); err != nil {
		return nil, errors.WithStack(err)
	}

	return config, nil
}

func (cache *CommentRedisCache) DeleteTenantConfig(tenantID domain.TenantID) error {
	key := utils.GetRedisKey(CommentTenantConfigKey)
	field := fmt.Sprintf("%d", tenantID)
	if err := cache.client.HDel(context.Background(), key, field).Err(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (cache *CommentRedisCache) SetPlateConfig(config *domain.PlateConfig) error {
	key := utils.GetRedisKey(commentPlateConfigKey)
	field := fmt.Sprintf("%d-%d", config.TenantID, config.Plate.ID)
	data, err := json.Marshal(config)
	if err != nil {
		return errors.WithStack(err)
	}

	pipeline := cache.client.Pipeline()

	pipeline.HSet(context.Background(), key, field, data)
	pipeline.HExpire(context.Background(), key, commentPlateConfigExpired, field).Err()
	_, err = pipeline.Exec(context.Background())
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (cache *CommentRedisCache) GetPlateConfig(tenantID domain.TenantID, plateID int64) (*domain.PlateConfig, error) {
	key := utils.GetRedisKey(commentPlateConfigKey)
	field := fmt.Sprintf("%d-%d", tenantID, plateID)

	result, err := cache.client.HGet(context.Background(), key, field).Result()
	if err != nil {
		return nil, errors.WithStack(codes.ErrCommentPlateConfigCacheMissing)
	}

	config := new(domain.PlateConfig)

	if err := json.Unmarshal([]byte(result), config); err != nil {
		return nil, errors.WithStack(err)
	}

	return config, nil
}

func (cache *CommentRedisCache) DeletePlateConfig(tenantID domain.TenantID, plateID int64) error {
	key := utils.GetRedisKey(commentPlateConfigKey)
	field := fmt.Sprintf("%d-%d", tenantID, plateID)
	if err := cache.client.HDel(context.Background(), key, field).Err(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
