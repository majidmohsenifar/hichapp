package metrics

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	DBTotalQueries = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_queries_total",
			Help: "Total number of DB queries executed",
		},
		[]string{"query_type", "status"},
	)

	DBQueryDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "DB query duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"query_type"},
	)
)

type DBTracer struct{}

// TraceQueryStart captures the SQL query string and stores the start time in the context.
func (d *DBTracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	// Store the SQL query and start time in the context
	ctx = context.WithValue(ctx, "sql_query", data.SQL)
	ctx = context.WithValue(ctx, "query_start", time.Now())
	return ctx
}

// TraceQueryEnd captures the duration of the query execution and logs the metrics.
func (d *DBTracer) TraceQueryEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	start, ok := ctx.Value("query_start").(time.Time)
	if !ok {
		return
	}

	// Get the SQL query from context
	sqlQuery, _ := ctx.Value("sql_query").(string)
	duration := time.Since(start).Seconds()

	status := "success"
	if data.Err != nil {
		status = "error"
	}

	// Log the query and status
	DBTotalQueries.WithLabelValues(sqlQuery, status).Inc()
	DBQueryDuration.WithLabelValues(sqlQuery).Observe(duration)
}
