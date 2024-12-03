package gapi

import (
	"context"
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

	buyOrder, err := s.repository.GetBuyOrder(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch buy order")
	}

	resp := &pb.GetBuyOrderResponse{
		BuyOrder: convertBuyOrder(buyOrder),
	}
	return resp, nil
}

func validateGetBuyOrderRequest(req *pb.GetBuyOrderRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateTraderId(req.GetId()); err != nil {
		violations = append(violations, fieldViolation("id", err))
	}

	return violations
}
