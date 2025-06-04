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
	var storageHost, storagePort, migrationsPath, storageUser, storageUserPass string

	flag.StringVar(&storageHost, "storage-host", "", "storage host")
	flag.StringVar(&storagePort, "storage-port", "", "storage port")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&storageUser, "storage-user", "", "storage user")
	flag.StringVar(&storageUserPass, "storage-user-pass", "", "storage user password")
	flag.Parse()

	if storageHost == "" {
		panic("storage-host is required")
	}
	if storagePort == "" {
		panic("storage-port is required")
	}
	if migrationsPath == "" {
		panic("migrations-path is required")
	}
	if storageUser == "" {
		panic("storage-user is required")
	}
	if storageUserPass == "" {
		panic("storage-user-pass is required")
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprint("postgres://"+storageUser+":"+storageUserPass+"@"+storageHost+":"+storagePort+"/sso?sslmode=disable"),
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")

			return
		}
		panic(err)
	}

	fmt.Println("migrations applied successfully")
}
