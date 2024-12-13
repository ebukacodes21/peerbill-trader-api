package db

import (
	"context"
	"testing"

	"github.com/ebukacodes21/peerbill-trader-api/utils"

	"github.com/stretchr/testify/require"
)

func TestCreateTradePair(t *testing.T) {
	trader := createRandomTrader(t)

	args := CreateTradePairParams{
		Username: trader.Username,
		Crypto:   utils.RandomCrypto(),
		Fiat:     utils.RandomFiat(),
		BuyRate:  utils.RandomFloat(20000.00, 50000.00, 2),
		SellRate: utils.RandomFloat(120000.00, 550000.00, 2),
	}

	tradePair, err := testQueries.CreateTradePair(context.Background(), args)
	require.NoError(t, err)
	require.Equal(t, trader.Username, tradePair.Username)
	require.NotEmpty(t, tradePair)
}
