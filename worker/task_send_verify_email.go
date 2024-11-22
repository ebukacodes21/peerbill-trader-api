package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	send_verify_email_task = "task:send_verify_email"
)

type SendVerifyEmailPayload struct {
	Username string `json:"username"`
}

func (rtd *RedisTaskDistributor) DistributeTaskSendVerifyEmail(ctx context.Context, payload *SendVerifyEmailPayload, opts ...asynq.Option) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload %w", err)
	}
	task := asynq.NewTask(send_verify_email_task, []byte(data), opts...)

	info, err := rtd.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to queue task")
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("queue", info.Queue).Int("max_retries", info.MaxRetry).Msg("message enqueued")
	return nil
}

func (rtp *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload SendVerifyEmailPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to marshal payload %w", asynq.SkipRetry)
	}

	trader, err := rtp.repository.GetTrader(ctx, payload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("trader not found %w", asynq.SkipRetry)
		}

		return fmt.Errorf("failed to get trader")
	}
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("email", trader.Email).Msg("message enqueued")
	return nil
}
