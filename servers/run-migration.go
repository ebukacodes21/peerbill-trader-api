package servers

import (
	"log"

	"github.com/golang-migrate/migrate"
)

func RunDBMigration(url string, source string) {
	migration, err := migrate.New(url, source)
	if err != nil {
		log.Fatal(err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	log.Print("migration successful")
}
