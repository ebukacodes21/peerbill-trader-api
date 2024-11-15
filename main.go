package main

import (
	"database/sql"
	"log"
	"peerbill-trader-server/api"
	db "peerbill-trader-server/db/sqlc"
	"peerbill-trader-server/utils"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal(err)
	}

	repository := db.NewRepository(conn)
	server, err := api.NewServer(config, repository)
	if err != nil {
		log.Fatal(err)
	}

	err = server.StartServer(config.ServerAddr)
	if err != nil {
		log.Fatal(err)
	}
}
