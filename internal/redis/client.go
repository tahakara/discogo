package redisclient

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Client defines the operations that a Redis client can perform.
type Client interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte, expiration time.Duration) error
	Delete(key string) error
	Add(key string, value []byte, expiration time.Duration) error
	Replace(key string, value []byte, expiration time.Duration) error
	Increment(key string, delta int64) (int64, error)
	Decrement(key string, delta int64) (int64, error)
	FlushAll() error
	Ping() error
	Close() error
	FindKeys(pattern string) ([]string, error)
}

type client struct {
	rdb *redis.Client
	ctx context.Context
}

// New creates a new Redis client.
func New(addr string, password string, db int) Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &client{
		rdb: rdb,
		ctx: context.Background(),
	}
}

func (c *client) Get(key string) ([]byte, error) {
	val, err := c.rdb.Get(c.ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	return val, err
}

func (c *client) Set(key string, value []byte, expiration time.Duration) error {
	return c.rdb.Set(c.ctx, key, value, expiration).Err()
}

func (c *client) Delete(key string) error {
	return c.rdb.Del(c.ctx, key).Err()
}

func (c *client) Add(key string, value []byte, expiration time.Duration) error {
	// NX: Only set the key if it does not already exist
	return c.rdb.SetNX(c.ctx, key, value, expiration).Err()
}

func (c *client) Replace(key string, value []byte, expiration time.Duration) error {
	// XX: Only set the key if it already exists
	return c.rdb.SetXX(c.ctx, key, value, expiration).Err()
}

func (c *client) Increment(key string, delta int64) (int64, error) {
	return c.rdb.IncrBy(c.ctx, key, delta).Result()
}

func (c *client) Decrement(key string, delta int64) (int64, error) {
	return c.rdb.DecrBy(c.ctx, key, delta).Result()
}

func (c *client) FlushAll() error {
	return c.rdb.FlushAll(c.ctx).Err()
}

func (c *client) Ping() error {
	return c.rdb.Ping(c.ctx).Err()
}

func (c *client) Close() error {
	return c.rdb.Close()
}

// FindKeys returns all keys matching the given pattern (use with care in production).
func (c *client) FindKeys(pattern string) ([]string, error) {
	var (
		cursor uint64
		keys   []string
	)
	for {
		var scanKeys []string
		var err error
		scanKeys, cursor, err = c.rdb.Scan(c.ctx, cursor, pattern, 100).Result()
		if err != nil {
			return nil, err
		}
		keys = append(keys, scanKeys...)
		if cursor == 0 {
			break
		}
	}
	return keys, nil
}
