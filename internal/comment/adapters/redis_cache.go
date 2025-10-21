package adapters

import (
	"context"
	"encoding/json"
	"fmt"
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

	return &CommentRedisCache{client: client}
}

const CommentTenantConfigKey = "comment:tenant:config"
const CommentTenantConfigExpired = 1 * time.Hour

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
		if err == redis.Nil {
			return nil, errors.WithStack(codes.ErrCommentTenantConfigCacheMissing)
		}
		return nil, errors.WithStack(err)
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

const commentPlateID = "comment:plate:id"
const commentPlateIDExpired = 1 * time.Hour

func (cache *CommentRedisCache) SetPlateID(tenantID domain.TenantID, belongKey string, id int64) error {
	key := utils.GetRedisKey(commentPlateID)
	field := fmt.Sprintf("%d-%s", tenantID, belongKey)

	pipeline := cache.client.Pipeline()
	pipeline.HSet(context.Background(), key, field, id)
	pipeline.HExpire(context.Background(), key, commentPlateIDExpired, field).Err()
	_, err := pipeline.Exec(context.Background())
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (cache *CommentRedisCache) GetPlateID(tenantID domain.TenantID, belongKey string) (int64, error) {
	key := utils.GetRedisKey(commentPlateID)
	field := fmt.Sprintf("%d-%s", tenantID, belongKey)

	result, err := cache.client.HGet(context.Background(), key, field).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, errors.WithStack(codes.ErrCommentPlateIDCacheMissing)
		}
		return 0, errors.WithStack(err)
	}

	plateID, err := strconv.ParseInt(result, 10, 64)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return plateID, nil
}

func (cache *CommentRedisCache) DeletePlateID(tenantID domain.TenantID, belongKey string) error {
	key := utils.GetRedisKey(commentPlateID)
	field := fmt.Sprintf("%d-%s", tenantID, belongKey)
	if err := cache.client.HDel(context.Background(), key, field).Err(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

const commentPlateConfigKey = "comment:plate:config"
const commentPlateConfigExpired = 1 * time.Hour

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
		if err == redis.Nil {
			return nil, errors.WithStack(codes.ErrCommentPlateConfigCacheMissing)
		}
		return nil, errors.WithStack(err)
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

const commentLikeKey = "comment:like"
const commentLikeExpired = 30 * time.Hour * 24 // 30缓存 场景足够

func (cache *CommentRedisCache) GetLikeStatus(tenantID domain.TenantID, userID int64, commentID int64) (domain.LikeStatus, error) {
	preKey := utils.GetRedisKey(commentLikeKey)
	key := fmt.Sprintf("%s:%d-%d", preKey, tenantID, userID)
	exists, err := cache.client.SIsMember(context.Background(), key, commentID).Result()
	if err != nil {
		return false, errors.WithStack(err)
	}

	likeStatus := new(domain.LikeStatus)

	if exists {
		likeStatus.Like()
		return *likeStatus, nil // true，表示已点赞
	}
	likeStatus.UnLike()
	return *likeStatus, nil // false，表示未点赞
}

func (cache *CommentRedisCache) AddLike(tenantID domain.TenantID, userID int64, commentID int64) error {
	preKey := utils.GetRedisKey(commentLikeKey)
	key := fmt.Sprintf("%s-%d-%d", preKey, tenantID, userID)
	if err := cache.client.SAdd(context.Background(), key, commentID).Err(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (cache *CommentRedisCache) RemoveLike(tenantID domain.TenantID, userID int64, commentID int64) error {
	preKey := utils.GetRedisKey(commentLikeKey)
	key := fmt.Sprintf("%s-%d-%d", preKey, tenantID, userID)
	if err := cache.client.SRem(context.Background(), key, commentID).Err(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (cache *CommentRedisCache) GetLikeMap(tenantID domain.TenantID, userID int64, commentIds []int64) (map[int64]struct{}, error) {
	if len(commentIds) == 0 {
		return nil, nil
	}

	preKey := utils.GetRedisKey(commentLikeKey)
	key := fmt.Sprintf("%s-%d-%d", preKey, tenantID, userID)

	members := utils.Int64SliceToInterface(commentIds)

	results, err := cache.client.SMIsMember(context.Background(), key, members...).Result()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	likeMap := make(map[int64]struct{})
	for i, exists := range results {
		if exists {
			likeMap[commentIds[i]] = struct{}{}
		}
	}

	return likeMap, nil
}
