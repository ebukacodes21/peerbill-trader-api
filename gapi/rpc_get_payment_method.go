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

func (s *Server) GetPaymentMethod(ctx context.Context, req *pb.GetPaymentMethodRequest) (*pb.GetPaymentMethodResponse, error) {
	violations := validateGetPaymentRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	args := db.GetPaymentMethodParams{
		Crypto:   req.GetCrypto(),
		Fiat:     req.GetFiat(),
		Username: req.Username,
	}
	method, err := s.repository.GetPaymentMethod(ctx, args)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "no payment method found")
	}

	resp := &pb.GetPaymentMethodResponse{
		PaymentMethod: convertPaymentMethod(method),
	}

	return resp, nil
}

func validateGetPaymentRequest(req *pb.GetPaymentMethodRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateCrypto(req.GetCrypto()); err != nil {
		violations = append(violations, fieldViolation("crypto", err))
	}

	if err := validate.ValidateFiat(req.GetFiat()); err != nil {
		violations = append(violations, fieldViolation("fiat", err))
	}

	if err := validate.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	return violations
}