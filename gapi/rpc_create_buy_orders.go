package gapi

import (
	"context"
	db "peerbill-trader-api/db/sqlc"
	"peerbill-trader-api/pb"
	"peerbill-trader-api/validate"
	"peerbill-trader-api/worker"
	"time"

	"github.com/hibiken/asynq"
	pg "github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateBuyOrder(ctx context.Context, req *pb.CreateBuyOrderRequest) (*pb.CreateBuyOrderResponse, error) {
	violations := validateCreateBuyOrderRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	args := db.CreateBuyOrderTxParams{
		CreateBuyOrderParams: db.CreateBuyOrderParams{
			WalletAddress: req.GetWallet(),
			FiatAmount:    float64(req.GetPayAmount()),
			CryptoAmount:  float64(req.GetReceiveAmount()),
			Crypto:        req.GetCrypto(),
			Fiat:          req.GetFiat(),
			Rate:          float64(req.GetSellRate()),
			Username:      req.GetUsername(),
			Duration:      time.Now().Add(30 * time.Minute),
		},
		AfterCreate: func(buyOrder db.BuyOrder) error {
			payload := worker.SendBuyOrderEmailPayload{
				Username:     buyOrder.Username,
				Fiat:         buyOrder.Fiat,
				Crypto:       buyOrder.Crypto,
				CryptoAmount: buyOrder.CryptoAmount,
				FiatAmount:   buyOrder.FiatAmount,
			}

			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.Critical),
			}
			return s.taskDistributor.DistributeTaskSendBuyOrderEmail(ctx, &payload, opts...)
		},
	}

	result, err := s.repository.CreateBuyOrderTx(ctx, args)
	if err != nil {
		if pgErr, ok := err.(*pg.Error); ok {
			switch pgErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.Internal, err.Error())
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create buy order %s ", err)
	}

	resp := &pb.CreateBuyOrderResponse{
		BuyOrder: convertBuyOrder(result.BuyOrder),
	}

	return resp, nil
}

func validateCreateBuyOrderRequest(req *pb.CreateBuyOrderRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateWalletAddres(req.GetWallet()); err != nil {
		violations = append(violations, fieldViolation("wallet", err))
	}

	if err := validate.ValidateNumber(req.GetPayAmount()); err != nil {
		violations = append(violations, fieldViolation("pay_amount", err))
	}

	if err := validate.ValidateNumber(req.GetReceiveAmount()); err != nil {
		violations = append(violations, fieldViolation("receive_amount", err))
	}

	if err := validate.ValidateCrypto(req.GetCrypto()); err != nil {
		violations = append(violations, fieldViolation("crypto", err))
	}

	if err := validate.ValidateFiat(req.GetFiat()); err != nil {
		violations = append(violations, fieldViolation("fiat", err))
	}

	if err := validate.ValidateNumber(req.GetSellRate()); err != nil {
		violations = append(violations, fieldViolation("sell_rate", err))
	}

	if err := validate.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	return violations
}
