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

func convertTraders(traders []db.Trader) []*pb.Trader {
	var pbTraders []*pb.Trader
	for _, trader := range traders {
		pbTrader := &pb.Trader{
			FirstName: trader.FirstName,
			LastName:  trader.LastName,
			Username:  trader.Username,
			Email:     trader.Email,
			Country:   trader.Country,
			Phone:     trader.Phone,
			CreatedAt: timestamppb.New(trader.CreatedAt),
		}
		pbTraders = append(pbTraders, pbTrader)
	}
	return pbTraders
}
