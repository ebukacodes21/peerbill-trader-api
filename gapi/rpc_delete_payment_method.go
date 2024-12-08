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

func (s *Server) DeletePaymentMethod(ctx context.Context, req *pb.DeletePaymentMethodRequest) (*pb.DeletePaymentMethodResponse, error) {
	authPayload, err := s.authorizeTrader(ctx, []string{"user", "admin"})
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateDeletePaymentMethodRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	if authPayload.Role != "admin" && authPayload.Username != req.GetUsername() {
		return nil, status.Errorf(codes.PermissionDenied, "not authorized to delete trader payment method")
	}

	args := db.DeletePaymentMethodParams{
		ID:       req.GetId(),
		Username: req.GetUsername(),
	}

	err = s.repository.DeletePaymentMethod(ctx, args)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to delete trader payment method")
	}

	methods, err := s.repository.GetPaymentMethods(ctx, authPayload.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot fetch trader pairs")
	}

	resp := &pb.DeletePaymentMethodResponse{
		PaymentMethods: convertPaymentMethods(methods),
	}

	return resp, nil
}

func validateDeletePaymentMethodRequest(req *pb.DeletePaymentMethodRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateTraderId(req.GetId()); err != nil {
		violations = append(violations, fieldViolation("id", err))
	}

	if err := validate.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	return violations
}
