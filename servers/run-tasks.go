package servers

import (
	"context"
	"log"
	db "peerbill-trader-api/db/sqlc"
	"peerbill-trader-api/mail"
	"peerbill-trader-api/utils"
	"peerbill-trader-api/worker"

	"github.com/hibiken/asynq"
	"golang.org/x/sync/errgroup"
)

func RunTaskProcessor(group *errgroup.Group, ctx context.Context, options asynq.RedisClientOpt, repository db.DatabaseContract, config utils.Config) {
	log.Print("running processor")
	mailer := mail.NewGmailSender(config.EmailSender, config.EmailAddress, config.EmailPassword)
	taskProcessor := worker.NewRedisTaskProcessor(options, repository, mailer)

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
