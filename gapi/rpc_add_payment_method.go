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

func (s *Server) AddPaymentMethod(ctx context.Context, req *pb.AddPaymentMethodRequest) (*pb.AddPaymentMethodResponse, error) {
	authPayload, err := s.authorizeTrader(ctx, []string{"user", "admin"})
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateAddPaymentMethodRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	if authPayload.Role != "admin" && authPayload.Username != req.GetUsername() {
		return nil, status.Errorf(codes.PermissionDenied, "cannot update trader info")
	}

	tpArgs := &db.GetTradePairParams{
		Crypto:   req.GetCrypto(),
		Fiat:     req.GetFiat(),
		Username: req.GetUsername(),
	}

	tradePair, err := s.repository.GetTradePair(ctx, *tpArgs)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "trade pair does not exist")
	}

	args := &db.CreatePaymentMethodParams{
		Crypto:        req.GetCrypto(),
		Fiat:          req.GetFiat(),
		Username:      req.GetUsername(),
		AccountHolder: req.GetAccountHolder(),
		AccountNumber: req.GetAccountNumber(),
		WalletAddress: req.GetWalletAddress(),
		BankName:      req.GetBankName(),
		TradePairID:   tradePair.ID,
	}

	_, err = s.repository.CreatePaymentMethod(ctx, *args)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create payment method %s", err)
	}

	methods, err := s.repository.GetPaymentMethods(ctx, authPayload.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch payment methods")
	}

	resp := &pb.AddPaymentMethodResponse{
		PaymentMethods: convertPaymentMethods(methods),
	}

	return resp, nil
}

func validateAddPaymentMethodRequest(req *pb.AddPaymentMethodRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateCrypto(req.GetCrypto()); err != nil {
		violations = append(violations, fieldViolation("crypto", err))
	}

	if err := validate.ValidateFiat(req.GetFiat()); err != nil {
		violations = append(violations, fieldViolation("fiat", err))
	}

	if err := validate.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := validate.ValidateFirstname(req.GetAccountHolder()); err != nil {
		violations = append(violations, fieldViolation("account_holder", err))
	}

	if err := validate.ValidateFirstname(req.GetBankName()); err != nil {
		violations = append(violations, fieldViolation("bank_name", err))
	}

	if err := validate.ValidateFirstname(req.GetAccountNumber()); err != nil {
		violations = append(violations, fieldViolation("account_number", err))
	}

	if err := validate.ValidateWalletAddress(req.GetWalletAddress()); err != nil {
		violations = append(violations, fieldViolation("wallet_address", err))
	}

	return violations
}
