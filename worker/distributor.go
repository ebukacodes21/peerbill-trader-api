package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeTaskSendVerifyEmail(ctx context.Context, payload *SendEmailPayload, opts ...asynq.Option) error
	DistributeTaskSendForgotEmail(ctx context.Context, payload *SendEmailPayload, opts ...asynq.Option) error
	DistributeTaskSendOrderEmail(ctx context.Context, payload *SendOrderEmailPayload, opts ...asynq.Option) error
	DistributeTaskUpdateOrder(ctx context.Context, payload *UpdateOrderPayload, opts ...asynq.Option) error
	DistributeTaskUpdateOrders(ctx context.Context, payload *UpdateOrderPayload, opts ...asynq.Option) error
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
