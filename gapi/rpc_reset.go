package gapi

import (
	"context"
	"database/sql"
	"log"

	db "github.com/ebukacodes21/peerbill-trader-api/db/sqlc"
	"github.com/ebukacodes21/peerbill-trader-api/pb"
	"github.com/ebukacodes21/peerbill-trader-api/utils"
	"github.com/ebukacodes21/peerbill-trader-api/validate"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Reset(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	violations := validateResetPasswordRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	log.Print(req.Token)
	payload, err := s.token.VerifyToken(req.GetToken())
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "expired token: %s", err)
	}

	trader, err := s.repository.GetTrader(ctx, payload.Username)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "trader not found: %s", err)
	}

	hash, err := utils.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	updateArgs := db.UpdateTraderParams{
		ID: trader.ID,
		Password: sql.NullString{
			Valid:  true,
			String: hash,
		},
	}

	_, err = s.repository.UpdateTrader(ctx, updateArgs)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Password reset failed: %s", err)
	}

	message := "Password updated successfully"
	resp := &pb.ResetPasswordResponse{
		Message: message,
	}
	return resp, nil
}

func validateResetPasswordRequest(req *pb.ResetPasswordRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateToken(req.GetToken()); err != nil {
		violations = append(violations, fieldViolation("token", err))
	}

	if err := validate.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return violations
}
