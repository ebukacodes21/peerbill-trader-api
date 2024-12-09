package gapi

import (
	"context"
	db "peerbill-trader-api/db/sqlc"
	"peerbill-trader-api/pb"
	"peerbill-trader-api/validate"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) DeleteTraderPair(ctx context.Context, req *pb.DeleteTradePairRequest) (*pb.DeleteTradePairResponse, error) {
	authPayload, err := s.authorizeTrader(ctx, []string{"user", "admin"})
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateDeleteTradePairRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	if authPayload.Role != "admin" && authPayload.Username != req.GetUsername() {
		return nil, status.Errorf(codes.PermissionDenied, "not authorized to delete trade pair")
	}

	args := db.DeleteTradePairParams{
		ID:       req.GetId(),
		Username: req.GetUsername(),
	}

	err = s.repository.DeleteTradePair(ctx, args)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to delete trade pair")
	}

	traderPairs, err := s.repository.GetTraderPairs(ctx, authPayload.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot fetch trader pairs")
	}

	resp := &pb.DeleteTradePairResponse{
		TradePairs: convertTradePairs(traderPairs),
	}

	return resp, nil
}

func validateDeleteTradePairRequest(req *pb.DeleteTradePairRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateId(req.GetId()); err != nil {
		violations = append(violations, fieldViolation("id", err))
	}

	if err := validate.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	return violations
}
