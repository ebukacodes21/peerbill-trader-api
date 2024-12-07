package gapi

import (
	"context"
	"fmt"
	"peerbill-trader-api/token"
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	authorizationBearer = "bearer"
)

func (s *Server) authorizeTrader(ctx context.Context, accessibleRoles []string) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata here")
	}

	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	authType := strings.ToLower(fields[0])
	if authType != authorizationBearer {
		return nil, fmt.Errorf("invalid authorization type")
	}

	accessToken := fields[1]
	payload, err := s.token.VerifyToken(accessToken)
	if err != nil {
		return nil, err
	}

	if !hasAccess(payload.Role, accessibleRoles) {
		return nil, fmt.Errorf("permission denied")
	}

	return payload, nil
}

func hasAccess(role string, accessibleRoles []string) bool {
	for _, v := range accessibleRoles {
		if role == v {
			return true
		}
	}

	return false
}
