package gapi

import (
	"fmt"
	db "peerbill-server/db/sqlc"
	"peerbill-server/pb"
	"peerbill-server/token"
	"peerbill-server/utils"
	"peerbill-server/worker"
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
