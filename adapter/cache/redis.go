package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client   *redis.Client
	ParseURL bool
}

var _ Storage = (*Redis)(nil)

func (r *Redis) Load(cfg *Config) Storage {
	opts := &redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DBID,
	}
	var err error
	if r.ParseURL {
		opts, err = redis.ParseURL(cfg.URL)
		log.Print("error while connecting redis", err)
	}
	r.client = redis.NewClient(opts)

	return r
}

func (r *Redis) Get(ctx context.Context, key string) (Record, error) {
	val, err := r.client.Get(ctx, key).Result()
	rec := Record{}
	rec.Value = val
	rec.Key = key

	return rec, err
}

func (r *Redis) Set(ctx context.Context, rec Record) (bool, error) {
	val := reflect.ValueOf(rec.Value)
	switch val.Kind() {
	case reflect.Struct:
		marshaledValue, err := json.Marshal(rec.Value)
		if err != nil {
			return false, fmt.Errorf("failed to marshal value: %w", err)
		}
		rec.Value = string(marshaledValue)
	default:
	}
	err := r.client.Set(ctx, rec.Key, rec.Value, rec.Expiry).Err()

	return err == nil, err
}
