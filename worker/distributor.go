package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeTaskSendVerifyEmail(ctx context.Context, payload *SendEmailPayload, opts ...asynq.Option) error
	DistributeTaskSendForgotEmail(ctx context.Context, payload *SendEmailPayload, opts ...asynq.Option) error
	DistributeTaskSendBuyOrderEmail(ctx context.Context, payload *SendBuyOrderEmailPayload, opts ...asynq.Option) error
	DistributeTaskRejectBuyOrder(ctx context.Context, payload *RejectBuyOrderPayload, opts ...asynq.Option) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOptions asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOptions)
	return &RedisTaskDistributor{
		client: client,
	}
}
