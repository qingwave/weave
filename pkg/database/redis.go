package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/qingwave/weave/pkg/config"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

var (
	RedisDisableError = errors.New("redis disable")
)

type RedisDB struct {
	enable bool
	*redis.Client
}

func NewRedisClient(conf *config.RedisConfig) (*RedisDB, error) {
	if !conf.Enable {
		logrus.Info("redis disable")
		return &RedisDB{}, nil
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password: conf.Password,
		DB:       0,
	})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}

	return &RedisDB{
		enable: true,
		Client: rdb,
	}, nil
}

func (rdb *RedisDB) Endable() bool {
	return rdb.enable
}

func (rdb *RedisDB) HGet(key, field string, obj interface{}) error {
	if !rdb.enable {
		return RedisDisableError
	}

	return rdb.Client.HGet(context.Background(), key, field).Scan(obj)
}

func (rdb *RedisDB) HSet(key, field string, val interface{}) error {
	if !rdb.enable {
		return nil
	}

	return rdb.Client.HSet(context.Background(), key, field, val).Err()
}

func (rdb *RedisDB) HDel(key string, fields ...string) error {
	if !rdb.enable {
		return nil
	}

	return rdb.Client.HDel(context.Background(), key, fields...).Err()
}
