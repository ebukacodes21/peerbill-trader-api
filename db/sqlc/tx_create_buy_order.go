package db

import "context"

type CreateBuyOrderTxParams struct {
	CreateBuyOrderParams
	AfterCreate func(buyOrder BuyOrder) error
}

type CreatBuyOrderTxResponse struct {
	BuyOrder BuyOrder
}

func (r *Repository) CreateBuyOrderTx(ctx context.Context, args CreateBuyOrderTxParams) (CreatBuyOrderTxResponse, error) {
	var result CreatBuyOrderTxResponse

	err := r.execTx(ctx, func(queries *Queries) error {
		var err error
		result.BuyOrder, err = queries.CreateBuyOrder(ctx, args.CreateBuyOrderParams)
		if err != nil {
			return err
		}

		return args.AfterCreate(result.BuyOrder)
	})

	return result, err
}
