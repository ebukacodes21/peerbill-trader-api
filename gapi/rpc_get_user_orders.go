package gapi

import (
	"context"

	"peerbill-trader-api/pb"
	"peerbill-trader-api/validate"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetUserOrders(ctx context.Context, req *pb.GetUserOrdersRequest) (*pb.GetUserOrdersResponse, error) {
	violations := validateGetUserOrdersRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	// Fetch initial orders to send the first response immediately
	orders, err := s.repository.GetUserOrders(ctx, req.GetUserAddress())

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch user orders")
	}

	// Reverse the orders only once
	reversedOrders := convertOrders(orders)
	for i, j := 0, len(reversedOrders)-1; i < j; i, j = i+1, j-1 {
		reversedOrders[i], reversedOrders[j] = reversedOrders[j], reversedOrders[i]
	}

	resp := &pb.GetUserOrdersResponse{
		Orders: reversedOrders,
	}

	return resp, nil
}

func validateGetUserOrdersRequest(req *pb.GetUserOrdersRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateWalletAddress(req.GetUserAddress()); err != nil {
		violations = append(violations, fieldViolation("user_address", err))
	}

	return violations
}
