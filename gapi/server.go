package gapi

import (
	"fmt"
	db "peerbill-trader-server/db/sqlc"
	"peerbill-trader-server/pb"
	"peerbill-trader-server/token"
	"peerbill-trader-server/utils"
	"peerbill-trader-server/worker"
)

type Server struct {
	pb.UnimplementedPeerBillTraderServer
	repository      db.DatabaseContract
	token           token.TokenMaker
	config          utils.Config
	taskDistributor worker.TaskDistributor
}

func NewServer(config utils.Config, r db.DatabaseContract, taskDistributor worker.TaskDistributor) (*Server, error) {
	token, err := token.NewToken(config.TokenKey)
	if err != nil {
		return nil, fmt.Errorf("unable to create token maker%w", err)
	}

	server := &Server{
		config:          config,
		repository:      r,
		token:           token,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
