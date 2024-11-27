package main

import (
	"context"

	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hibiken/asynq"
	"golang.org/x/sync/errgroup"

	db "peerbill-trader-api/db/sqlc"

	"peerbill-trader-api/servers"
	"peerbill-trader-api/utils"
	"peerbill-trader-api/worker"
)

/*
*
signals to look out for a graceful shutdown
*/
var signals = []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGINT}

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), signals...)
	defer stop()

	repository := db.NewRepository(conn)
	servers.RunDBMigration(config.MigrationURL, config.DBSource)
	// message broker
	redisOption := asynq.RedisClientOpt{
		Addr: config.REDISServerAddr,
	}

	group, ctx := errgroup.WithContext(ctx)

	taskDistributor := worker.NewRedisTaskDistributor(redisOption)
	servers.RunGatewayServer(group, ctx, config, repository, taskDistributor)
	servers.RunTaskProcessor(group, ctx, redisOption, repository, config)
	servers.RunGrpcServer(group, ctx, config, repository, taskDistributor)
	servers.RunWebSocketServer(group, ctx, config)

	// wait bfr exiting main fn
	err = group.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
