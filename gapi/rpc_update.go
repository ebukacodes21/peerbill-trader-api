package gapi

import (
	"context"
	"database/sql"
	"time"

	db "github.com/ebukacodes21/peerbill-trader-api/db/sqlc"
	"github.com/ebukacodes21/peerbill-trader-api/pb"
	"github.com/ebukacodes21/peerbill-trader-api/validate"
	"github.com/ebukacodes21/peerbill-trader-api/worker"

	"github.com/hibiken/asynq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateTrader(ctx context.Context, req *pb.UpdateTraderRequest) (*pb.UpdateTraderResponse, error) {
	authPayload, err := s.authorizeTrader(ctx, []string{"user", "admin"})
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateUpdateTraderRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	if authPayload.Role != "admin" && authPayload.TraderID != req.GetTraderId() {
		return nil, status.Errorf(codes.PermissionDenied, "cannot update trader info")
	}

	// Prepare the update parameters
	args := db.UpdateTraderParams{
		ID: req.GetTraderId(),
		FirstName: sql.NullString{
			String: req.GetFirstName(),
			Valid:  req.FirstName != nil,
		},
		LastName: sql.NullString{
			String: req.GetLastName(),
			Valid:  req.LastName != nil,
		},
		Username: sql.NullString{
			String: req.GetUsername(),
			Valid:  req.Username != nil,
		},
		Email: sql.NullString{
			String: req.GetEmail(),
			Valid:  req.Email != nil,
		},
		Country: sql.NullString{
			String: req.GetCountry(),
			Valid:  req.Country != nil,
		},
		Phone: sql.NullString{
			String: req.GetPhone(),
			Valid:  req.Phone != nil,
		},
	}

	// compare old and new email
	trader, err := s.repository.GetTrader(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "trader not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to fetch trader")
	}

	// Check if the email has changed
	if req.GetEmail() != "" && trader.Email != req.GetEmail() {
		args.IsVerified = sql.NullBool{
			Valid: true,
			Bool:  false,
		}

		payload := &worker.SendEmailPayload{
			Username: authPayload.Username,
		}
		opts := []asynq.Option{
			asynq.MaxRetry(10),
			asynq.ProcessIn(10 * time.Second),
			asynq.Queue(worker.Critical),
		}

		err = s.taskDistributor.DistributeTaskSendVerifyEmail(ctx, payload, opts...)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to send verification email")
		}
	}

	// Perform the update operation
	updatedTrader, err := s.repository.UpdateTrader(ctx, args)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update trader %s", err)
	}

	resp := &pb.UpdateTraderResponse{
		Trader: convert(updatedTrader),
	}

	return resp, nil
}

func validateUpdateTraderRequest(req *pb.UpdateTraderRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateId(req.GetTraderId()); err != nil {
		violations = append(violations, fieldViolation("trader_id", err))
	}

	if req.FirstName != nil {
		if err := validate.ValidateFirstname(req.GetFirstName()); err != nil {
			violations = append(violations, fieldViolation("first_name", err))
		}
	}

	if req.LastName != nil {
		if err := validate.ValidateLastname(req.GetLastName()); err != nil {
			violations = append(violations, fieldViolation("last_name", err))
		}
	}

	if req.Username != nil {
		if err := validate.ValidateUsername(req.GetUsername()); err != nil {
			violations = append(violations, fieldViolation("username", err))
		}
	}

	if req.Email != nil {
		if err := validate.ValidateEmail(req.GetEmail()); err != nil {
			violations = append(violations, fieldViolation("email", err))
		}
	}

	if req.Phone != nil {
		if err := validate.ValidatePhone(req.GetPhone()); err != nil {
			violations = append(violations, fieldViolation("phone", err))
		}
	}
	if req.Country != nil {

		if err := validate.ValidateCountry(req.GetCountry()); err != nil {
			violations = append(violations, fieldViolation("country", err))
		}
	}

	return violations
}
