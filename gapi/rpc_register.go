package gapi

import (
	"context"
	db "peerbill-server/db/sqlc"
	"peerbill-server/pb"
	"peerbill-server/utils"
	"peerbill-server/validate"
	"peerbill-server/worker"
	"time"

	"github.com/hibiken/asynq"
	pg "github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) RegisterTrader(ctx context.Context, req *pb.RegisterTraderRequest) (*pb.RegisterTraderResponse, error) {
	violations := validateRegisterTraderRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	hash, err := utils.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password")
	}

	args := db.CreateTraderTxParams{
		CreateTraderParams: db.CreateTraderParams{
			FirstName: req.GetFirstName(),
			LastName:  req.GetLastName(),
			Username:  req.GetUsername(),
			Password:  hash,
			Email:     req.GetEmail(),
			Country:   req.GetCountry(),
			Phone:     req.GetPhone(),
		},
		AfterCreate: func(trader db.Trader) error {
			payload := worker.SendVerifyEmailPayload{
				Username: trader.Username,
			}

			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.Critical),
			}
			return s.taskDistributor.DistributeTaskSendVerifyEmail(ctx, &payload, opts...)
		},
	}

	result, err := s.repository.CreateTraderTx(ctx, args)
	if err != nil {
		if pgErr, ok := err.(*pg.Error); ok {
			switch pgErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.Internal, "already exists")
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user")
	}

	resp := &pb.RegisterTraderResponse{
		Trader: convert(result.Trader),
	}

	return resp, nil
}

func validateRegisterTraderRequest(req *pb.RegisterTraderRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateFirstname(req.GetFirstName()); err != nil {
		violations = append(violations, fieldViolation("first_name", err))
	}

	if err := validate.ValidateLastname(req.GetLastName()); err != nil {
		violations = append(violations, fieldViolation("last_name", err))
	}

	if err := validate.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := validate.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	if err := validate.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	if err := validate.ValidatePhone(req.GetPhone()); err != nil {
		violations = append(violations, fieldViolation("phone", err))
	}

	if err := validate.ValidateCountry(req.GetCountry()); err != nil {
		violations = append(violations, fieldViolation("country", err))
	}

	return violations
}
