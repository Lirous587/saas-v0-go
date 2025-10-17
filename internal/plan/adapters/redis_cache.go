package adapters

import (
	"context"
	"saas/internal/common/utils"
	"saas/internal/plan/domain"

	"github.com/redis/go-redis/v9"
)

type PlanRedisCache struct {
	client *redis.Client
}

func NewPlanRedisCache() domain.PlanCache {
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

	return &PlanRedisCache{client: client}
}
