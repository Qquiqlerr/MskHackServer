package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var storagePath, migrationsPath, direction string
	flag.StringVar(&storagePath, "storage-path", "", "Path to a directory containing the migration files")
	flag.StringVar(&migrationsPath, "migrations-path", "", "Path to a directory containing the migration files")
	flag.StringVar(&direction, "direction", "", "Direction to migrate up")
	flag.Parse()
	if direction == "" {
		panic(errors.New("direction is required"))
	}
	if storagePath == "" {
		panic("storage-path is required")
	}
	if migrationsPath == "" {
		panic("migrations-path is required")
	}
	m, err := migrate.New("file://"+migrationsPath, storagePath)
	if err != nil {
		panic(err)
	}
	if direction == "up" {
		if err := m.Up(); err != nil {
			if !errors.Is(err, migrate.ErrNoChange) {
				panic(err)
			}
		}
	}
	if direction == "down" {
		if err := m.Down(); err != nil {
			if !errors.Is(err, migrate.ErrNoChange) {
				panic(err)
			}
		}
	}

	fmt.Println("Migrations applied")
}
