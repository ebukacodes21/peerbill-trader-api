package gapi

import (
	"context"
	"database/sql"
	db "peerbill-trader-api/db/sqlc"
	"peerbill-trader-api/pb"
	"peerbill-trader-api/validate"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdatePaymentMethod(ctx context.Context, req *pb.UpdatePaymentMethodRequest) (*pb.UpdatePaymentMethodResponse, error) {
	authPayload, err := s.authorizeTrader(ctx, []string{"user", "admin"})
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateUpdatePaymentMethodRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	if authPayload.Role != "admin" && authPayload.Username != req.GetUsername() {
		return nil, status.Errorf(codes.PermissionDenied, "not authorized to update trader payment method info")
	}

	args := db.UpdatePaymentMethodParams{
		ID:       req.GetId(),
		Username: req.GetUsername(),
		BankName: sql.NullString{
			Valid:  true,
			String: req.GetBankName(),
		},
		AccountHolder: sql.NullString{
			Valid:  true,
			String: req.GetAccountHolder(),
		},
		AccountNumber: sql.NullString{
			Valid:  true,
			String: req.GetAccountNumber(),
		},
		WalletAddress: sql.NullString{
			Valid:  true,
			String: req.GetWalletAddress(),
		},
	}

	err = s.repository.UpdatePaymentMethod(ctx, args)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot update trader payment method info")
	}

	methods, err := s.repository.GetPaymentMethods(ctx, authPayload.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot fetch trader payment method")
	}

	resp := &pb.UpdatePaymentMethodResponse{
		PaymentMethods: convertPaymentMethods(methods),
	}

	return resp, nil
}

func validateUpdatePaymentMethodRequest(req *pb.UpdatePaymentMethodRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateId(req.GetId()); err != nil {
		violations = append(violations, fieldViolation("id", err))
	}

	if err := validate.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if req.Crypto != nil {
		if err := validate.ValidateCrypto(req.GetCrypto()); err != nil {
			violations = append(violations, fieldViolation("crypto", err))
		}
	}

	if req.Fiat != nil {
		if err := validate.ValidateFiat(req.GetFiat()); err != nil {
			violations = append(violations, fieldViolation("fiat", err))
		}
	}

	if req.AccountHolder != nil {
		if err := validate.ValidateFirstname(req.GetAccountHolder()); err != nil {
			violations = append(violations, fieldViolation("account_holder", err))
		}
	}

	if req.BankName != nil {
		if err := validate.ValidateFirstname(req.GetBankName()); err != nil {
			violations = append(violations, fieldViolation("bank_name", err))
		}
	}

	if req.AccountNumber != nil {
		if err := validate.ValidateFirstname(req.GetAccountNumber()); err != nil {
			violations = append(violations, fieldViolation("account_number", err))
		}
	}

	if req.WalletAddress != nil {
		if err := validate.ValidateWalletAddress(req.GetWalletAddress()); err != nil {
			violations = append(violations, fieldViolation("wallet_address", err))
		}
	}
	return violations
}
