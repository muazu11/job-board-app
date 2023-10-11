package db

import (
	"context"
	_ "embed"

	"github.com/z0ne-dev/mgx"
)

var (
	//go:embed migrations/init.sql
	initMg string
)

func (d DB) migrate() error {
	migrator, err := mgx.New(mgx.Migrations(
		mgx.NewRawMigration("init", initMg),
	))
	if err != nil {
		return err
	}
	conn, err := d.pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	err = migrator.Migrate(context.Background(), conn)
	if err != nil {
		return err
	}
	return nil
}
