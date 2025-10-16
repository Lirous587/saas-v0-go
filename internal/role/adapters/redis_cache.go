package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"saas/internal/common/utils"
	"saas/internal/role/domain"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type RoleRedisCache struct {
	client *redis.Client
}

func NewRoleRedisCache() domain.RoleCache {
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

	return &RoleRedisCache{client: client}
}

const userRoleInTenantPreKey = "user_role_in_tenant"

const userRoleInTenantExpire = time.Hour * 2

// const userRoleInTenantExpire = time.Second * 1

func (cache *RoleRedisCache) buildRoleKey(userID, tenantID int64) string {
	key := utils.GetRedisKey(userRoleInTenantPreKey)
	return fmt.Sprintf("%s:%d:%d", key, userID, tenantID)
}

func (cache *RoleRedisCache) GetUserRoleInTenant(userID, tenantID int64) (*domain.Role, error) {
	key := cache.buildRoleKey(userID, tenantID)
	result, err := cache.client.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	role := new(domain.Role)
	if err = json.Unmarshal([]byte(result), role); err != nil {
		return nil, errors.WithMessage(err, "json Unmarshal失败")
	}

	return role, nil
}

func (cache *RoleRedisCache) SetUserRoleInTenant(userID, tenantID int64, role *domain.Role) error {
	key := cache.buildRoleKey(userID, tenantID)
	data, err := json.Marshal(role)
	if err != nil {
		return errors.WithMessage(err, "json Marshal失败")
	}
	return cache.client.Set(context.Background(), key, string(data), userRoleInTenantExpire).Err()
}
