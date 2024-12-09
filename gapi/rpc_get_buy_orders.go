package gapi

import (
	"context"
	"peerbill-trader-api/pb"
	"peerbill-trader-api/validate"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetBuyOrders(ctx context.Context, req *pb.GetBuyOrdersRequest) (*pb.GetBuyOrdersResponse, error) {
	authPayload, err := s.authorizeTrader(ctx, []string{"user", "admin"})
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateGetBuyOrdersRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	if authPayload.Role != "admin" && authPayload.Username != req.GetUsername() {
		return nil, status.Errorf(codes.PermissionDenied, "not authorized to get buy orders")
	}

	// Fetch initial buy orders to send the first response immediately
	buyOrders, err := s.repository.GetOrders(ctx, req.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch buy orders")
	}

	// Reverse the buy orders only once
	reversedOrders := convertOrders(buyOrders)
	for i, j := 0, len(reversedOrders)-1; i < j; i, j = i+1, j-1 {
		reversedOrders[i], reversedOrders[j] = reversedOrders[j], reversedOrders[i]
	}

	resp := &pb.GetBuyOrdersResponse{
		BuyOrders: reversedOrders,
	}

	return resp, nil
}

func validateGetBuyOrdersRequest(req *pb.GetBuyOrdersRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	return violations
}
