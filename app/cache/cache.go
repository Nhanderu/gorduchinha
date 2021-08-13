package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/paemuri/gorduchinha/app/constant"
	"github.com/paemuri/gorduchinha/app/contract"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type redisCache struct {
	redis             *redis.Client
	prefix            string
	defaultExpiration time.Duration
}

func New(
	url string,
	db int,
	prefix string,
	defaultExpiration time.Duration,
) (contract.CacheManager, error) {

	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	opt.DB = db

	return redisCache{
		redis:             redis.NewClient(opt),
		prefix:            prefix,
		defaultExpiration: defaultExpiration,
	}, nil
}

func (r redisCache) buildKey(key string) string {
	return r.prefix + ":" + key
}

func (r redisCache) ClientPool() *redis.Client {
	return r.redis
}

func (r redisCache) Prefix() string {
	return r.prefix
}

func (r redisCache) Get(key string) ([]byte, error) {

	val, err := r.redis.Get(context.Background(), r.buildKey(key)).Bytes()
	if err == redis.Nil {
		return val, errors.WithStack(constant.NewErrorCacheMiss())
	}
	if err != nil {
		return val, errors.WithStack(err)
	}

	return val, nil
}

func (r redisCache) Set(key string, data []byte) error {

	err := r.redis.Set(context.Background(), r.buildKey(key), data, r.defaultExpiration).Err()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r redisCache) GetJSON(key string, data interface{}) error {

	val, err := r.Get(key)
	if err != nil {
		return errors.WithStack(err)
	}

	err = json.Unmarshal(val, &data)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r redisCache) SetJSON(key string, data interface{}) error {

	dataString, err := json.Marshal(data)
	if err != nil {
		return errors.WithStack(err)
	}

	err = r.Set(key, dataString)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r redisCache) GetExpiration(key string) (time.Duration, error) {

	expiration, err := r.redis.TTL(context.Background(), r.buildKey(key)).Result()
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return expiration, nil
}

func (r redisCache) SetExpiration(key string, expiration time.Duration) error {

	err := r.redis.Expire(context.Background(), r.buildKey(key), expiration).Err()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r redisCache) Invalidate(key string) error {

	err := r.redis.Del(context.Background(), r.buildKey(key)).Err()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r redisCache) CleanAll() error {

	keys, err := r.redis.Keys(context.Background(), r.buildKey("*")).Result()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return errors.WithStack(err)
	}

	if len(keys) > 0 {
		err = r.redis.Del(context.Background(), keys...).Err()
	}
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
