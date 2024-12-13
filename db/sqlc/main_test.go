package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/ebukacodes21/peerbill-trader-api/utils"

	_ "github.com/lib/pq"
)

var testDb *sql.DB
var testQueries *Queries

func TestMain(M *testing.M) {
	var err error
	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal(err)
	}

	testDb, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal(err)
	}

	testQueries = New(testDb)
	os.Exit(M.Run())
}
