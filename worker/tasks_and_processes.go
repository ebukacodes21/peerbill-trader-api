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
	send_verify_email_task     = "task:send_verify_email"
	send_forgot_email_task     = "task:send_forgot_email"
	send_buy_order_email_task  = "task:send_buy_order_email"
	send_reject_buy_order_task = "task:send_reject_buy_order"
)

type SendEmailPayload struct {
	Username string `json:"username"`
}

type SendBuyOrderEmailPayload struct {
	Username     string  `json:"username"`
	Crypto       string  `db:"crypto" json:"crypto"`
	Fiat         string  `db:"fiat" json:"fiat"`
	CryptoAmount float64 `db:"crypto_amount" json:"crypto_amount"`
	FiatAmount   float64 `db:"fiat_amount" json:"fiat_amount"`
}

type RejectBuyOrderPayload struct {
	ID       int64  `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
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

func (rtd *RedisTaskDistributor) DistributeTaskSendBuyOrderEmail(ctx context.Context, payload *SendBuyOrderEmailPayload, opts ...asynq.Option) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload %w", err)
	}
	task := asynq.NewTask(send_buy_order_email_task, []byte(data), opts...)

	info, err := rtd.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to queue task")
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("queue", info.Queue).Int("max_retries", info.MaxRetry).Msg("message enqueued")
	return nil
}

func (rtd *RedisTaskDistributor) DistributeTaskRejectBuyOrder(ctx context.Context, payload *RejectBuyOrderPayload, opts ...asynq.Option) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload %w", err)
	}
	task := asynq.NewTask(send_reject_buy_order_task, []byte(data), opts...)

	info, err := rtd.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to queue task")
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("queue", info.Queue).Int("max_retries", info.MaxRetry).Msg("reject buy order message enqueued")
	return nil
}

// processor ==========================================================
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

	url := fmt.Sprintf("http://localhost:3000/auth/verify?trader_id=%d&verification_code=%s", trader.ID, trader.VerificationCode)

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

	accessToken, _, err := rtp.token.CreateToken(trader.Username, trader.ID, trader.Role, rtp.config.TokenAccess)
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

func (rtp *RedisTaskProcessor) ProcessTaskBuyOrderEmail(ctx context.Context, task *asynq.Task) error {
	var payload SendBuyOrderEmailPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return asynq.SkipRetry
	}

	trader, err := rtp.repository.GetTrader(ctx, payload.Username)
	if err != nil {
		return fmt.Errorf("failed to get trader with username %s: %w", payload.Username, err)
	}

	buyOrders, err := rtp.repository.GetBuyOrders(ctx, trader.Username)
	if err != nil {
		return fmt.Errorf("failed to fetch updated buy orders: %w", err)
	}

	for i, j := 0, len(buyOrders)-1; i < j; i, j = i+1, j-1 {
		buyOrders[i], buyOrders[j] = buyOrders[j], buyOrders[i]
	}
	// Broadcast the updated buy orders to all connected WebSocket clinets
	err = rtp.wsManager.Broadcast(buyOrders, "get-buy-orders")
	if err != nil {
		return fmt.Errorf("failed to broadcast buy orders via WebSocket: %w", err)
	}

	// Construct the URL to the login page
	url := "http://localhost:3000/auth/signin"

	// Prepare the email content
	subject := "Buy Request"
	content := fmt.Sprintf(`Hello %s,<br/>
	You have a pending buy request for %.8f %s. <br/>
	You will receive an equivalent of %.2f %s after a successful transaction.<br/>
	Kindly <a href="%s">click this link to log in to your account and complete the transaction</a>.<br/>
	`,
		trader.Username,
		payload.CryptoAmount,
		payload.Crypto,
		payload.FiatAmount,
		payload.Fiat,
		url)

	// Send the email
	to := []string{trader.Email}
	if err := rtp.mailer.SendMail(subject, content, to, nil, nil, nil); err != nil {
		// Log and return a more specific error message
		return fmt.Errorf("failed to send buy order email to %s: %w", trader.Email, err)
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("email", trader.Email).
		Msg("Buy order processed")

	return nil
}

func (rtp *RedisTaskProcessor) ProcessTaskRejectBuyOrder(ctx context.Context, task *asynq.Task) error {
	var payload RejectBuyOrderPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return asynq.SkipRetry
	}

	buyOrder, err := rtp.repository.GetBuyOrder(ctx, payload.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch updated buy order: %w", err)
	}

	// Broadcast the updated buy orders to all connected WebSocket client
	err = rtp.wsManager.Broadcast(buyOrder, "reject-buy-order")
	if err != nil {
		return fmt.Errorf("failed to broadcast buy orders via WebSocket: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Msg("order processed")

	return nil
}
