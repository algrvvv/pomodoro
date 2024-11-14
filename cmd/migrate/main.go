package main

import (
	"database/sql"
	"errors"
	"flag"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/algrvvv/pomodoro/internal/config"
	"github.com/algrvvv/pomodoro/internal/database"
)

var (
	down = flag.Bool("down", false, "down migrations")
	up   = flag.Bool("up", false, "down migrations")
)

func main() {
	flag.Parse()
	if (*up && *down) || (!*up && !*down) {
		log.Fatalf("invalid flags")
	}

	if err := config.Parse("config.yml"); err != nil {
		log.Fatal("failed to load config")
	}

	connString := database.GetConnectionString()
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatalf("failed to open database: %s", err)
	}
	defer db.Close()

	log.Println("successfully connected to database")

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("failed to initialize driver: %s", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations/",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("failed to initialize migrations: %s", err)
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			log.Println("failed to close src migrations: %v", srcErr)
		}
		if dbErr != nil {
			log.Println("failed to close db migrations: %v", dbErr)
		}

		if srcErr == nil && dbErr == nil {
			log.Println("successfully closed migrations")
		}
	}()

	if v, dirty, err := m.Version(); err == nil && dirty {
		log.Println("last version: %v is dirty", v)

		if err = m.Force(int(v)); err != nil {
			log.Fatalf("failed to force: %v", err)
		}

		if err = m.Steps(-1); err != nil {
			log.Fatalf("failed to rollback migrations: %v", err)
		}

		log.Println("successfully cleaning dirty migrations")
	} else if err != nil {
		if !errors.Is(err, migrate.ErrNilVersion) {
			log.Fatalf("failed to get mirgations version: %v", err)
		}
	}

	if *up {
		log.Println("start up migrations")
		if err = m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				log.Println("Migrations have not changes")
			} else {
				log.Fatalf("failed to run migrations: %v", err)
			}
		} else {
			log.Println("successfully applied migrations")
		}
	} else if *down {
		log.Println("start down migrations")
		if err = m.Down(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				log.Println("Migrations have not changes")
			} else {
				log.Fatalf("failed to run migrations: %v", err)
			}
		} else {
			log.Println("successfully applied migrations")
		}
	}
}
