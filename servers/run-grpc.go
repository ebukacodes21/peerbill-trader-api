package servers

import (
	"context"
	"log"
	"net"
	db "peerbill-trader-api/db/sqlc"
	"peerbill-trader-api/gapi"
	"peerbill-trader-api/pb"
	"peerbill-trader-api/utils"
	"peerbill-trader-api/worker"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func RunGrpcServer(group *errgroup.Group, ctx context.Context, config utils.Config, repository db.DatabaseContract, td worker.TaskDistributor) {
	server, err := gapi.NewServer(config, repository, td)
	if err != nil {
		log.Fatal(err)
	}

	logger := grpc.UnaryInterceptor(gapi.Logger)

	grpcServer := grpc.NewServer(logger)
	pb.RegisterPeerbillTraderServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddr)
	if err != nil {
		log.Fatal(err)
	}

	// starting server in go routine
	group.Go(func() error {
		log.Print("listening...", config.GRPCServerAddr)
		err = grpcServer.Serve(listener)
		if err != nil {
			log.Fatal(err)
		}

		return nil
	})

	// graceful shutdown
	group.Go(func() error {
		<-ctx.Done()
		log.Print("gracefully shutting down...")

		grpcServer.GracefulStop()
		log.Print("grpc server shutdown.. goodbye")

		return nil
	})
}
