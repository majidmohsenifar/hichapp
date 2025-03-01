package metrics

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
)

var (
	RedisTotalCommands = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "redis_commands_total",
			Help: "Total Redis commands executed",
		},
		[]string{"command", "status"},
	)

	RedisDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "redis_query_duration_seconds",
			Help:    "Redis query duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"command"},
	)

	RedisCacheHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "redis_cache_hit_total",
			Help: "Total cache hits",
		},
	)

	RedisCacheMisses = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "redis_cache_miss_total",
			Help: "Total cache misses",
		},
	)
)

type RedisMetricsHook struct{}

func (h *RedisMetricsHook) DialHook(next redis.DialHook) redis.DialHook {
	return next
}

func (h *RedisMetricsHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		start := time.Now()
		err := next(ctx, cmd)
		duration := time.Since(start).Seconds()

		RedisDuration.WithLabelValues(cmd.Name()).Observe(duration)
		status := "success"
		if cmd.Err() != nil {
			status = "error"
		}
		RedisTotalCommands.WithLabelValues(cmd.Name(), status).Inc()

		if cmd.Name() == "get" {
			if cmd.Err() == redis.Nil {
				RedisCacheMisses.Inc()
			} else if cmd.Err() == nil {
				RedisCacheHits.Inc()
			}
		}

		return err
	}
}

func (h *RedisMetricsHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		start := time.Now()
		err := next(ctx, cmds)
		duration := time.Since(start).Seconds()

		for _, cmd := range cmds {
			RedisDuration.WithLabelValues(cmd.Name()).Observe(duration)
			status := "success"
			if cmd.Err() != nil {
				status = "error"
			}
			RedisTotalCommands.WithLabelValues(cmd.Name(), status).Inc()
		}
		return err
	}
}
