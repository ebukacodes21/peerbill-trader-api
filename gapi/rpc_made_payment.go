package gapi

import (
	"context"
	"database/sql"
	db "peerbill-trader-api/db/sqlc"
	"peerbill-trader-api/pb"
	"peerbill-trader-api/validate"
	"time"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) MadePayment(ctx context.Context, req *pb.MadePaymentRequest) (*pb.MadePaymentResponse, error) {
	violations := validateMadePaymentRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	args := db.GetOrderParams{
		ID:        req.GetId(),
		OrderType: req.GetOrderType(),
	}

	order, err := s.repository.GetOrder(ctx, args)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "no order found %s", err)
	}

	if order.IsReceived || time.Now().After(order.Duration.Add(30*time.Minute)) {
		// move crypto from escrow to user address

		args := db.UpdateOrderParams{
			ID:       req.GetId(),
			Username: req.GetUsername(),
			IsCompleted: sql.NullBool{
				Valid: true,
				Bool:  true,
			},
		}

		err := s.repository.UpdateOrder(ctx, args)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to update order %s", err)
		}
		resp := &pb.MadePaymentResponse{
			Message: "crypto has been sent to your wallet",
		}
		return resp, nil
	}
	return nil, status.Errorf(codes.FailedPrecondition, "order is not completed and 30 minutes have not passed since the order's duration")
}

func validateMadePaymentRequest(req *pb.MadePaymentRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateId(req.GetId()); err != nil {
		violations = append(violations, fieldViolation("id", err))
	}

	if err := validate.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := validate.ValidateType(req.GetOrderType()); err != nil {
		violations = append(violations, fieldViolation("order_type", err))
	}

	if err := validate.ValidateWalletAddress(req.GetUserAddress()); err != nil {
		violations = append(violations, fieldViolation("user_address", err))
	}

	return violations
}
