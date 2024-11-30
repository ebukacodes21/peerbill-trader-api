package gapi

import (
	"context"
	"peerbill-trader-api/pb"
	"peerbill-trader-api/validate"

	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) LogoutTrader(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	violations := validateLogoutRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	err := s.repository.Logout(ctx, uuid.MustParse(req.GetSessionId()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to delete session")
	}

	message := "logout successful"
	resp := &pb.LogoutResponse{
		Message: message,
	}
	return resp, nil
}

func validateLogoutRequest(req *pb.LogoutRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateToken(req.GetSessionId()); err != nil {
		violations = append(violations, fieldViolation("session_id", err))
	}

	return violations
}
