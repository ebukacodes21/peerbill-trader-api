package db

import "context"

type CreateTraderTxParams struct {
	CreateTraderParams
	AfterCreate func(trader Trader) error
}

type CreateTraderTxResponse struct {
	Trader Trader
}

func (r *Repository) CreateTraderTx(ctx context.Context, args CreateTraderTxParams) (CreateTraderTxResponse, error) {
	var result CreateTraderTxResponse

	err := r.execTx(ctx, func(queries *Queries) error {
		var err error
		result.Trader, err = queries.CreateTrader(ctx, args.CreateTraderParams)
		if err != nil {
			return err
		}

		return args.AfterCreate(result.Trader)
	})

	return result, err
}
