package worker

import (
	"context"
	db "peerbill-trader-server/db/sqlc"
	"peerbill-trader-server/mail"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server     *asynq.Server
	repository db.DatabaseContract
	mailer     mail.EmailSender
}

const (
	Critical = "critical"
	Default  = "default"
)

func NewRedisTaskProcessor(redisOpt asynq.RedisConnOpt, repository db.DatabaseContract, mailer mail.EmailSender) TaskProcessor {
	server := asynq.NewServer(redisOpt,
		asynq.Config{
			Queues: map[string]int{
				Critical: 10,
				Default:  5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Error().Err(err).Str("task_type", task.Type()).Bytes("payload", task.Payload()).Msg("task failed")
			}),
			// Logger: ,
		})
	return &RedisTaskProcessor{server: server, repository: repository, mailer: mailer}
}

func (rtp *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(send_verify_email_task, rtp.ProcessTaskSendVerifyEmail)

	return rtp.server.Start(mux)
}
