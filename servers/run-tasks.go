package servers

import (
	"context"
	"log"

	db "github.com/ebukacodes21/peerbill-trader-api/db/sqlc"
	"github.com/ebukacodes21/peerbill-trader-api/mail"
	"github.com/ebukacodes21/peerbill-trader-api/socket"
	"github.com/ebukacodes21/peerbill-trader-api/token"
	"github.com/ebukacodes21/peerbill-trader-api/utils"
	"github.com/ebukacodes21/peerbill-trader-api/worker"

	"github.com/hibiken/asynq"
	"golang.org/x/sync/errgroup"
)

func RunTaskProcessor(group *errgroup.Group, ctx context.Context, options asynq.RedisClientOpt, repository db.DatabaseContract, config utils.Config, token token.TokenMaker, manager *socket.WebSocketManager) {
	log.Print("running processor")
	mailer := mail.NewGmailSender(config.EmailSender, config.EmailAddress, config.EmailPassword)
	taskProcessor := worker.NewRedisTaskProcessor(config, options, repository, mailer, token, manager)

	err := taskProcessor.Start()
	if err != nil {
		log.Fatal("failed to process tasks")
	}

	group.Go(func() error {
		<-ctx.Done()
		log.Print("gracefully shutting down...")

		taskProcessor.Shutdown()
		log.Print("task processor shutdown.. goodbye")

		return nil
	})
}
