package gapi

import (
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

func convertPaymentMethods(methods []db.PaymentMethod) []*pb.PaymentMethod {
	var pbMethods []*pb.PaymentMethod
	for _, method := range methods {
		pbMethods = append(pbMethods, &pb.PaymentMethod{
			Id:            method.ID,
			Crypto:        method.Crypto,
			Fiat:          method.Fiat,
			Username:      method.Username,
			AccountHolder: method.AccountHolder,
			AccountNumber: method.AccountNumber,
			BankName:      method.BankName,
			WalletAddress: method.WalletAddress,
		})
	}
	return pbMethods
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
		FiatAmount:    float32(buyOrder.FiatAmount),
		CryptoAmount:  float32(buyOrder.CryptoAmount),
		Rate:          float32(buyOrder.Rate),
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
	for _, buyOrder := range buyOrders {
		pbBuyOrders = append(pbBuyOrders, &pb.BuyOrder{
			Id:            buyOrder.ID,
			WalletAddress: buyOrder.WalletAddress,
			Crypto:        buyOrder.Crypto,
			Fiat:          buyOrder.Fiat,
			Username:      buyOrder.Username,
			FiatAmount:    float32(buyOrder.FiatAmount),
			CryptoAmount:  float32(buyOrder.CryptoAmount),
			Rate:          float32(buyOrder.Rate),
			IsAccepted:    buyOrder.IsAccepted,
			IsRejected:    buyOrder.IsRejected,
			IsCompleted:   buyOrder.IsCompleted,
			IsExpired:     buyOrder.IsExpired,
			Duration:      timestamppb.New(buyOrder.Duration),
			CreatedAt:     timestamppb.New(buyOrder.CreatedAt),
		})
	}

	return pbBuyOrders
}
