package db

import (
	"context"
	"database/sql"
)

type VerifyEmailTxParams struct {
	EmailId    int64
	SecretCode string
}

type VerifyEmailTxResponse struct {
	Trader      Trader
	VerifyEmail VerifyEmail
}

func (r *Repository) VerifyEmailTx(ctx context.Context, args VerifyEmailTxParams) (VerifyEmailTxResponse, error) {
	var result VerifyEmailTxResponse

	err := r.execTx(ctx, func(queries *Queries) error {
		var err error
		result.VerifyEmail, err = queries.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID:         args.EmailId,
			SecretCode: args.SecretCode,
		})

		if err != nil {
			return err
		}

		result.Trader, err = queries.UpdateTrader(ctx, UpdateTraderParams{
			ID: result.VerifyEmail.UserID,
			IsVerified: sql.NullBool{
				Bool:  true,
				Valid: true,
			},
		})

		return err
	})
	return result, err
}
