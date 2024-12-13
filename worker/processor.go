package worker

import (
	"context"

	db "github.com/ebukacodes21/peerbill-trader-api/db/sqlc"
	"github.com/ebukacodes21/peerbill-trader-api/mail"
	"github.com/ebukacodes21/peerbill-trader-api/socket"
	"github.com/ebukacodes21/peerbill-trader-api/token"
	"github.com/ebukacodes21/peerbill-trader-api/utils"

	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type TaskProcessor interface {
	Start() error
	Shutdown()
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
	ProcessTaskSendForgotEmail(ctx context.Context, task *asynq.Task) error
	ProcessTaskOrderEmail(ctx context.Context, task *asynq.Task) error
	ProcessTaskUpdateOrder(ctx context.Context, task *asynq.Task) error
	ProcessTaskUpdateOrders(ctx context.Context, task *asynq.Task) error
	Get(ctx context.Context, value string) (string, error)
	Set(key string, value string) error
}

type RedisTaskProcessor struct {
	server     *asynq.Server
	redis      *redis.Client
	repository db.DatabaseContract
	mailer     mail.EmailSender
	config     utils.Config
	token      token.TokenMaker
	wsManager  *socket.WebSocketManager
}

const (
	Critical = "critical"
	Default  = "default"
)

func NewRedisTaskProcessor(config utils.Config, redisOpt asynq.RedisConnOpt, repository db.DatabaseContract, mailer mail.EmailSender, token token.TokenMaker, wsManager *socket.WebSocketManager) TaskProcessor {
	redisCache := redis.NewClient(&redis.Options{
		Addr: config.REDISServerAddr,
	})

	server := asynq.NewServer(redisOpt,
		asynq.Config{
			Queues: map[string]int{
				Critical: 10,
				Default:  5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Error().Err(err).Str("task_type", task.Type()).Bytes("payload", task.Payload()).Msg("task failed")
			}),
			Logger: NewLogger(),
		})
	return &RedisTaskProcessor{server: server, repository: repository, mailer: mailer, redis: redisCache, config: config, token: token, wsManager: wsManager}
}

func (rtp *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(send_verify_email_task, rtp.ProcessTaskSendVerifyEmail)
	mux.HandleFunc(send_forgot_email_task, rtp.ProcessTaskSendForgotEmail)
	mux.HandleFunc(send_order_email_task, rtp.ProcessTaskOrderEmail)
	mux.HandleFunc(send_update_order_task, rtp.ProcessTaskUpdateOrder)
	mux.HandleFunc(send_update_orders_task, rtp.ProcessTaskUpdateOrders)

	return rtp.server.Start(mux)
}

func (rtp *RedisTaskProcessor) Shutdown() {
	rtp.server.Shutdown()
}
