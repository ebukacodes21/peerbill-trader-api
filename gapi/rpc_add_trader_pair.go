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

func (s *Server) AddTraderPair(ctx context.Context, req *pb.AddTradePairRequest) (*pb.AddTradePairResponse, error) {
	authPayload, err := s.authorizeTrader(ctx, []string{"user", "admin"})
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateAddTradePairRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	if authPayload.Role != "admin" && authPayload.Username != req.GetUsername() {
		return nil, status.Errorf(codes.PermissionDenied, "cannot update trader info")
	}

	args := &db.CreateTradePairParams{
		Username: req.Username,
		BuyRate:  float64(req.BuyRate),
		SellRate: float64(req.SellRate),
		Crypto:   req.Crypto,
		Fiat:     req.Fiat,
	}

	_, err = s.repository.CreateTradePair(ctx, *args)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create trade pair%s", err)
	}

	tradePairs, err := s.repository.GetTraderPairs(ctx, authPayload.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch trade pairs")
	}

	resp := &pb.AddTradePairResponse{
		TradePairs: convertTradePairs(tradePairs),
	}

	return resp, nil
}

func validateAddTradePairRequest(req *pb.AddTradePairRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateCrypto(req.GetCrypto()); err != nil {
		violations = append(violations, fieldViolation("crypto", err))
	}

	if err := validate.ValidateFiat(req.GetFiat()); err != nil {
		violations = append(violations, fieldViolation("fiat", err))
	}

	if err := validate.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := validate.ValidateNumber(req.GetBuyRate()); err != nil {
		violations = append(violations, fieldViolation("buy_rate", err))
	}

	if err := validate.ValidateNumber(req.GetSellRate()); err != nil {
		violations = append(violations, fieldViolation("sell_rate", err))
	}

	return violations
}
