package servers

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// RunDBMigration runs the migration process for the given source and database URL
func RunDBMigration(source string, url string) {
	migration, err := migrate.New(source, url)
	if err != nil {
		log.Fatal("Unable to start migration:", err)
	}

	// Run migrations
	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Error running migration:", err)
	}

	log.Println("Migration completed successfully")
}
