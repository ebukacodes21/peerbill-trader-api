package gapi

import (
	"context"
	// "encoding/json"
	// "fmt"
	// "log"

	// "github.com/go-redis/redis/v8"

	db "peerbill-trader-api/db/sqlc"
	"peerbill-trader-api/pb"
	"peerbill-trader-api/validate"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TradersWithDetails struct {
	Trader    db.Trader
	TradePair db.TradePair
}

func (s *Server) GetTraders(ctx context.Context, req *pb.GetTradersRequest) (*pb.GetTradersResponse, error) {
	violations := validateGetTradersRequests(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	// cacheKey := fmt.Sprintf("traders:%s:%s", req.Crypto, req.Fiat)
	// cachedResult, err := s.taskProcessor.Get(ctx, cacheKey)
	// if err != nil {
	// 	if err != redis.Nil {
	// 		return nil, fmt.Errorf("redis error: %v", err)
	// 	}
	// } else if cachedResult != "" {
	// 	var cachedTraders []TradersWithDetails
	// 	if err := json.Unmarshal([]byte(cachedResult), &cachedTraders); err != nil {
	// 		return nil, fmt.Errorf("failed to unmarshal cached data: %v", err)
	// 	}

	// 	resp := &pb.GetTradersResponse{
	// 		Result: convertTradersWithDetails(cachedTraders),
	// 	}
	// 	log.Print("cached effort")
	// 	return resp, nil

	// }

	args := db.GetTradePairsParams{
		Crypto: req.GetCrypto(),
		Fiat:   req.GetFiat(),
	}
	tradePairs, err := s.repository.GetTradePairs(ctx, args)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "no trade pair found: %s", err)
	}

	tradersMap := make(map[string]db.Trader)
	for _, tradePair := range tradePairs {
		trader, err := s.repository.GetTrader(ctx, tradePair.Username)
		if err != nil {
			return nil, status.Errorf(codes.NotFound, "no trader found for username %s: %s", tradePair.Username, err)
		}

		tradersMap[tradePair.Username] = trader
	}

	var tradersWithDetails []TradersWithDetails
	for _, tradePair := range tradePairs {
		trader, exists := tradersMap[tradePair.Username]
		if !exists {
			return nil, status.Errorf(codes.NotFound, "no trader found for trade pair username: %s", tradePair.Username)
		}

		traderWithDetails := TradersWithDetails{
			Trader:    trader,
			TradePair: tradePair,
		}
		tradersWithDetails = append(tradersWithDetails, traderWithDetails)
	}

	// cachedData, err := json.Marshal(tradersWithDetails)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to marshal traders data: %v", err)
	// }
	// if err := s.taskProcessor.Set(cacheKey, string(cachedData)); err != nil {
	// 	return nil, fmt.Errorf("failed to cache traders data: %v", err)
	// }

	resp := &pb.GetTradersResponse{
		Result: convertTradersWithDetails(tradersWithDetails),
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
