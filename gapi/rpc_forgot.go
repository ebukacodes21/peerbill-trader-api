package gapi

import (
	"context"
	"peerbill-trader-api/pb"
	"peerbill-trader-api/validate"
	"peerbill-trader-api/worker"
	"time"

	"github.com/hibiken/asynq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Forgot(ctx context.Context, req *pb.ForgotPasswordRequest) (*pb.ForgotPasswordResponse, error) {
	violations := validateForgotlRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	trader, err := s.repository.GetTraderEmail(ctx, req.GetEmail())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "no trader found with such email: %s", err)
	}

	payload := worker.SendEmailPayload{
		Username: trader.Username,
	}

	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.Critical),
	}
	err = s.taskDistributor.DistributeTaskSendForgotEmail(ctx, &payload, opts...)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to distribute task: %s", err)
	}

	message := "a link has been sent to the email you provided"
	resp := &pb.ForgotPasswordResponse{
		Message: message,
	}

	return resp, nil
}

func validateForgotlRequest(req *pb.ForgotPasswordRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("user_id", err))
	}

	return violations
}
