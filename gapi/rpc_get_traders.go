package gapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"

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
	// Validate the incoming request
	violations := validateGetTradersRequests(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	// Create cache key based on request parameters
	cacheKey := fmt.Sprintf("traders:%s:%s", req.Crypto, req.Fiat)

	// Try to get the response from the Redis cache
	cachedResult, err := s.taskProcessor.Get(ctx, cacheKey)
	if err != nil {
		if err != redis.Nil {
			// If there’s a Redis error other than a cache miss, return the error
			return nil, fmt.Errorf("redis error: %v", err)
		}
		// If there’s a cache miss (redis.Nil), continue with database call
	} else if cachedResult != "" {
		// Cache hit, unmarshal the cached response into the correct struct
		var cachedTraders []TradersWithDetails
		if err := json.Unmarshal([]byte(cachedResult), &cachedTraders); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cached data: %v", err)
		}
		log.Print("result from cache")

		// Return cached data
		resp := &pb.GetTradersResponse{
			Result: convertTradersWithDetails(cachedTraders),
		}
		return resp, nil
	}

	// Cache miss, need to fetch from database
	args := db.GetTradePairsParams{
		Crypto: req.Crypto,
		Fiat:   req.Fiat,
	}
	tradePairs, err := s.repository.GetTradePairs(ctx, args)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "no trade pair found: %s", err)
	}

	// Map to store traders by username
	tradersMap := make(map[string]db.Trader)

	// Fetch each trader and map by username for quick lookup
	for _, tradePair := range tradePairs {
		trader, err := s.repository.GetTrader(ctx, tradePair.Username)
		if err != nil {
			return nil, status.Errorf(codes.NotFound, "no trader found for username %s: %s", tradePair.Username, err)
		}

		// Store trader by username to avoid fetching multiple times
		tradersMap[tradePair.Username] = trader
	}

	var tradersWithDetails []TradersWithDetails

	// Iterate over the trade pairs and map each to the correct trader
	for _, tradePair := range tradePairs {
		trader, exists := tradersMap[tradePair.Username]
		if !exists {
			return nil, status.Errorf(codes.NotFound, "no trader found for trade pair username: %s", tradePair.Username)
		}

		// Create a TradersWithDetails struct and add it to the slice
		traderWithDetails := TradersWithDetails{
			Trader:    trader,
			TradePair: tradePair,
		}

		// Append the TradersWithDetails struct to the slice
		tradersWithDetails = append(tradersWithDetails, traderWithDetails)
	}

	// Prepare the response
	resp := &pb.GetTradersResponse{
		Result: convertTradersWithDetails(tradersWithDetails),
	}

	// Cache the result in Redis for future use (with an expiration time, e.g., 1 hour)
	cachedData, err := json.Marshal(tradersWithDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal traders data: %v", err)
	}

	// Set the cache with a TTL (time-to-live) of 1 hour
	if err := s.taskProcessor.Set(cacheKey, string(cachedData)); err != nil {
		return nil, fmt.Errorf("failed to cache traders data: %v", err)
	}

	// Return the response
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
