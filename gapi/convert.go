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

func convertPaymentMethod(method db.PaymentMethod) *pb.PaymentMethod {
	return &pb.PaymentMethod{
		Id:            method.ID,
		Crypto:        method.Crypto,
		Fiat:          method.Fiat,
		Username:      method.Username,
		AccountHolder: method.AccountHolder,
		AccountNumber: method.AccountNumber,
		BankName:      method.BankName,
		WalletAddress: method.WalletAddress,
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

func convertOrder(order db.Order) *pb.Order {

	return &pb.Order{
		Id:            order.ID,
		EscrowAddress: order.EscrowAddress,
		UserAddress:   order.UserAddress,
		OrderType:     order.OrderType,
		Crypto:        order.Crypto,
		Fiat:          order.Fiat,
		Username:      order.Username,
		FiatAmount:    float32(order.FiatAmount),
		CryptoAmount:  float32(order.CryptoAmount),
		Rate:          float32(order.Rate),
		IsAccepted:    order.IsAccepted,
		IsRejected:    order.IsRejected,
		IsCompleted:   order.IsCompleted,
		IsExpired:     order.IsExpired,
		Duration:      timestamppb.New(order.Duration),
		CreatedAt:     timestamppb.New(order.CreatedAt),
	}
}

func convertOrders(Orders []db.Order) []*pb.Order {
	var pbOrders []*pb.Order
	for _, order := range Orders {
		pbOrders = append(pbOrders, &pb.Order{
			Id:            order.ID,
			EscrowAddress: order.EscrowAddress,
			UserAddress:   order.UserAddress,
			OrderType:     order.OrderType,
			Crypto:        order.Crypto,
			Fiat:          order.Fiat,
			Username:      order.Username,
			FiatAmount:    float32(order.FiatAmount),
			CryptoAmount:  float32(order.CryptoAmount),
			Rate:          float32(order.Rate),
			IsAccepted:    order.IsAccepted,
			IsRejected:    order.IsRejected,
			IsCompleted:   order.IsCompleted,
			IsExpired:     order.IsExpired,
			Duration:      timestamppb.New(order.Duration),
			CreatedAt:     timestamppb.New(order.CreatedAt),
		})
	}

	return pbOrders
}
