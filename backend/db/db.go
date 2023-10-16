package db

import (
	"context"
	"fmt"
	"jobboard/backend/services"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sanggonlee/gosq"
)

type Config struct {
	User     string
	Password string
	Port     int
	Host     string
}

type DB struct {
	pool *pgxpool.Pool
}

func New(config Config, servicesConfig services.Config) DB {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/postgres",
		config.User, config.Password, config.Host, config.Port,
	)
	pgxConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		panic(err)
	}
	pgxConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe
	pool, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		panic(err)
	}

	db := DB{pool: pool}
	err = db.migrate(servicesConfig)
	if err != nil {
		panic(err)
	}

	return db
}

func (d DB) Query(ctx context.Context, dest any, template string, templateArgs any, args ...any) error {
	sql, err := gosq.Compile(template, templateArgs)
	if err != nil {
		return err
	}
	return pgxscan.Select(ctx, d.pool, dest, sql, args...)
}

func (d DB) QueryOne(ctx context.Context, dest any, template string, templateArgs any, args ...any) error {
	sql, err := gosq.Compile(template, templateArgs)
	if err != nil {
		return err
	}
	return pgxscan.Get(ctx, d.pool, dest, sql, args...)
}

func (d DB) Exec(ctx context.Context, template string, templateArgs any, args ...any) error {
	sql, err := gosq.Compile(template, templateArgs)
	if err != nil {
		return err
	}
	_, err = d.pool.Exec(ctx, sql, args...)
	return err
}
