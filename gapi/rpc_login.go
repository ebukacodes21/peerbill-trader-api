package gapi

import (
	"context"
	"database/sql"
	"log"
	db "peerbill-trader-server/db/sqlc"
	"peerbill-trader-server/pb"
	"peerbill-trader-server/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) LoginTrader(ctx context.Context, req *pb.LoginTraderRequest) (*pb.LoginTraderResponse, error) {
	trader, err := s.repository.GetTrader(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "trader not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find trader")
	}

	err = utils.VerifyPassword(trader.Password, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "incorrect password")
	}

	accessToken, accessPayload, err := s.token.CreateToken(trader.Username, s.config.TokenAccess)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to create access token")
	}

	refreshToken, refreshPayload, err := s.token.CreateToken(trader.Username, s.config.RefreshAccess)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to create refresh token")
	}

	session, err := s.repository.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     trader.Username,
		RefreshToken: refreshToken,
		UserAgent:    "",
		ClientIp:     "",
		IsBlocked:    false,
		ExpiredAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		log.Print(err)
		return nil, status.Errorf(codes.Internal, "failed to create session")
	}

	resp := &pb.LoginTraderResponse{
		Trader:                convert(trader),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
		SessionId:             session.ID.String(),
	}

	return resp, nil
}