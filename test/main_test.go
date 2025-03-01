package test

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	infraredis "github.com/majidmohsenifar/hichapp/infra/redis"

	"github.com/caarlos0/env"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/majidmohsenifar/hichapp/cmd"
	"github.com/majidmohsenifar/hichapp/config"
	"github.com/majidmohsenifar/hichapp/infra/db"
	"github.com/majidmohsenifar/hichapp/repository"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type TestApp struct {
	address string
	db      *pgxpool.Pool
	repo    *repository.Queries
	redis   redis.UniversalClient
	cfg     *config.Config
}

func (app *TestApp) close() {
	app.db.Close()
	app.redis.Close()
}

func (app *TestApp) CreatePollWithOptionsAndTags(ctx context.Context, title string, options []string, tags []string) repository.Poll {
	poll, err := app.repo.CreatePoll(ctx, app.db, title)
	if err != nil {
		panic(err)
	}
	optionParams := make([]repository.CreateOptionParams, 0, len(options))
	for _, option := range options {
		optionParams = append(optionParams, repository.CreateOptionParams{
			PollID:  poll.ID,
			Content: option,
		})
	}
	_, err = app.repo.CreateOption(ctx, app.db, optionParams)
	if err != nil {
		panic(err)
	}

	tagIDs, err := app.repo.CreateTags(ctx, app.db, tags)
	if err != nil {
		panic(err)
	}

	createPollTagsParams := make([]repository.CreatePollTagParams, len(tags))
	for i, t := range tagIDs {
		createPollTagsParams[i] = repository.CreatePollTagParams{
			PollID: poll.ID,
			TagID:  t,
		}
	}
	_, err = app.repo.CreatePollTag(ctx, app.db, createPollTagsParams)
	if err != nil {
		panic(err)
	}
	return poll
}

func spawn_app() *TestApp {
	ctx := context.Background()
	cfg := config.Get()
	if err := env.Parse(cfg); err != nil {
		panic(err)
	}
	newDSN := configureDB(ctx, cfg.PostgresDSN)
	cfg.PostgresDSN = newDSN

	cfg.HttpAddress = "127.0.0.1:0"
	reg := prometheus.NewRegistry()

	httpApp, err := cmd.BuildHttpServer(ctx, cfg, reg)
	if err != nil {
		panic(err)
	}
	go httpApp.Run()
	db, err := db.NewDBClient(ctx, cfg.PostgresDSN)
	if err != nil {
		panic(err)
	}
	repo := repository.New()
	address := fmt.Sprintf("http://127.0.0.1:%d", httpApp.Port())
	redisClient, _ := infraredis.NewRedisClient(cfg.RedisDSN)
	return &TestApp{
		address: address,
		db:      db,
		repo:    repo,
		redis:   redisClient,
		cfg:     cfg,
	}
}

func configureDB(ctx context.Context, dbDSN string) string {
	dbName := uuid.New().String()
	dbName = "db_" + strings.ReplaceAll(dbName, "-", "_")
	dbURL, err := url.Parse(dbDSN)
	if err != nil {
		panic(err)
	}
	password, _ := dbURL.User.Password()
	dbURLWithoutDatabase := fmt.Sprintf("postgres://%s:%s@%s:%s?sslmode=disable", dbURL.User.Username(), password, dbURL.Hostname(), dbURL.Port())

	dbClient, err := db.NewDBClient(ctx, dbURLWithoutDatabase)
	if err != nil {
		panic(err)
	}
	_, err = dbClient.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s;", dbName))
	if err != nil {
		panic(err)
	}

	dbDSN = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbURL.User.Username(), password, dbURL.Hostname(), dbURL.Port(), dbName)

	dbClient, err = db.NewDBClient(ctx, dbDSN)
	if err != nil {
		panic(err)
	}

	//run migrate here
	m, err := migrate.New(
		"file://../db/migrations",
		dbDSN)
	if err != nil {
		panic(err)
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic(err)
	}
	return dbDSN
}
