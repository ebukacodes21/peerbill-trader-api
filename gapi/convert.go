package gapi

import (
	"log"
	db "peerbill-trader-api/db/sqlc"
	"peerbill-trader-api/pb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func convert(trader db.Trader) *pb.Trader {
	return &pb.Trader{
		TraderId:  trader.ID,
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
		Id:       tradePair.ID,
		Crypto:   tradePair.Crypto,
		Fiat:     tradePair.Fiat,
		Username: tradePair.Username,
		BuyRate:  float32(tradePair.BuyRate),
		SellRate: float32(tradePair.SellRate),
	}
}

func convertTradePairs(tradePairs []db.TradePair) []*pb.TraderPair {
	var pbTradePairs []*pb.TraderPair
	for _, tradePair := range tradePairs {
		pbTradePairs = append(pbTradePairs, &pb.TraderPair{
			Id:       tradePair.ID,
			Crypto:   tradePair.Crypto,
			Fiat:     tradePair.Fiat,
			Username: tradePair.Username,
			BuyRate:  float32(tradePair.BuyRate),
			SellRate: float32(tradePair.SellRate),
		})
	}
	return pbTradePairs
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

func convertBuyOrder(buyOrder db.BuyOrder) *pb.BuyOrder {
	return &pb.BuyOrder{
		Id:            buyOrder.ID,
		WalletAddress: buyOrder.WalletAddress,
		Crypto:        buyOrder.Crypto,
		Fiat:          buyOrder.Fiat,
		Username:      buyOrder.Username,
		PayAmount:     float32(buyOrder.CryptoAmount),
		ReceiveAmount: float32(buyOrder.FiatAmount),
		SellRate:      float32(buyOrder.Rate),
		IsAccepted:    buyOrder.IsAccepted,
		IsRejected:    buyOrder.IsRejected,
		IsCompleted:   buyOrder.IsCompleted,
		IsExpired:     buyOrder.IsExpired,
		Duration:      timestamppb.New(buyOrder.Duration),
		CreatedAt:     timestamppb.New(buyOrder.CreatedAt),
	}
}

func convertBuyOrders(buyOrders []db.BuyOrder) []*pb.BuyOrder {
	var pbBuyOrders []*pb.BuyOrder
	log.Print(buyOrders)
	for _, buyOrder := range buyOrders {
		pbBuyOrders = append(pbBuyOrders, &pb.BuyOrder{
			Id:            buyOrder.ID,
			WalletAddress: buyOrder.WalletAddress,
			Crypto:        buyOrder.Crypto,
			Fiat:          buyOrder.Fiat,
			Username:      buyOrder.Username,
			PayAmount:     float32(buyOrder.CryptoAmount),
			ReceiveAmount: float32(buyOrder.FiatAmount),
			SellRate:      float32(buyOrder.Rate),
			IsAccepted:    buyOrder.IsAccepted,
			IsRejected:    buyOrder.IsRejected,
			IsCompleted:   buyOrder.IsCompleted,
			IsExpired:     buyOrder.IsExpired,
			Duration:      timestamppb.New(buyOrder.Duration),
			CreatedAt:     timestamppb.New(buyOrder.CreatedAt),
		})
	}

	log.Print(pbBuyOrders, " here")
	return pbBuyOrders
}
