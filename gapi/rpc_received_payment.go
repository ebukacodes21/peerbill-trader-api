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

func (s *Server) ReceivePayment(ctx context.Context, req *pb.ReceivedPaymentRequest) (*pb.ReceivedPaymentResponse, error) {
	authPayload, err := s.authorizeTrader(ctx, []string{"user", "admin"})
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateReceivedPaymentRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	if authPayload.Role != "admin" && authPayload.Username != req.GetUsername() {
		return nil, status.Errorf(codes.PermissionDenied, "not authorized to access route")
	}

	// update timer if trader has not received fiat or mark order as completed
	// distribute to connected websockets
	if err := s.updateOrderDurationOrReceivedState(ctx, req); err != nil {
		return nil, err
	}

	// Fetch updated orders and reverse them
	orders, err := s.fetchAndReverseOrders(ctx, authPayload.Username)
	if err != nil {
		return nil, err
	}

	resp := &pb.ReceivedPaymentResponse{
		Orders: orders,
	}
	return resp, nil
}

func validateReceivedPaymentRequest(req *pb.ReceivedPaymentRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateId(req.GetId()); err != nil {
		violations = append(violations, fieldViolation("id", err))
	}

	if err := validate.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := validate.ValidateType(req.GetOrderType()); err != nil {
		violations = append(violations, fieldViolation("order_type", err))
	}

	if err := validate.ValidateBool(req.GetReceivedPayment()); err != nil {
		violations = append(violations, fieldViolation("received_payment", err))
	}

	return violations
}

func (s *Server) updateOrderDurationOrReceivedState(ctx context.Context, req *pb.ReceivedPaymentRequest) error {
	if !req.GetReceivedPayment() {
		args := db.UpdateOrderParams{
			ID:       req.GetId(),
			Username: req.GetUsername(),
			Duration: sql.NullTime{
				Valid: true,
				Time:  time.Now().Add(30 * time.Minute),
			},
		}
		if err := s.repository.UpdateOrder(ctx, args); err != nil {
			return status.Errorf(codes.Internal, "unable to update order %s", err)
		}
	} else {
		args := db.UpdateOrderParams{
			ID:       req.GetId(),
			Username: req.GetUsername(),
			IsReceived: sql.NullBool{
				Valid: true,
				Bool:  true,
			},
		}
		if err := s.repository.UpdateOrder(ctx, args); err != nil {
			return status.Errorf(codes.Internal, "unable to update order %s", err)
		}
	}

	payload := worker.UpdateOrderPayload{
		ID:        req.GetId(),
		Username:  req.GetUsername(),
		OrderType: req.GetOrderType(),
	}
	if err := s.taskDistributor.DistributeTaskUpdateOrder(ctx, &payload); err != nil {
		return status.Errorf(codes.Internal, "unable to distribute upload task %s", err)
	}

	return nil
}

func (s *Server) fetchAndReverseOrders(ctx context.Context, username string) ([]*pb.Order, error) {
	orders, err := s.repository.GetOrders(ctx, username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to get orders %s", err)
	}

	// Reverse the orders
	reversedOrders := convertOrders(orders)
	for i, j := 0, len(reversedOrders)-1; i < j; i, j = i+1, j-1 {
		reversedOrders[i], reversedOrders[j] = reversedOrders[j], reversedOrders[i]
	}

	return reversedOrders, nil
}
