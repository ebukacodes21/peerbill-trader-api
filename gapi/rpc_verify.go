package gapi

import (
	"context"
	"database/sql"

	db "peerbill-trader-api/db/sqlc"
	"peerbill-trader-api/pb"
	"peerbill-trader-api/validate"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	violations := validateVerifyEmailRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	verifyArgs := db.FindTraderParams{
		ID:               req.GetUserId(),
		VerificationCode: req.GetVerificationCode(),
	}

	trader, err := s.repository.FindTrader(ctx, verifyArgs)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "trader not found: %s", err)
	}

	if trader.IsVerified {
		return nil, status.Errorf(codes.AlreadyExists, "trader already verified: %s", err)
	}

	updateArgs := db.UpdateTraderParams{
		ID: trader.ID,
		IsVerified: sql.NullBool{
			Valid: true,
			Bool:  true,
		},
	}

	result, err := s.repository.UpdateTrader(ctx, updateArgs)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "verification failed: %s", err)
	}

	resp := &pb.VerifyEmailResponse{
		IsVerified: result.IsVerified,
	}
	return resp, nil
}

func validateVerifyEmailRequest(req *pb.VerifyEmailRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateUserId(req.GetUserId()); err != nil {
		violations = append(violations, fieldViolation("user_id", err))
	}

	if err := validate.ValidateCode(req.GetVerificationCode()); err != nil {
		violations = append(violations, fieldViolation("verification_code", err))
	}

	return violations
}
