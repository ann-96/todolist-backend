package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	rediscli "github.com/go-redis/redis/v8"
)

type SessionCache interface {
	GetSession(string) *int
	CreateSession(string, int) error
	DeleteSession(string) error
}

type sessionCache struct {
	client *rediscli.Client
	ctx    context.Context
}

func NewSessionCache(address string, ctx context.Context) *sessionCache {
	rdb := rediscli.NewClient(&rediscli.Options{
		Addr: address,
		DB:   0,
	})

	return &sessionCache{
		client: rdb,
		ctx:    ctx,
	}
}

func (r *sessionCache) GetSession(key string) *int {
	key = "userid_" + key
	raw, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return nil
	}

	res, err := strconv.Atoi(raw)
	if err != nil {
		return nil
	}

	return &res
}

func (r *sessionCache) CreateSession(key string, id int) error {
	validFor := 30 * time.Hour * 24 // The authorization session has the TTL of 30 days
	key = "userid_" + key
	val := fmt.Sprintf("%v", id)
	err := r.client.Set(r.ctx, key, val, validFor).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *sessionCache) DeleteSession(key string) error {
	key = "userid_" + key
	err := r.client.Del(r.ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}
