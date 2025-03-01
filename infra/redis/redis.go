package redis

import (
	"github.com/majidmohsenifar/hichapp/infra/metrics"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(redisAddress string) (redis.UniversalClient, error) {
	opts, err := redis.ParseURL(redisAddress)
	if err != nil {
		return nil, err
	}
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{
			opts.Addr,
		},
		ClientName:            opts.ClientName,
		DB:                    opts.DB,
		Username:              opts.Username,
		Password:              opts.Password,
		SentinelUsername:      opts.Username,
		SentinelPassword:      opts.Password,
		MaxRetries:            opts.MaxRetries,
		MinRetryBackoff:       opts.MinRetryBackoff,
		MaxRetryBackoff:       opts.MaxRetryBackoff,
		DialTimeout:           opts.DialTimeout,
		ReadTimeout:           opts.ReadTimeout,
		WriteTimeout:          opts.WriteTimeout,
		ContextTimeoutEnabled: opts.ContextTimeoutEnabled,
		PoolFIFO:              opts.PoolFIFO,
		PoolSize:              opts.PoolSize,
		PoolTimeout:           opts.PoolTimeout,
		MinIdleConns:          opts.MinIdleConns,
		MaxIdleConns:          opts.MaxIdleConns,
		ConnMaxIdleTime:       opts.ConnMaxIdleTime,
		ConnMaxLifetime:       opts.ConnMaxLifetime,
		MaxRedirects:          opts.MaxRetries,
		ReadOnly:              false,
		RouteByLatency:        false,
		RouteRandomly:         false,
	})

	rdb.AddHook(&metrics.RedisMetricsHook{})
	return rdb, nil

}
