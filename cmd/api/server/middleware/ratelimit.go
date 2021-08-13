package middleware

import (
	"time"

	"github.com/paemuri/gorduchinha/app/contract"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/fasthttp"
	"github.com/ulule/limiter/v3/drivers/store/redis"
)

func RateLimit(
	cache contract.CacheManager,
	period time.Duration,
	limit int64,
) RequestMiddleware {

	var storeOptions limiter.StoreOptions
	storeOptions.Prefix = cache.Prefix() + ":server-rate-limit"
	store, err := redis.NewStoreWithOptions(cache.ClientPool(), storeOptions)
	if err != nil {
		panic(err)
	}

	var rate limiter.Rate
	rate.Period = period
	rate.Limit = limit
	mw := fasthttp.NewMiddleware(limiter.New(store, rate))

	return mw.Handle
}
