package adapters

import (
	"context"
	"fmt"
	"saas/internal/common/utils"
	"saas/internal/img/domain"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type ImgRedisCache struct {
	client *redis.Client
}

func NewImgRedisCache() domain.ImgMsgQueue {
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

	return &ImgRedisCache{client: client}
}

const (
	keyImgDeleteQueueKey = "img:delete"
	deleteImgExpire      = time.Hour * 24 * 7
	//deleteImgExpire = time.Second * 20
)

func (c *ImgRedisCache) buildDeletedQueueKey(tenantID domain.TenantID, imgID int64) string {
	return fmt.Sprintf("%s:%d:%d", utils.GetRedisKey(keyImgDeleteQueueKey), tenantID, imgID)
}

// AddToDeleteQueue 软删除时写入 redis 并设置过期
func (c *ImgRedisCache) AddToDeleteQueue(tenantID domain.TenantID, imgID int64) error {
	key := c.buildDeletedQueueKey(tenantID, imgID)
	return c.client.SetEx(context.Background(), key, "", deleteImgExpire).Err()
}

// ListenDeleteQueue 后台监听 key 过期事件
func (c *ImgRedisCache) ListenDeleteQueue(onExpire func(tenantID domain.TenantID, imgID int64)) {
	pubsub := c.client.PSubscribe(context.Background(), "__keyevent@0__:expired")
	defer pubsub.Close() // 确保资源释放

	preKey := utils.GetRedisKey(keyImgDeleteQueueKey) + ":"

	for msg := range pubsub.Channel() {
		// 解析 key: img:delete:tenantID:imgID
		if strings.HasPrefix(msg.Payload, preKey) {
			parts := strings.Split(msg.Payload, ":")
			if len(parts) == 4 { // 确保格式正确：img:delete:tenantID:imgID
				tenantID, err1 := strconv.ParseInt(parts[2], 10, 64)
				imgID, err2 := strconv.ParseInt(parts[3], 10, 64)
				if err1 == nil && err2 == nil {
					onExpire(domain.TenantID(tenantID), imgID)
				} else {
					zap.L().Error("Failed to parse tenantID or imgID from expired key",
						zap.String("key", msg.Payload),
						zap.Error(err1),
						zap.Error(err2),
					)
				}
			}
		}
	}
}

// RemoveFromDeleteQueue 从删除队列中移除指定图片ID
func (c *ImgRedisCache) RemoveFromDeleteQueue(tenantID domain.TenantID, imgID int64) error {
	key := c.buildDeletedQueueKey(tenantID, imgID)
	return c.client.Del(context.Background(), key).Err()
}
