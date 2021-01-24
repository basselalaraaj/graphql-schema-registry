package registry

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	// RedisDB client for redis
	RedisDB redisClient
	ctx     = context.Background()
)

type redisClientInterface interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd
}

type redisClient struct {
	client redisClientInterface
}

func (r *redisClient) CreateRedisClient() {
	r.client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func (s *SchemaRegistry) setSchema() error {
	value, _ := json.Marshal(s)
	err := RedisDB.client.Set(ctx, s.ServiceName, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func getSchema(service string) (*SchemaRegistry, error) {
	val2, err := RedisDB.client.Get(ctx, service).Result()
	if err == redis.Nil {
		return &SchemaRegistry{}, err
	} else if err != nil {
		return &SchemaRegistry{}, err
	} else {
		var serviceSchema SchemaRegistry
		err := json.Unmarshal([]byte(val2), &serviceSchema)
		if err != nil {
			return &SchemaRegistry{}, err
		}
		return &serviceSchema, nil
	}
}

func scanSchemas() (*[]string, error) {
	var cursor uint64
	serviceSchemas, _, err := RedisDB.client.Scan(ctx, cursor, "*", 100).Result()
	if err != nil {
		return &[]string{}, err
	}

	return &serviceSchemas, nil
}
