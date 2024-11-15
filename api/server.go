package api

import (
	db "peerbill-trader-server/db/sqlc"
	"peerbill-trader-server/utils"

	"github.com/gin-gonic/gin"
)

type Server struct {
	repository db.DatabaseContract
	config     utils.Config
	router     *gin.Engine
}

func NewServer(config utils.Config, r db.DatabaseContract) (*Server, error) {
	server := &Server{config: config, repository: r}

	server.setupRouter()
	return server, nil
}

func (s *Server) StartServer(addr string) error {
	return s.router.Run(addr)
}

func (s *Server) setupRouter() {
	router := gin.Default()
	router.POST("/api/register-trader", s.RegisterTrader)
	s.router = router
}

func errorRes(err error) gin.H {
	return gin.H{"err": err}
}
