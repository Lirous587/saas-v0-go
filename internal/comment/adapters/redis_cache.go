package adapters

import (
	"context"
	"os"
	"saas/internal/comment/domain"
	"strconv"

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

// const CommentTenantConfigKey = "tenant:comment:config"
// const commentConfigKey = "tenant:comment:config"

func (cache *CommentRedisCache) SetTenantCommentClientToken(tenantID domain.TenantID, clientToken string) error {

	return nil
}

func (cache *CommentRedisCache) GetTenantCommentClientToken(tenantID domain.TenantID, benlongKey domain.BelongKey) (string, error) {

	return "", nil
}

func (cache *CommentRedisCache) SetCommentClientToken(tenantID domain.TenantID, belongKey domain.BelongKey, clientToken string) error {
	return nil
}

func (cache *CommentRedisCache) GetCommentClientToken(tenantID domain.TenantID, benlongKey domain.BelongKey) (string, error) {
	return "", nil
}
