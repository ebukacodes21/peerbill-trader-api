package gapi

import (
	"context"
	"database/sql"

	db "peerbill-trader-api/db/sqlc"
	"peerbill-trader-api/pb"
	"peerbill-trader-api/utils"
	"peerbill-trader-api/validate"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) LoginTrader(ctx context.Context, req *pb.LoginTraderRequest) (*pb.LoginTraderResponse, error) {
	violations := validateLoginTraderRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	trader, err := s.repository.GetTrader(ctx, req.GetUsername())
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

	if !trader.IsVerified {
		return nil, status.Errorf(codes.Unauthenticated, "trader yet not verified")
	}

	accessToken, accessPayload, err := s.token.CreateToken(trader.Username, trader.ID, trader.Role, s.config.TokenAccess)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to create access token")
	}

	refreshToken, refreshPayload, err := s.token.CreateToken(trader.Username, trader.ID, trader.Role, s.config.RefreshAccess)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to create refresh token")
	}

	metaData := s.extractMetaData(ctx)
	session, err := s.repository.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     trader.Username,
		TraderID:     trader.ID,
		RefreshToken: refreshToken,
		UserAgent:    metaData.UserAgent,
		ClientIp:     metaData.ClientIp,
		IsBlocked:    false,
		ExpiredAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session")
	}

	// attach trade pairs to response
	traderPairs, err := s.repository.GetTraderPairs(ctx, trader.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch trader pairs")
	}

	methods, err := s.repository.GetPaymentMethods(ctx, trader.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch payment methods")
	}

	resp := &pb.LoginTraderResponse{
		Trader:                convert(trader),
		TraderPairs:           convertTradePairs(traderPairs),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
		SessionId:             session.ID.String(),
		PaymentMethods:        convertPaymentMethods(methods),
	}

	return resp, nil
}

func validateLoginTraderRequest(req *pb.LoginTraderRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := validate.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return violations
}
