package api

import (
	"fmt"
	db "peerbill-server/db/sqlc"
	"peerbill-server/token"
	"peerbill-server/utils"

	"github.com/gin-gonic/gin"
)

type Server struct {
	repository db.DatabaseContract
	config     utils.Config
	router     *gin.Engine
	token      token.TokenMaker
}

func NewServer(config utils.Config, r db.DatabaseContract) (*Server, error) {
	token, err := token.NewToken(config.TokenKey)
	if err != nil {
		return nil, fmt.Errorf("unable to create token maker%w", err)
	}

	server := &Server{
		config:     config,
		repository: r,
		token:      token,
	}

	server.setupRouter()
	return server, nil
}

func (s *Server) StartServer(addr string) error {
	return s.router.Run(addr)
}

func (s *Server) setupRouter() {
	router := gin.Default()

	// authRouter := router.Group("/").Use(authMiddleware(s.token))
	router.POST("/api/register-trader", s.RegisterTrader)
	router.POST("/api/login-trader", s.LoginTrader)
	router.POST("/api/refresh-access", s.refreshAccess)
	s.router = router
}

func errorRes(err error) gin.H {
	return gin.H{"err": err}
}
