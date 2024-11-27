package gapi

import (
	"context"
	// db "peerbill-trader-api/db/sqlc"
	db "peerbill-trader-api/db/sqlc"
	"peerbill-trader-api/pb"
	"peerbill-trader-api/validate"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetTraders(ctx context.Context, req *pb.GetTradersRequest) (*pb.GetTradersResponse, error) {
	violations := validateGetTradersRequests(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	args := db.GetTradePairsParams{
		Crypto: req.Crypto,
		Fiat:   req.Fiat,
	}
	tradePairs, err := s.repository.GetTradePairs(ctx, args)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "no trade pair found: %s", err)
	}

	var traders []db.Trader
	for _, tradePair := range tradePairs {
		trader, err := s.repository.GetTrader(ctx, tradePair.Username)
		if err != nil {
			return nil, status.Errorf(codes.NotFound, "no trader found: %s", err)
		}

		traders = append(traders, trader)
	}

	resp := &pb.GetTradersResponse{
		Trader: convertTraders(traders),
	}
	return resp, nil
}

func validateGetTradersRequests(req *pb.GetTradersRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateCrypto(req.Crypto); err != nil {
		violations = append(violations, fieldViolation("crypto", err))
	}

	if err := validate.ValidateFiat(req.Fiat); err != nil {
		violations = append(violations, fieldViolation("fiat", err))
	}

	return violations
}
