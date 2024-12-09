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

func (s *Server) GetBuyOrder(ctx context.Context, req *pb.GetBuyOrderRequest) (*pb.GetBuyOrderResponse, error) {
	violations := validateGetBuyOrderRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	args := db.GetOrderParams{
		ID:        req.GetId(),
		OrderType: req.GetOrderType(),
	}

	buyOrder, err := s.repository.GetOrder(ctx, args)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch buy order")
	}

	resp := &pb.GetBuyOrderResponse{
		BuyOrder: convertOrder(buyOrder),
	}
	return resp, nil
}

func validateGetBuyOrderRequest(req *pb.GetBuyOrderRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateId(req.GetId()); err != nil {
		violations = append(violations, fieldViolation("id", err))
	}

	if err := validate.ValidateType(req.GetOrderType()); err != nil {
		violations = append(violations, fieldViolation("order_type", err))
	}

	return violations
}
