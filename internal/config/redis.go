package config

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	// RedisClient 全局Redis客户端
	RedisClient *redis.Client
)

// InitRedis 初始化Redis连接
func InitRedis(config *RedisConfig) error {
	client := redis.NewClient(&redis.Options{
		Addr:     config.GetRedisAddr(),
		Password: config.Password,
		DB:       config.DB,
		PoolSize: 10,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("连接Redis失败: %w", err)
	}

	RedisClient = client
	return nil
}

// CloseRedis 关闭Redis连接
func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}

// GetRedisClient 获取Redis客户端
func GetRedisClient() *redis.Client {
	return RedisClient
}

// SetCache 设置缓存
func SetCache(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if RedisClient == nil {
		return fmt.Errorf("Redis客户端未初始化")
	}
	
	return RedisClient.Set(ctx, key, value, expiration).Err()
}

// GetCache 获取缓存
func GetCache(ctx context.Context, key string) (string, error) {
	if RedisClient == nil {
		return "", fmt.Errorf("Redis客户端未初始化")
	}
	
	return RedisClient.Get(ctx, key).Result()
}

// DeleteCache 删除缓存
func DeleteCache(ctx context.Context, key string) error {
	if RedisClient == nil {
		return fmt.Errorf("Redis客户端未初始化")
	}
	
	return RedisClient.Del(ctx, key).Err()
}

// ExistsCache 检查缓存是否存在
func ExistsCache(ctx context.Context, key string) (bool, error) {
	if RedisClient == nil {
		return false, fmt.Errorf("Redis客户端未初始化")
	}
	
	result, err := RedisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	
	return result > 0, nil
} 