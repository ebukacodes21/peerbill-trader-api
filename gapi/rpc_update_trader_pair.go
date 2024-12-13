package gapi

import (
	"context"
	"database/sql"

	db "github.com/ebukacodes21/peerbill-trader-api/db/sqlc"
	"github.com/ebukacodes21/peerbill-trader-api/pb"
	"github.com/ebukacodes21/peerbill-trader-api/validate"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateTraderPair(ctx context.Context, req *pb.UpdateTradePairRequest) (*pb.UpdateTradePairResponse, error) {
	authPayload, err := s.authorizeTrader(ctx, []string{"user", "admin"})
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateUpdateTradePairRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	if authPayload.Role != "admin" && authPayload.Username != req.GetUsername() {
		return nil, status.Errorf(codes.PermissionDenied, "not authorized to update trader info")
	}

	// Prepare the update parameters
	args := db.UpdateTradePairParams{
		ID:       req.GetId(),
		Username: req.GetUsername(),
		Crypto: sql.NullString{
			String: req.GetCrypto(),
			Valid:  req.Crypto != nil,
		},
		Fiat: sql.NullString{
			String: req.GetFiat(),
			Valid:  req.Fiat != nil,
		},
		BuyRate: sql.NullFloat64{
			Float64: float64(req.GetBuyRate()),
			Valid:   req.BuyRate != nil,
		},
		SellRate: sql.NullFloat64{
			Float64: float64(req.GetSellRate()),
			Valid:   req.SellRate != nil,
		},
	}

	err = s.repository.UpdateTradePair(ctx, args)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot update trader info")
	}

	traderPairs, err := s.repository.GetTraderPairs(ctx, authPayload.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot fetch trader pairs")
	}

	resp := &pb.UpdateTradePairResponse{
		TradePairs: convertTradePairs(traderPairs),
	}

	return resp, nil
}

func validateUpdateTradePairRequest(req *pb.UpdateTradePairRequest) (violations []*errdetails.BadRequest_FieldViolation) {
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

	if req.BuyRate != nil {
		if err := validate.ValidateNumber(req.GetBuyRate()); err != nil {
			violations = append(violations, fieldViolation("buy_rate", err))
		}
	}

	if req.SellRate != nil {
		if err := validate.ValidateNumber(req.GetSellRate()); err != nil {
			violations = append(violations, fieldViolation("sell_rate", err))
		}
	}
	return violations
}
