package gapi

import (
	"context"

	db "github.com/ebukacodes21/peerbill-trader-api/db/sqlc"
	"github.com/ebukacodes21/peerbill-trader-api/pb"
	"github.com/ebukacodes21/peerbill-trader-api/validate"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
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
		return nil, status.Errorf(codes.NotFound, "order not found %s", err)
	}

	resp := &pb.GetOrderResponse{
		Order: convertOrder(buyOrder),
	}
	return resp, nil
}

func validateGetBuyOrderRequest(req *pb.GetOrderRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateId(req.GetId()); err != nil {
		violations = append(violations, fieldViolation("id", err))
	}

	if err := validate.ValidateType(req.GetOrderType()); err != nil {
		violations = append(violations, fieldViolation("order_type", err))
	}

	return violations
}
