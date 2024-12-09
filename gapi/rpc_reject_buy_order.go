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
		return nil, status.Errorf(codes.PermissionDenied, "not authorized to reject buy order")
	}

	args := db.UpdateOrderParams{
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

	err = s.repository.UpdateOrder(ctx, args)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to reject buy order")
	}

	payload := worker.UpdateOrderPayload{
		ID:        req.GetId(),
		Username:  req.GetUsername(),
		OrderType: req.GetOrderType(),
	}
	err = s.taskDistributor.DistributeTaskUpdateOrder(ctx, &payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to distribute task")
	}

	// Reverse the buy orders array
	orders, err := s.repository.GetOrders(ctx, authPayload.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch buy orders")
	}

	reversedOrders := convertOrders(orders)
	for i, j := 0, len(reversedOrders)-1; i < j; i, j = i+1, j-1 {
		reversedOrders[i], reversedOrders[j] = reversedOrders[j], reversedOrders[i]
	}

	resp := &pb.RejectBuyOrderResponse{
		BuyOrders: reversedOrders,
	}

	return resp, nil
}

func validateRejectBuyOrderRequest(req *pb.RejectBuyOrderRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateId(req.GetId()); err != nil {
		violations = append(violations, fieldViolation("id", err))
	}

	if err := validate.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := validate.ValidateType(req.GetOrderType()); err != nil {
		violations = append(violations, fieldViolation("order_type", err))
	}
	return violations
}
