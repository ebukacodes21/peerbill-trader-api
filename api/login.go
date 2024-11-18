package api

import (
	"database/sql"
	"net/http"

	"peerbill-trader-server/utils"

	"github.com/gin-gonic/gin"
)

type loginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginUserResponse struct {
	AccessToken string         `db:"access_token" json:"access_token"`
	Trader      TraderResponse `db:"trader" json:"trader"`
}

func (s *Server) LoginTrader(ctx *gin.Context) {
	var loginReq loginUserRequest
	if err := ctx.ShouldBindJSON(&loginReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorRes(err))
		return
	}

	trader, err := s.repository.GetTrader(ctx, loginReq.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorRes(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorRes(err))
		return
	}

	err = utils.VerifyPassword(trader.Password, loginReq.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorRes(err))
		return
	}

	accessToken, err := s.token.CreateToken(trader.Username, s.config.TokenAccess)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorRes(err))
		return
	}

	resp := loginUserResponse{
		AccessToken: accessToken,
		Trader:      newTraderResponse(trader),
	}

	ctx.JSON(http.StatusOK, resp)
}
