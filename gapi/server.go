package gapi

import (
	"fmt"

	db "github.com/ebukacodes21/peerbill-trader-api/db/sqlc"
	"github.com/ebukacodes21/peerbill-trader-api/pb"
	"github.com/ebukacodes21/peerbill-trader-api/token"
	"github.com/ebukacodes21/peerbill-trader-api/utils"
	"github.com/ebukacodes21/peerbill-trader-api/worker"
)

type Server struct {
	pb.UnimplementedPeerbillTraderServer
	repository      db.DatabaseContract
	token           token.TokenMaker
	config          utils.Config
	taskDistributor worker.TaskDistributor
	taskProcessor   worker.TaskProcessor
}

func NewServer(config utils.Config, r db.DatabaseContract, taskDistributor worker.TaskDistributor, taskProcessor worker.TaskProcessor) (*Server, error) {
	token, err := token.NewToken(config.TokenKey)
	if err != nil {
		return nil, fmt.Errorf("unable to create token maker%w", err)
	}

	server := &Server{
		config:          config,
		repository:      r,
		token:           token,
		taskDistributor: taskDistributor,
		taskProcessor:   taskProcessor,
	}

	return server, nil
}
