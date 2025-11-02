package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"saas/internal/comment/domain"
	"saas/internal/common/reskit/codes"
	"saas/internal/common/utils"
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
	data, err := json.Marshal(config)
	if err != nil {
		return errors.WithStack(err)
	}

	pipeline := cache.client.Pipeline()
	pipeline.HSet(context.Background(), key, config.TenantID, data)
	pipeline.HExpire(context.Background(), key, CommentTenantConfigExpired, config.TenantID.String()).Err()
	_, err = pipeline.Exec(context.Background())
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (cache *CommentRedisCache) GetTenantConfig(tenantID domain.TenantID) (*domain.TenantConfig, error) {
	key := utils.GetRedisKey(CommentTenantConfigKey)

	result, err := cache.client.HGet(context.Background(), key, tenantID.String()).Result()
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
	if err := cache.client.HDel(context.Background(), key, tenantID.String()).Err(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

const commentPlateID = "comment:plate:id"
const commentPlateIDExpired = 1 * time.Hour

func (cache *CommentRedisCache) SetPlateID(tenantID domain.TenantID, belongKey string, plateID domain.PlateID) error {
	key := utils.GetRedisKey(commentPlateID)
	field := fmt.Sprintf("%s-%s", tenantID, belongKey)

	pipeline := cache.client.Pipeline()
	pipeline.HSet(context.Background(), key, field, plateID)
	pipeline.HExpire(context.Background(), key, commentPlateIDExpired, field).Err()
	_, err := pipeline.Exec(context.Background())
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (cache *CommentRedisCache) GetPlateID(tenantID domain.TenantID, belongKey string) (domain.PlateID, error) {
	key := utils.GetRedisKey(commentPlateID)
	field := fmt.Sprintf("%s-%s", tenantID, belongKey)

	plateID, err := cache.client.HGet(context.Background(), key, field).Result()
	if err != nil {
		if err == redis.Nil {
			return "", errors.WithStack(codes.ErrCommentPlateIDCacheMissing)
		}
		return "", errors.WithStack(err)
	}

	return domain.PlateID(plateID), nil
}

func (cache *CommentRedisCache) DeletePlateID(tenantID domain.TenantID, belongKey string) error {
	key := utils.GetRedisKey(commentPlateID)
	field := fmt.Sprintf("%s-%s", tenantID, belongKey)
	if err := cache.client.HDel(context.Background(), key, field).Err(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

const commentPlateConfigKey = "comment:plate:config"
const commentPlateConfigExpired = 1 * time.Hour

func (cache *CommentRedisCache) SetPlateConfig(config *domain.PlateConfig) error {
	key := utils.GetRedisKey(commentPlateConfigKey)
	field := fmt.Sprintf("%s-%s", config.TenantID, config.Plate.ID)
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

func (cache *CommentRedisCache) GetPlateConfig(tenantID domain.TenantID, plateID domain.PlateID) (*domain.PlateConfig, error) {
	key := utils.GetRedisKey(commentPlateConfigKey)
	field := fmt.Sprintf("%s-%s", tenantID, plateID)

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

func (cache *CommentRedisCache) DeletePlateConfig(tenantID domain.TenantID, plateID domain.PlateID) error {
	key := utils.GetRedisKey(commentPlateConfigKey)
	field := fmt.Sprintf("%s-%s", tenantID, plateID)
	if err := cache.client.HDel(context.Background(), key, field).Err(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

const commentLikeKey = "comment:like"
const commentLikeExpired = 30 * time.Hour * 24 // 30缓存 场景足够

func (cache *CommentRedisCache) buildLikeKey(tenantID domain.TenantID, userID domain.UserID) string {
	preKey := utils.GetRedisKey(commentLikeKey)
	return fmt.Sprintf("%s:%s:%s", preKey, tenantID, userID)
}

func (cache *CommentRedisCache) GetLikeStatus(tenantID domain.TenantID, userID domain.UserID, commentID domain.CommentID) (domain.LikeStatus, error) {
	key := cache.buildLikeKey(tenantID, userID)

	exists, err := cache.client.SIsMember(context.Background(), key, commentID.String()).Result()
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

func (cache *CommentRedisCache) AddLike(tenantID domain.TenantID, userID domain.UserID, commentID domain.CommentID) error {
	key := cache.buildLikeKey(tenantID, userID)

	if err := cache.client.SAdd(context.Background(), key, commentID.String()).Err(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (cache *CommentRedisCache) RemoveLike(tenantID domain.TenantID, userID domain.UserID, commentID domain.CommentID) error {
	key := cache.buildLikeKey(tenantID, userID)

	if err := cache.client.SRem(context.Background(), key, commentID.String()).Err(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (cache *CommentRedisCache) GetLikeMap(tenantID domain.TenantID, userID domain.UserID, commentIDs []domain.CommentID) (map[domain.CommentID]struct{}, error) {
	if len(commentIDs) == 0 {
		return nil, nil
	}

	key := cache.buildLikeKey(tenantID, userID)

	commentIDsStr := domain.CommentIDs(commentIDs).ToStringSlice()

	members := utils.StringSliceToInterface(commentIDsStr)

	results, err := cache.client.SMIsMember(context.Background(), key, members...).Result()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	likeMap := make(map[domain.CommentID]struct{})
	for i, exists := range results {
		if exists {
			likeMap[commentIDs[i]] = struct{}{}
		}
	}

	return likeMap, nil
}
