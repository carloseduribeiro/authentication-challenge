package db

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

func RunMigrations(databaseURL string) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		log.Fatalf("error open connection to apply migration: %s", err)
	}
	defer db.Close()
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("could not init driver: %s", err)
	}
	defer driver.Close()
	m, err := migrate.NewWithDatabaseInstance("file://../internal/infra/database/migrations", "pgx", driver)
	if err != nil {
		log.Fatalf("could not apply the migration: %s", err)
	}
	m.Up()
	defer m.Close()
}
