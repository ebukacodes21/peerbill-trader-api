package gapi

import (
	"context"
	"database/sql"
	"time"

	db "github.com/ebukacodes21/peerbill-trader-api/db/sqlc"
	"github.com/ebukacodes21/peerbill-trader-api/pb"
	"github.com/ebukacodes21/peerbill-trader-api/validate"
	"github.com/ebukacodes21/peerbill-trader-api/worker"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.UpdateOrderResponse, error) {
	violations := validateUpdateOrderRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	updateParams := db.UpdateOrderParams{
		ID:       req.GetId(),
		Username: req.GetUsername(),
	}

	if req.GetOrderType() == "sell" && req.AccountHolder != nil {
		updateParams.BankName = sql.NullString{
			Valid:  true,
			String: req.GetBankName(),
		}
		updateParams.AccountNumber = sql.NullString{
			Valid:  true,
			String: req.GetAccountNumber(),
		}
		updateParams.AccountHolder = sql.NullString{
			Valid:  true,
			String: req.GetAccountHolder(),
		}
	} else {
		updateParams.Duration = sql.NullTime{
			Valid: true,
			Time:  time.Now(),
		}
		updateParams.IsCompleted = sql.NullBool{
			Valid: true,
			Bool:  true,
		}
		updateParams.IsExpired = sql.NullBool{
			Valid: true,
			Bool:  req.GetIsExpired(),
		}
	}

	err := s.repository.UpdateOrder(ctx, updateParams)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update order")
	}

	payload := worker.UpdateOrderPayload{
		ID:        req.GetId(),
		Username:  req.GetUsername(),
		OrderType: req.GetOrderType(),
	}
	err = s.taskDistributor.DistributeTaskUpdateOrders(ctx, &payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to distribute task")
	}

	resp := &pb.UpdateOrderResponse{
		Message: "order completed",
	}
	return resp, nil
}

func validateUpdateOrderRequest(req *pb.UpdateOrderRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateId(req.GetId()); err != nil {
		violations = append(violations, fieldViolation("id", err))
	}

	if err := validate.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := validate.ValidateBool(req.GetIsExpired()); err != nil {
		violations = append(violations, fieldViolation("is_expired", err))
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
