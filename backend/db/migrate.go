package db

import (
	"context"
	_ "embed"
	"jobboard/backend/services"
	"strings"

	"github.com/z0ne-dev/mgx"
	"golang.org/x/crypto/bcrypt"
)

var (
	//go:embed migrations/init.sql
	initMg string
	//go:embed migrations/admin_user.sql
	adminUserMg string
)

func (d DB) migrate(config services.Config) error {
	err := preprocessMigrations(config)
	if err != nil {
		return err
	}
	migrator, err := mgx.New(mgx.Migrations(
		mgx.NewRawMigration("init", initMg),
		mgx.NewRawMigration("adminUser", adminUserMg),
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

func preprocessMigrations(config services.Config) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(config.AdminPassword), 10)
	if err != nil {
		return err
	}
	replacer := strings.NewReplacer(
		"__ADMIN_EMAIL__", config.AdminEmail,
		"__ADMIN_PASSWORD_HASH__", string(passwordHash),
	)
	adminUserMg = replacer.Replace(adminUserMg)

	return nil
}
