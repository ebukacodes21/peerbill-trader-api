package db

import "context"

type CreateOrderTxParams struct {
	CreateOrderParams
	AfterCreate func(order Order) error
}

type CreateOrderTxResponse struct {
	Order Order
}

func (r *Repository) CreateOrderTx(ctx context.Context, args CreateOrderTxParams) (CreateOrderTxResponse, error) {
	var result CreateOrderTxResponse

	err := r.execTx(ctx, func(queries *Queries) error {
		var err error
		result.Order, err = queries.CreateOrder(ctx, args.CreateOrderParams)
		if err != nil {
			return err
		}

		return args.AfterCreate(result.Order)
	})

	return result, err
}
