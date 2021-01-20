package registry

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
)

var redisDB *redis.Client

var ctx = context.Background()

func init() {
	redisDB = getRedisClient()
}

func getRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func (s *SchemaRegistry) setSchema() error {
	value, _ := json.Marshal(s)
	err := redisDB.Set(ctx, s.ServiceName, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func getSchema(service string) (*SchemaRegistry, error) {
	val2, err := redisDB.Get(ctx, service).Result()
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
	serviceSchemas, _, err := redisDB.Scan(ctx, cursor, "*", 100).Result()
	if err != nil {
		return &[]string{}, err
	}

	return &serviceSchemas, nil
}
