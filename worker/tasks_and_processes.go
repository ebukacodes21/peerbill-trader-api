package worker

import (
	"context"

	db "github.com/ebukacodes21/peerbill-trader-api/db/sqlc"

	// "database/sql"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	send_verify_email_task  = "task:send_verify_email"
	send_forgot_email_task  = "task:send_forgot_email"
	send_order_email_task   = "task:send_order_email"
	send_update_order_task  = "task:send_update_order"
	send_update_orders_task = "task:send_update_orders"
)

type SendEmailPayload struct {
	Username string `json:"username"`
}

type SendOrderEmailPayload struct {
	Username  string `json:"username"`
	OrderType string `db:"order_type" json:"order_type"`
}

type UpdateOrderPayload struct {
	ID        int64  `db:"id" json:"id"`
	Username  string `db:"username" json:"username"`
	OrderType string `db:"order_type" json:"order_type"`
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

func (rtd *RedisTaskDistributor) DistributeTaskSendOrderEmail(ctx context.Context, payload *SendOrderEmailPayload, opts ...asynq.Option) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload %w", err)
	}
	task := asynq.NewTask(send_order_email_task, []byte(data), opts...)

	info, err := rtd.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to queue task")
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("queue", info.Queue).Int("max_retries", info.MaxRetry).Msg("message enqueued")
	return nil
}

func (rtd *RedisTaskDistributor) DistributeTaskUpdateOrder(ctx context.Context, payload *UpdateOrderPayload, opts ...asynq.Option) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload %w", err)
	}
	task := asynq.NewTask(send_update_order_task, []byte(data), opts...)

	info, err := rtd.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to queue task")
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("queue", info.Queue).Int("max_retries", info.MaxRetry).Msg("update order message enqueued")
	return nil
}

func (rtd *RedisTaskDistributor) DistributeTaskUpdateOrders(ctx context.Context, payload *UpdateOrderPayload, opts ...asynq.Option) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload %w", err)
	}
	task := asynq.NewTask(send_update_orders_task, []byte(data), opts...)

	info, err := rtd.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to queue task")
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("queue", info.Queue).Int("max_retries", info.MaxRetry).Msg("update orders message enqueued")
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

func (rtp *RedisTaskProcessor) ProcessTaskOrderEmail(ctx context.Context, task *asynq.Task) error {
	var payload SendOrderEmailPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return asynq.SkipRetry
	}

	trader, err := rtp.repository.GetTrader(ctx, payload.Username)
	if err != nil {
		return fmt.Errorf("failed to get trader with username %s: %w", payload.Username, err)
	}

	orders, err := rtp.repository.GetOrders(ctx, trader.Username)
	if err != nil {
		return fmt.Errorf("failed to fetch updated orders: %w", err)
	}

	for i, j := 0, len(orders)-1; i < j; i, j = i+1, j-1 {
		orders[i], orders[j] = orders[j], orders[i]
	}
	// Broadcast the updated orders to all connected WebSocket clinets
	err = rtp.wsManager.Broadcast(orders, "get-orders")
	if err != nil {
		return fmt.Errorf("failed to broadcast orders via WebSocket: %w", err)
	}

	// Construct the URL to the login page
	url := "http://localhost:3000/auth/signin"

	// Prepare the email content
	subject := "Order Request"
	content := fmt.Sprintf(`Hello %s,<br/>
	You have a pending %s order request. <br/>
	Kindly <a href="%s">click this link to log in to your account and complete the transaction</a>.<br/>
	`,
		trader.Username,
		payload.OrderType,
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

func (rtp *RedisTaskProcessor) ProcessTaskUpdateOrder(ctx context.Context, task *asynq.Task) error {
	var payload UpdateOrderPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return asynq.SkipRetry
	}

	args := db.GetOrderParams{
		ID:        payload.ID,
		OrderType: payload.OrderType,
	}
	order, err := rtp.repository.GetOrder(ctx, args)
	if err != nil {
		return fmt.Errorf("failed to fetch updated order: %w", err)
	}

	// Broadcast the updated orders to all connected WebSocket client
	err = rtp.wsManager.Broadcast(order, "update-order")
	if err != nil {
		return fmt.Errorf("failed to broadcast orders via WebSocket: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Msg("order processed")

	return nil
}

func (rtp *RedisTaskProcessor) ProcessTaskUpdateOrders(ctx context.Context, task *asynq.Task) error {
	var payload UpdateOrderPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return asynq.SkipRetry
	}

	orders, err := rtp.repository.GetOrders(ctx, payload.Username)
	if err != nil {
		return fmt.Errorf("failed to fetch updated orders: %w", err)
	}

	// Broadcast the updated orders to all connected WebSocket client
	err = rtp.wsManager.Broadcast(orders, "get-orders")
	if err != nil {
		return fmt.Errorf("failed to broadcast orders via WebSocket: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Msg("orders processed")

	return nil
}
