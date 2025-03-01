package cmd

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/caarlos0/env"
	"github.com/go-playground/validator/v10"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/majidmohsenifar/hichapp/config"
	"github.com/majidmohsenifar/hichapp/handler/api"
	"github.com/majidmohsenifar/hichapp/handler/api/router"
	"github.com/majidmohsenifar/hichapp/infra/db"
	"github.com/majidmohsenifar/hichapp/infra/logger"
	"github.com/majidmohsenifar/hichapp/infra/redis"
	"github.com/majidmohsenifar/hichapp/repository"
	"github.com/majidmohsenifar/hichapp/service/limiter"
	"github.com/majidmohsenifar/hichapp/service/poll"
	"github.com/majidmohsenifar/hichapp/service/statistic"
	"github.com/majidmohsenifar/hichapp/service/tag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

type HttpApp struct {
	httpServer *http.Server
	listener   net.Listener
	port       int
	closeFuncs []func() error
}

func (app *HttpApp) Run() {
	app.httpServer.Serve(app.listener)
}

func (app *HttpApp) Port() int {
	return app.port
}

func RunHttpServer() {
	ctx := context.Background()
	cfg := config.Get()
	if err := env.Parse(cfg); err != nil {
		slog.Error("cannot parse env", "err", err)
		os.Exit(1)
	}
	logger := logger.NewLogger()
	slog.SetDefault(logger)
	//TODO: handle prom later
	reg := prometheus.NewRegistry()
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)
	httpApp, err := BuildHttpServer(ctx, cfg, reg)
	if err != nil {
		slog.Error("cannot build http server", "err", err)
		os.Exit(1)
	}
	httpApp.Run()
	//closing resources
	for _, f := range httpApp.closeFuncs {
		f()
	}
}

func BuildHttpServer(ctx context.Context, cfg *config.Config, reg *prometheus.Registry) (*HttpApp, error) {
	logger := logger.NewLogger()
	slog.SetDefault(logger)
	db, err := db.NewDBClient(ctx, cfg.PostgresDSN)
	if err != nil {
		slog.Error("cannot create db client", "err", err)
		os.Exit(1)
	}
	closeFuncs := []func() error{
		func() error {
			db.Close()
			return nil
		},
	}

	//run migrate here
	m, err := migrate.New(
		"file://../db/migrations",
		cfg.PostgresDSN)
	if err != nil {
		slog.Error("cannot create migrations", "err", err)
		os.Exit(1)
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error("cannot run migrations", "err", err)
		os.Exit(1)
	}

	repo := repository.New()
	redisClient, err := redis.NewRedisClient(cfg.RedisDSN)
	closeFuncs = append(closeFuncs, redisClient.Close)
	tagService := tag.New(db, repo, redisClient)
	userVoteLimiter := limiter.NewUserVoteLimiter(redisClient)
	validator := validator.New()
	pollService := poll.New(db, repo, tagService, userVoteLimiter)
	pollHandler := api.NewPollHandler(pollService, validator)
	statsService := statistic.New(db, repo)
	statsHandler := api.NewStatsHandler(statsService, validator)
	router := router.New(
		pollHandler,
		statsHandler,
		cfg,
		reg,
	)

	ln, err := net.Listen("tcp", cfg.HttpAddress)
	if err != nil {
		return nil, err
	}

	httpServer := &http.Server{
		Addr:    ln.Addr().String(),
		Handler: router,
	}

	closeFuncs = append(closeFuncs, ln.Close)
	port := ln.Addr().(*net.TCPAddr).Port

	return &HttpApp{
		httpServer: httpServer,
		listener:   ln,
		port:       port,
		closeFuncs: closeFuncs,
	}, nil

}
