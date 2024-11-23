package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type refreshAccessRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type refreshAccessResponse struct {
	AccessToken           string    ` json:"access_token"`
	AccessTokenExpiration time.Time ` json:"access_token_expiration"`
}

func (s *Server) refreshAccess(ctx *gin.Context) {
	var refreshReq refreshAccessRequest
	if err := ctx.ShouldBindJSON(&refreshReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorRes(err))
		return
	}

	payload, err := s.token.VerifyToken(refreshReq.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorRes(err))
		return
	}

	session, err := s.repository.GetSession(ctx, payload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorRes(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorRes(err))
		return
	}

	if session.IsBlocked {
		err := fmt.Errorf("blocked session")
		ctx.JSON(http.StatusUnauthorized, errorRes(err))
		return
	}

	if session.Username != payload.Username {
		err := fmt.Errorf("incorrect user")
		ctx.JSON(http.StatusUnauthorized, errorRes(err))
		return
	}

	if session.RefreshToken != refreshReq.RefreshToken {
		err := fmt.Errorf("mismatch session")
		ctx.JSON(http.StatusUnauthorized, errorRes(err))
		return
	}

	if time.Now().After(session.ExpiredAt) {
		err := fmt.Errorf("expired session")
		ctx.JSON(http.StatusUnauthorized, errorRes(err))
		return
	}

	accessToken, accessPayload, err := s.token.CreateToken(payload.Username, payload.Role, s.config.TokenAccess)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorRes(err))
		return
	}

	resp := refreshAccessResponse{
		AccessToken:           accessToken,
		AccessTokenExpiration: accessPayload.ExpiredAt,
	}

	ctx.JSON(http.StatusOK, resp)
}
