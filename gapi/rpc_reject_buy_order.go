package gapi

import (
	"context"
	"database/sql"
	db "peerbill-trader-api/db/sqlc"
	"peerbill-trader-api/pb"
	"peerbill-trader-api/validate"
	"peerbill-trader-api/worker"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) RejectBuyOrder(ctx context.Context, req *pb.RejectBuyOrderRequest) (*pb.RejectBuyOrderResponse, error) {
	authPayload, err := s.authorizeTrader(ctx, []string{"user", "admin"})
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateRejectBuyOrderRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	if authPayload.Role != "admin" && authPayload.Username != req.GetUsername() {
		return nil, status.Errorf(codes.PermissionDenied, "not authorized to reject buy orders")
	}

	args := db.UpdateBuyOrderParams{
		ID:       req.GetId(),
		Username: req.GetUsername(),
		IsRejected: sql.NullBool{
			Valid: true,
			Bool:  true,
		},
		IsCompleted: sql.NullBool{
			Valid: true,
			Bool:  true,
		},
	}

	_, err = s.repository.UpdateBuyOrder(ctx, args)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to reject buy orders")
	}

	payload := worker.UpdateBuyOrderPayload{
		ID:       req.GetId(),
		Username: req.GetUsername(),
	}
	err = s.taskDistributor.DistributeTaskUpdateBuyOrder(ctx, &payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to distribute task")
	}

	// Reverse the buy orders array
	buyOrders, err := s.repository.GetBuyOrders(ctx, authPayload.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch buy orders")
	}

	reversedBuyOrders := convertBuyOrders(buyOrders)
	for i, j := 0, len(reversedBuyOrders)-1; i < j; i, j = i+1, j-1 {
		reversedBuyOrders[i], reversedBuyOrders[j] = reversedBuyOrders[j], reversedBuyOrders[i]
	}

	resp := &pb.RejectBuyOrderResponse{
		BuyOrders: reversedBuyOrders,
	}

	return resp, nil
}

func validateRejectBuyOrderRequest(req *pb.RejectBuyOrderRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateTraderId(req.GetId()); err != nil {
		violations = append(violations, fieldViolation("id", err))
	}

	if err := validate.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	return violations
}
