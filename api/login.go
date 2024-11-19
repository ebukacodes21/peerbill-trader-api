package api

import (
	"database/sql"
	"net/http"
	"time"

	db "peerbill-trader-server/db/sqlc"
	"peerbill-trader-server/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type loginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginUserResponse struct {
	SessionId              uuid.UUID      `json:"session_id"`
	AccessToken            string         ` json:"access_token"`
	AccessTokenExpiration  time.Time      ` json:"access_token_expiration"`
	RefreshToken           string         `json:"refresh_token"`
	RefreshTokenExpiration time.Time      ` json:"refresh_token_expiration"`
	Trader                 TraderResponse `db:"trader" json:"trader"`
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

	accessToken, accessPayload, err := s.token.CreateToken(trader.Username, s.config.TokenAccess)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorRes(err))
		return
	}

	refreshToken, refreshPayload, err := s.token.CreateToken(trader.Username, s.config.RefreshAccess)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorRes(err))
		return
	}

	session, err := s.repository.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     trader.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiredAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorRes(err))
		return
	}

	resp := loginUserResponse{
		SessionId:              session.ID,
		AccessToken:            accessToken,
		AccessTokenExpiration:  accessPayload.ExpiredAt,
		RefreshToken:           refreshToken,
		RefreshTokenExpiration: refreshPayload.ExpiredAt,
		Trader:                 newTraderResponse(trader),
	}

	ctx.JSON(http.StatusOK, resp)
}
