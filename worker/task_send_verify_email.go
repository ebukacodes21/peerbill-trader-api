package worker

import (
	"context"

	// "database/sql"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	send_verify_email_task = "task:send_verify_email"
	send_forgot_email_task = "task:send_forgot_email"
)

type SendEmailPayload struct {
	Username string `json:"username"`
}

func (rtd *RedisTaskDistributor) DistributeTaskSendVerifyEmail(ctx context.Context, payload *SendEmailPayload, opts ...asynq.Option) error {
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

func (rtd *RedisTaskDistributor) DistributeTaskSendForgotEmail(ctx context.Context, payload *SendEmailPayload, opts ...asynq.Option) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload %w", err)
	}
	task := asynq.NewTask(send_forgot_email_task, []byte(data), opts...)

	info, err := rtd.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to queue task")
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("queue", info.Queue).Int("max_retries", info.MaxRetry).Msg("message enqueued")
	return nil
}

func (rtp *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload SendEmailPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to marshal payload %w", asynq.SkipRetry)
	}

	trader, err := rtp.repository.GetTrader(ctx, payload.Username)
	if err != nil {
		// if err == sql.ErrNoRows {
		// 	return fmt.Errorf("trader not found %w", asynq.SkipRetry)
		// }

		return fmt.Errorf("failed to get trader")
	}

	url := fmt.Sprintf("http://localhost:3000/auth/verify?user_id=%d&verification_code=%s", trader.ID, trader.VerificationCode)

	subject := "Welcome to Peerbill"
	content := fmt.Sprintf(`Hello %s,<br/>
	You have registered on Peerbill as a Trader. <br/>
	Kindly <a href="%s">click this link to verify your email address</a>
	`, trader.Username, url)
	to := []string{trader.Email}

	if err := rtp.mailer.SendMail(subject, content, to, nil, nil, nil); err != nil {
		return fmt.Errorf("failed to send verify email")
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("email", trader.Email).Msg("message processed")
	return nil
}

func (rtp *RedisTaskProcessor) ProcessTaskSendForgotEmail(ctx context.Context, task *asynq.Task) error {
	var payload SendEmailPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to marshal payload %w", asynq.SkipRetry)
	}

	trader, err := rtp.repository.GetTrader(ctx, payload.Username)
	if err != nil {
		return fmt.Errorf("failed to get trader")
	}

	accessToken, _, err := rtp.token.CreateToken(trader.Username, trader.Role, rtp.config.TokenAccess)
	if err != nil {
		return fmt.Errorf("failed to create token")
	}

	url := fmt.Sprintf("http://localhost:3000/auth/reset?token=%s", accessToken)

	subject := "Password Reset"
	content := fmt.Sprintf(`Hello %s,<br/>
	You have requested to change your password. <br/>
	Kindly <a href="%s">click this link to reset your password</a>
	`, trader.Username, url)
	to := []string{trader.Email}

	if err := rtp.mailer.SendMail(subject, content, to, nil, nil, nil); err != nil {
		return fmt.Errorf("failed to send reset email")
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("email", trader.Email).Msg("message processed")
	return nil
}
