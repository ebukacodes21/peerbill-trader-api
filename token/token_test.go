package token

import (
	"testing"
	"time"

	"github.com/ebukacodes21/peerbill-trader-api/utils"

	"github.com/stretchr/testify/require"
)

func TestToken(t *testing.T) {
	token, err := NewToken(utils.RandomString(32))
	require.NoError(t, err)

	username := utils.RandomOwner()
	role := "user"
	duration := time.Minute

	issueAt := time.Now()
	expiredAt := issueAt.Add(duration)

	// create token
	newToken, payload, err := token.CreateToken(username, 1, role, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = token.VerifyToken(newToken)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issueAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredToken(t *testing.T) {
	token, err := NewToken(utils.RandomString(32))
	require.NoError(t, err)

	newToken, payload, err := token.CreateToken(utils.RandomOwner(), 1, "user", -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = token.VerifyToken(newToken)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}
