package gapi

import (
	db "peerbill-trader-api/db/sqlc"
	"peerbill-trader-api/pb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func convert(trader db.Trader) *pb.Trader {
	return &pb.Trader{
		FirstName: trader.FirstName,
		LastName:  trader.LastName,
		Username:  trader.Username,
		Email:     trader.Email,
		Country:   trader.Country,
		Phone:     trader.Phone,
		CreatedAt: timestamppb.New(trader.CreatedAt),
	}
}

func convertTradePair(tradePair db.TradePair) *pb.TraderPair {
	return &pb.TraderPair{
		Crypto:   tradePair.Crypto,
		Fiat:     tradePair.Fiat,
		Username: tradePair.Username,
		BuyRate:  float32(tradePair.BuyRate),
		SellRate: float32(tradePair.SellRate),
	}
}

func convertTradersWithDetails(tradersWithDetails []TradersWithDetails) []*pb.TraderWithDetails {
	var result []*pb.TraderWithDetails
	for _, item := range tradersWithDetails {
		pbTraderWithDetails := &pb.TraderWithDetails{
			Trader:    convert(item.Trader),
			TradePair: convertTradePair(item.TradePair),
		}
		result = append(result, pbTraderWithDetails)
	}
	return result
}
