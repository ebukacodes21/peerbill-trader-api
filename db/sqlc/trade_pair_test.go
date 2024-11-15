package db

import (
	"context"
	"peerbill-trader-server/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateTradePair(t *testing.T) {
	trader := createRandomTrader(t)

	args := CreateTradePairParams{
		TraderID:   trader.ID,
		BaseAsset:  utils.RandomFiat(),
		QuoteAsset: utils.RandomCrypto(),
		BuyRate:    utils.RandomFloat(20000.00, 50000.00, 2),
		SellRate:   utils.RandomFloat(120000.00, 550000.00, 2),
	}

	tradePair, err := testQueries.CreateTradePair(context.Background(), args)
	require.NoError(t, err)
	require.Equal(t, trader.ID, tradePair.TraderID)
	require.NotEmpty(t, tradePair)
}
