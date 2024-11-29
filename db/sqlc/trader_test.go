package db

import (
	"context"
	"log"
	"peerbill-trader-api/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateTrader(t *testing.T) {
	createRandomTrader(t)
}

func createRandomTrader(t *testing.T) Trader {
	password := utils.RandomString(8)
	code := utils.RandomString(32)
	hash, err := utils.HashPassword(password)
	if err != nil {
		log.Fatal(err)
	}

	arg := CreateTraderParams{
		FirstName:        utils.RandomOwner(),
		LastName:         utils.RandomOwner(),
		Username:         utils.RandomOwner(),
		Password:         hash,
		Email:            utils.RandomEmail(),
		Country:          utils.RandomOwner(),
		Phone:            utils.RandomPhone(),
		VerificationCode: code,
	}

	trader, err := testQueries.CreateTrader(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, trader)

	return trader
}
