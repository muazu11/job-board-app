package db

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
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

func New(config Config) DB {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/postgres",
		config.User, config.Password, config.Host, config.Port,
	)
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		panic(err)
	}

	db := DB{pool: pool}
	err = db.migrate()
	if err != nil {
		panic(err)
	}

	return db
}

func (d DB) Query(ctx context.Context, dest any, template string, args any) error {
	sql, err := gosq.Compile(template, args)
	if err != nil {
		return err
	}
	return pgxscan.Select(ctx, d.pool, dest, sql)
}

func (d DB) Exec(ctx context.Context, template string, args any) error {
	sql, err := gosq.Compile(template, args)
	if err != nil {
		return err
	}
	_, err = d.pool.Exec(ctx, sql)
	return err
}
