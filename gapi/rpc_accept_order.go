package gapi

import (
	"context"
	"strconv"

	"database/sql"

	db "github.com/ebukacodes21/peerbill-trader-api/db/sqlc"
	"github.com/ebukacodes21/peerbill-trader-api/pb"
	"github.com/ebukacodes21/peerbill-trader-api/utils"
	"github.com/ebukacodes21/peerbill-trader-api/validate"

	"github.com/ebukacodes21/peerbill-trader-api/worker"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) AcceptOrder(ctx context.Context, req *pb.AcceptOrderRequest) (*pb.AcceptOrderResponse, error) {
	authPayload, err := s.authorizeTrader(ctx, []string{"user", "admin"})
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateAcceptOrderRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	if authPayload.Role != "admin" && authPayload.Username != req.GetUsername() {
		return nil, status.Errorf(codes.PermissionDenied, "not authorized to accept orders")
	}

	balance := utils.CheckBalance(ctx, req.GetCrypto(), req.GetEscrowAddress())
	balanceFloat, err := strconv.ParseFloat(balance, 64)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse balance as float")
	}

	// <
	if float32(balanceFloat) > req.GetAmount() {
		return nil, status.Errorf(codes.InvalidArgument, "insufficient balance")
	}

	args := db.UpdateOrderParams{
		ID:       req.GetId(),
		Username: req.GetUsername(),
		IsAccepted: sql.NullBool{
			Valid: true,
			Bool:  true,
		},
	}

	err = s.repository.UpdateOrder(ctx, args)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to accept orders")
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

	// =========================================
	// Reverse the orders array
	orders, err := s.repository.GetOrders(ctx, authPayload.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch orders")
	}

	reversedOrders := convertOrders(orders)
	for i, j := 0, len(reversedOrders)-1; i < j; i, j = i+1, j-1 {
		reversedOrders[i], reversedOrders[j] = reversedOrders[j], reversedOrders[i]
	}

	resp := &pb.AcceptOrderResponse{
		Orders: reversedOrders,
	}

	return resp, nil
}

func validateAcceptOrderRequest(req *pb.AcceptOrderRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateId(req.GetId()); err != nil {
		violations = append(violations, fieldViolation("id", err))
	}

	if err := validate.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := validate.ValidateWalletAddress(req.GetEscrowAddress()); err != nil {
		violations = append(violations, fieldViolation("escrow_address", err))
	}

	if err := validate.ValidateCrypto(req.GetCrypto()); err != nil {
		violations = append(violations, fieldViolation("crypto", err))
	}

	if err := validate.ValidateNumber(req.GetAmount()); err != nil {
		violations = append(violations, fieldViolation("amount", err))
	}

	if err := validate.ValidateType(req.GetOrderType()); err != nil {
		violations = append(violations, fieldViolation("order_type", err))
	}
	return violations
}
