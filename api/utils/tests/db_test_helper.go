package testutils

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	pgMigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/SaeedAlian/econest/api/config"
)

const (
	testDBName = "econest_test"
)

func SetupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	conninfo := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		config.Env.DBUser,
		config.Env.DBPassword,
		config.Env.DBHost,
		config.Env.DBPort,
		testDBName,
	)

	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping test database: %v", err)
	}

	driver, err := pgMigrate.WithInstance(db, &pgMigrate.Config{})
	if err != nil {
		t.Fatalf("Failed to create migration driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../../db/migrate/migrations", // Adjust path as needed
		"postgres",
		driver,
	)
	if err != nil {
		t.Fatalf("Failed to create migration instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("Failed to apply migrations: %v", err)
	}

	t.Cleanup(func() {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			t.Logf("Warning: failed to rollback migrations: %v", err)
		}
		if err := db.Close(); err != nil {
			t.Logf("Warning: failed to close database connection: %v", err)
		}
	})

	return db
}
