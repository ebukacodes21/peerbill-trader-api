package api

import (
	"net/http"
	db "peerbill-trader-server/db/sqlc"
	"peerbill-trader-server/utils"
	"time"

	"github.com/gin-gonic/gin"
	pg "github.com/lib/pq"
)

type TraderRequest struct {
	FirstName string `db:"first_name" json:"first_name" binding:"required"`
	LastName  string `db:"last_name" json:"last_name" binding:"required"`
	Username  string `db:"username" json:"username" binding:"required,alphanum"`
	Password  string `db:"password" json:"password" binding:"required"`
	Email     string `db:"email" json:"email" binding:"required,email"`
	Country   string `db:"country" json:"country" binding:"required"`
	Phone     string `db:"phone" json:"phone" binding:"required"`
}

type TraderResponse struct {
	FirstName string    `db:"first_name" json:"first_name"`
	LastName  string    `db:"last_name" json:"last_name"`
	Username  string    `db:"username" json:"username"`
	Email     string    `db:"email" json:"email"`
	Country   string    `db:"country" json:"country"`
	Phone     string    `db:"phone" json:"phone"`
	CREATEDAT time.Time `db:"created_at" json:"created_at"`
}

func newTraderResponse(trader db.Trader) TraderResponse {
	return TraderResponse{
		FirstName: trader.FirstName,
		LastName:  trader.LastName,
		Username:  trader.Username,
		Email:     trader.Email,
		Country:   trader.Country,
		Phone:     trader.Phone,
		CREATEDAT: trader.CreatedAt,
	}
}

func (s *Server) RegisterTrader(ctx *gin.Context) {
	var traderReq TraderRequest
	if err := ctx.ShouldBindJSON(&traderReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorRes(err))
		return
	}

	hash, err := utils.HashPassword(traderReq.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorRes(err))
		return
	}

	args := db.CreateTraderParams{
		FirstName: traderReq.FirstName,
		LastName:  traderReq.LastName,
		Username:  traderReq.Username,
		Password:  hash,
		Email:     traderReq.Email,
		Country:   traderReq.Country,
		Phone:     traderReq.Phone,
	}

	trader, err := s.repository.CreateTrader(ctx, args)
	if err != nil {
		if pgErr, ok := err.(*pg.Error); ok {
			switch pgErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorRes(pgErr))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, errorRes(err))
		return
	}

	traderRes := newTraderResponse(trader)
	ctx.JSON(http.StatusCreated, traderRes)
}
