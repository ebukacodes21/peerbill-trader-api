package gapi

import (
	"context"
	"database/sql"
	"time"

	db "github.com/ebukacodes21/peerbill-trader-api/db/sqlc"
	"github.com/ebukacodes21/peerbill-trader-api/pb"
	"github.com/ebukacodes21/peerbill-trader-api/validate"
	"github.com/ebukacodes21/peerbill-trader-api/worker"
	pg "github.com/lib/pq"

	"github.com/hibiken/asynq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	violations := validateCreateOrderRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	args := db.CreateOrderTxParams{
		CreateOrderParams: db.CreateOrderParams{
			EscrowAddress: sql.NullString{
				Valid:  true,
				String: req.GetEscrowAddress(),
			},
			UserAddress:  req.GetUserAddress(),
			OrderType:    req.GetOrderType(),
			FiatAmount:   float64(req.GetFiatAmount()),
			CryptoAmount: float64(req.GetCryptoAmount()),
			Crypto:       req.GetCrypto(),
			Fiat:         req.GetFiat(),
			Rate:         float64(req.GetRate()),
			Username:     req.GetUsername(),
			Duration:     time.Now().Add(30 * time.Minute),
			BankName: sql.NullString{
				Valid:  true,
				String: req.GetBankName(),
			},
			AccountNumber: sql.NullString{
				Valid:  true,
				String: req.GetAccountNumber(),
			},
			AccountHolder: sql.NullString{
				Valid:  true,
				String: req.GetAccountHolder(),
			},
		},
		AfterCreate: func(buyOrder db.Order) error {
			payload := worker.SendOrderEmailPayload{
				Username:  buyOrder.Username,
				OrderType: buyOrder.OrderType,
			}

			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.Critical),
			}
			return s.taskDistributor.DistributeTaskSendOrderEmail(ctx, &payload, opts...)
		},
	}

	result, err := s.repository.CreateOrderTx(ctx, args)
	if err != nil {
		if pgErr, ok := err.(*pg.Error); ok {
			switch pgErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.Internal, "unqiue violation %s ", pgErr)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create order %s ", err)
	}

	resp := &pb.CreateOrderResponse{
		Order: convertOrder(result.Order),
	}

	return resp, nil
}

func validateCreateOrderRequest(req *pb.CreateOrderRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if req.EscrowAddress != nil {
		if err := validate.ValidateWalletAddress(req.GetEscrowAddress()); err != nil {
			violations = append(violations, fieldViolation("escrow_address", err))
		}
	}

	if err := validate.ValidateWalletAddress(req.GetUserAddress()); err != nil {
		violations = append(violations, fieldViolation("user_address", err))
	}

	if err := validate.ValidateNumber(req.GetCryptoAmount()); err != nil {
		violations = append(violations, fieldViolation("crypto_amount", err))
	}

	if err := validate.ValidateNumber(req.GetFiatAmount()); err != nil {
		violations = append(violations, fieldViolation("fiat_amount", err))
	}

	if err := validate.ValidateType(req.GetOrderType()); err != nil {
		violations = append(violations, fieldViolation("order_type", err))
	}

	if err := validate.ValidateFiat(req.GetFiat()); err != nil {
		violations = append(violations, fieldViolation("fiat", err))
	}

	if err := validate.ValidateFiat(req.GetCrypto()); err != nil {
		violations = append(violations, fieldViolation("crypto", err))
	}

	if err := validate.ValidateNumber(req.GetRate()); err != nil {
		violations = append(violations, fieldViolation("rate", err))
	}

	if err := validate.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if req.AccountHolder != nil {
		if err := validate.ValidateFirstname(req.GetAccountHolder()); err != nil {
			violations = append(violations, fieldViolation("account_holder", err))
		}
	}

	if req.BankName != nil {
		if err := validate.ValidateFirstname(req.GetBankName()); err != nil {
			violations = append(violations, fieldViolation("bank_name", err))
		}
	}

	if req.AccountNumber != nil {
		if err := validate.ValidateFirstname(req.GetAccountNumber()); err != nil {
			violations = append(violations, fieldViolation("account_number", err))
		}
	}

	return violations
}
