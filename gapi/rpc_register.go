package gapi

import (
	"context"
	db "peerbill-trader-server/db/sqlc"
	"peerbill-trader-server/pb"
	"peerbill-trader-server/utils"

	pg "github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) RegisterTrader(ctx context.Context, req *pb.RegisterTraderRequest) (*pb.RegisterTraderResponse, error) {
	hash, err := utils.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password")
	}

	args := db.CreateTraderParams{
		FirstName: req.GetFirstName(),
		LastName:  req.GetLastName(),
		Username:  req.GetUsername(),
		Password:  hash,
		Email:     req.GetEmail(),
		Country:   req.GetCountry(),
		Phone:     req.GetPhone(),
	}

	trader, err := s.repository.CreateTrader(ctx, args)
	if err != nil {
		if pgErr, ok := err.(*pg.Error); ok {
			switch pgErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.Internal, "already exists")
			}
		}

		return nil, status.Errorf(codes.Internal, "failed to create user")
	}

	resp := &pb.RegisterTraderResponse{
		Trader: convert(trader),
	}

	return resp, nil
}
