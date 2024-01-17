package db

import (
	"log"

	_ "github.com/golang-migrate/migrate/database/postgres"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

// func ExecMigrations(db *sqlx.DB) error {
// 	m, err := migrate.New("file://utils/db/migrations", os.Getenv("DB_URL"))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
// 		log.Fatal(err)
// 	}
// 	return nil
// }

func ExecMigrations(db *sqlx.DB) {
    driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
    if err != nil {
        log.Fatal(err)
    }

    m, err := migrate.NewWithDatabaseInstance(
        "file://utils/db/migrations", 
        "postgres", driver)
    if err != nil {
        log.Fatal(err)
    }

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatal(err)
    }
}