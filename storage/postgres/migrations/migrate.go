package migrations

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // for migration data
	_ "github.com/lib/pq"                                // for migration driver
)

// 	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
//		s.cfg.User, s.cfg.Password, s.cfg.Host, s.cfg.Port, s.cfg.Name)
func Migrate(connString, command, migrationsDir string) error {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return err
	}
	operation := func() error {
		return db.Ping()
	}
	bOff := backoff.NewExponentialBackOff()
	bOff.MaxElapsedTime = time.Second * 30
	err = backoff.Retry(operation, bOff)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %v", err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+filepath.Join(migrationsDir, "migrations"),
		"user_management_service_test", driver)
	if err != nil {
		return err
	}
	switch command {
	case "up":
		return m.Up()
	case "down":
		err = m.Down()
		if err == migrate.ErrNoChange {
			return nil
		}
		return err
	default:
		return errors.New("unknown command")
	}
}
