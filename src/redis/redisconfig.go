package redis


import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisService represents a Redis service.
type RedisService struct {
	client *redis.Client
}

// NewRedisService creates a new instance of RedisService.
func NewRedisService(addr, password string, db int) (*RedisService, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Test the connection
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &RedisService{client: client}, nil
}

// Close closes the Redis connection.
func (r *RedisService) Close() error {
	return r.client.Close()
}

// Set stores a key-value pair in Redis.
func (r *RedisService) Set(key string, data interface{}) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return r.client.Set(context.Background(), key, jsonData, 0).Err()
}

// SetEx stores a key-value pair in Redis with an expiration time.
func (r *RedisService) SetEx(key string, data interface{}, duration time.Duration) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return r.client.Set(context.Background(), key, jsonData, duration).Err()
}

// Get retrieves the value associated with the given key from Redis.
func (r *RedisService) Get(key string, target interface{}) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}

	jsonData, err := r.client.Get(context.Background(), key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(jsonData), target)
}

// Delete removes the specified key from Redis.
func (r *RedisService) Delete(key string) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}

	return r.client.Del(context.Background(), key).Err()
}

// Example usage:
func main() {
	// Create a new RedisService instance
	redisService, err := NewRedisService("localhost:6379", "", 0)
	if err != nil {
		fmt.Println("Error creating RedisService:", err)
		return
	}
	defer redisService.Close()

	// Example usage
	key := "example_key"
	data := map[string]interface{}{"field1": "value1", "field2": 42}

	// Set data in Redis
	err = redisService.Set(key, data)
	if err != nil {
		fmt.Println("Error setting key:", err)
		return
	}

	// Get data from Redis
	var retrievedData map[string]interface{}
	err = redisService.Get(key, &retrievedData)
	if err != nil {
		fmt.Println("Error getting value:", err)
		return
	}

	fmt.Println("Retrieved data:", retrievedData)
}
