package db

import (
	"context"
	"log"
	"project-server/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateTrader(t *testing.T) {
	createRandomTrader(t)
}

func createRandomTrader(t *testing.T) Trader {
	password := utils.RandomString(8)
	hash, err := utils.HashPassword(password)
	if err != nil {
		log.Fatal(err)
	}

	arg := CreateTraderParams{
		FirstName: utils.RandomOwner(),
		LastName:  utils.RandomOwner(),
		Username:  utils.RandomOwner(),
		Password:  hash,
		Email:     utils.RandomEmail(),
		Country:   utils.RandomOwner(),
		Phone:     utils.RandomPhone(),
	}

	trader, err := testQueries.CreateTrader(context.Background(), arg)
	require.NoError(t, err)

	return trader
}
