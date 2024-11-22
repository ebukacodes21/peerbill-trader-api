package worker

import (
	"context"
	db "peerbill-trader-server/db/sqlc"

	"github.com/hibiken/asynq"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server     *asynq.Server
	repository db.DatabaseContract
}

const (
	Critical = "critical"
	Default  = "default"
)

func NewRedisTaskProcessor(redisOpt asynq.RedisConnOpt, repository db.DatabaseContract) TaskProcessor {
	server := asynq.NewServer(redisOpt, asynq.Config{
		Queues: map[string]int{
			Critical: 10,
			Default:  5,
		},
	})
	return &RedisTaskProcessor{server: server, repository: repository}
}

func (rtp *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(send_verify_email_task, rtp.ProcessTaskSendVerifyEmail)

	return rtp.server.Start(mux)
}
