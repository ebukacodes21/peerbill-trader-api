package db

import (
	"database/sql"
	"log"
	"os"
	"project-server/utils"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries

func TestMain(M *testing.M) {
	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal(err)
	}
	testQueries = New(conn)
	os.Exit(M.Run())
}
