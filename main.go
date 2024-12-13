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

	db "github.com/ebukacodes21/peerbill-trader-api/db/sqlc"
	"github.com/ebukacodes21/peerbill-trader-api/mail"
	"github.com/ebukacodes21/peerbill-trader-api/socket"
	"github.com/ebukacodes21/peerbill-trader-api/token"

	"github.com/ebukacodes21/peerbill-trader-api/servers"
	"github.com/ebukacodes21/peerbill-trader-api/utils"
	"github.com/ebukacodes21/peerbill-trader-api/worker"
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

	repository := db.NewRepository(conn)
	servers.RunDBMigration(config.MigrationURL, config.DBSource)
	// message broker
	redisOption := asynq.RedisClientOpt{
		Addr: config.REDISServerAddr,
	}

	token, err := token.NewToken(config.TokenKey)
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), signals...)
	defer stop()
	group, ctx := errgroup.WithContext(ctx)

	manager := socket.NewWebSocketManager()
	mailer := mail.NewGmailSender(config.EmailSender, config.EmailAddress, config.EmailPassword)
	taskDistributor := worker.NewRedisTaskDistributor(redisOption)
	taskProcessor := worker.NewRedisTaskProcessor(config, redisOption, repository, mailer, token, manager)

	servers.RunGatewayServer(group, ctx, config, repository, taskDistributor, taskProcessor)
	servers.RunTaskProcessor(group, ctx, redisOption, repository, config, token, manager)
	servers.RunGrpcServer(group, ctx, config, repository, taskDistributor, taskProcessor)
	servers.RunWebSocketServer(group, ctx, config, manager)

	// wait bfr exiting main fn
	err = group.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
