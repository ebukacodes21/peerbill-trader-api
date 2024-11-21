package gapi

import (
	"context"
	db "peerbill-trader-server/db/sqlc"
	"peerbill-trader-server/pb"
	"peerbill-trader-server/utils"
	"peerbill-trader-server/validate"

	pg "github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) RegisterTrader(ctx context.Context, req *pb.RegisterTraderRequest) (*pb.RegisterTraderResponse, error) {
	violations := validateRegisterTraderRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

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

func validateRegisterTraderRequest(req *pb.RegisterTraderRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateFirstname(req.GetFirstName()); err != nil {
		violations = append(violations, fieldViolation(req.GetFirstName(), err))
	}

	if err := validate.ValidateLastname(req.GetLastName()); err != nil {
		violations = append(violations, fieldViolation(req.GetFirstName(), err))
	}

	if err := validate.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation(req.GetUsername(), err))
	}

	if err := validate.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation(req.GetEmail(), err))
	}

	if err := validate.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation(req.GetPassword(), err))
	}

	if err := validate.ValidatePhone(req.GetPhone()); err != nil {
		violations = append(violations, fieldViolation(req.GetPhone(), err))
	}

	if err := validate.ValidateCountry(req.GetCountry()); err != nil {
		violations = append(violations, fieldViolation(req.GetCountry(), err))
	}

	return violations
}
