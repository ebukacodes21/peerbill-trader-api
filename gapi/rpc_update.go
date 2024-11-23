package gapi

import (
	"context"
	"database/sql"
	db "peerbill-trader-server/db/sqlc"
	"peerbill-trader-server/pb"
	"peerbill-trader-server/utils"
	"peerbill-trader-server/validate"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateTrader(ctx context.Context, req *pb.UpdateTraderRequest) (*pb.UpdateTraderResponse, error) {
	authPayload, err := s.authorizeTrader(ctx, []string{"user", "admin"})
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateUpdateTraderRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	if authPayload.Role != "admin" && authPayload.Username != *req.Username {
		return nil, status.Errorf(codes.PermissionDenied, "cannot update user info: %s", err)
	}

	args := db.UpdateTraderParams{
		ID: req.GetId(),
		FirstName: sql.NullString{
			String: req.GetFirstName(),
			Valid:  req.FirstName != nil,
		},
		LastName: sql.NullString{
			String: req.GetLastName(),
			Valid:  req.LastName != nil,
		},
		Username: sql.NullString{
			String: req.GetUsername(),
			Valid:  req.Username != nil,
		},
		Email: sql.NullString{
			String: req.GetEmail(),
			Valid:  req.Email != nil,
		},
		Country: sql.NullString{
			String: req.GetCountry(),
			Valid:  req.Country != nil,
		},
		Phone: sql.NullString{
			String: req.GetPhone(),
			Valid:  req.Phone != nil,
		},
	}

	if req.Password != nil {
		hash, err := utils.HashPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password")
		}

		args.Password = sql.NullString{
			String: hash,
			Valid:  true,
		}
	}

	trader, err := s.repository.UpdateTrader(ctx, args)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}

		return nil, status.Errorf(codes.Internal, "failed to update user")
	}

	resp := &pb.UpdateTraderResponse{
		Trader: convert(trader),
	}

	return resp, nil
}

func validateUpdateTraderRequest(req *pb.UpdateTraderRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if req.FirstName != nil {
		if err := validate.ValidateFirstname(req.GetFirstName()); err != nil {
			violations = append(violations, fieldViolation(req.GetFirstName(), err))
		}
	}

	if req.LastName != nil {
		if err := validate.ValidateLastname(req.GetLastName()); err != nil {
			violations = append(violations, fieldViolation(req.GetFirstName(), err))
		}
	}

	if req.Username != nil {
		if err := validate.ValidateUsername(req.GetUsername()); err != nil {
			violations = append(violations, fieldViolation(req.GetUsername(), err))
		}
	}

	if req.Email != nil {
		if err := validate.ValidateEmail(req.GetEmail()); err != nil {
			violations = append(violations, fieldViolation(req.GetEmail(), err))
		}
	}

	if req.Password != nil {
		if err := validate.ValidatePassword(req.GetPassword()); err != nil {
			violations = append(violations, fieldViolation(req.GetPassword(), err))
		}
	}

	if req.Phone != nil {
		if err := validate.ValidatePhone(req.GetPhone()); err != nil {
			violations = append(violations, fieldViolation(req.GetPhone(), err))
		}
	}
	if req.Country != nil {

		if err := validate.ValidateCountry(req.GetCountry()); err != nil {
			violations = append(violations, fieldViolation(req.GetCountry(), err))
		}
	}

	return violations
}
