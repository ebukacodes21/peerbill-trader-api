package gapi

import (
	"fmt"
	db "peerbill-trader-api/db/sqlc"
	"peerbill-trader-api/pb"
	"peerbill-trader-api/token"
	"peerbill-trader-api/utils"
	"peerbill-trader-api/worker"
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
