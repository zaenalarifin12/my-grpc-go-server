package db

import (
	"database/sql"
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func Migrate(conn *sql.DB) {
	log.Println("Database migration start")

	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		log.Fatal("Error creating migration driver: ", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://db/migration", "postgres", driver)
	if err != nil {
		log.Fatal("Error creating migration instance: ", err)
	}

	if err := m.Down(); err != nil {
		log.Println("Database migration (down) failed: ", err)
	}

	if err := m.Up(); !errors.Is(err, migrate.ErrNoChange) && err != nil {
		log.Fatal("Database migration (up) failed: ", err)
	}

	log.Println("Database migration end")
}
