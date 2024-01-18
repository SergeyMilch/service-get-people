package db

import (
	"log"

	_ "github.com/golang-migrate/migrate/database/postgres"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

func ExecMigrations(db *sqlx.DB) {
    driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
    if err != nil {
        log.Fatal("Error creating database driver instance:", err)
    }

    m, err := migrate.NewWithDatabaseInstance(
        "file://utils/db/migrations", 
        "postgres", driver)
    if err != nil {
        log.Fatal("Error creating migration instance:", err)
    }

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatal("Error applying migrations:", err)
    }

    log.Println("Migrations applied successfully")
}