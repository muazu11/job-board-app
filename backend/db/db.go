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

func (d DB) Query(ctx context.Context, dest any, template string, templateArgs, args any) error {
	sql, err := gosq.Compile(template, templateArgs)
	if err != nil {
		return err
	}
	return pgxscan.Select(ctx, d.pool, dest, sql, args)
}

func (d DB) Exec(ctx context.Context, template string, templateArgs, args any) error {
	sql, err := gosq.Compile(template, templateArgs)
	if err != nil {
		return err
	}
	_, err = d.pool.Exec(ctx, sql, args)
	return err
}

func GetAll[T any](ctx context.Context, db DB, tableName string) ([]T, error) {
	var dest []T
	args := map[string]string{"table": tableName}
	err := db.Query(ctx, &dest, "SELECT * FROM .table", args)
	return dest, err
}

func GetById[T any](ctx context.Context, db DB, tableName string, id int) (T, error) {
	var dest T
	args := map[string]any{"table": tableName, "id": id}
	err := db.Query(ctx, &dest, "SELECT * FROM .table WHERE id = .id", args)
	return dest, err
}
func DeleteById(ctx context.Context, db DB, tableName string, id int) error {
	args := map[string]any{"table": tableName, "id": id}
	return db.Exec(ctx, "DELETE FROM .table WHERE id = .id", args)
}
