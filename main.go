package main

import (
	"database/sql"
	"log"
	"net"

	_ "github.com/lib/pq"

	// "peerbill-trader-server/api"
	"peerbill-trader-server/api"
	db "peerbill-trader-server/db/sqlc"
	"peerbill-trader-server/gapi"
	"peerbill-trader-server/pb"
	"peerbill-trader-server/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal(err)
	}

	repository := db.NewRepository(conn)
	runGrpcServer(config, repository)

}

func runGrpcServer(config utils.Config, repository db.DatabaseContract) {
	server, err := gapi.NewServer(config, repository)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPeerBillTraderServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("listening...")
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}

func runGinServer(config utils.Config, repository db.DatabaseContract) {
	server, err := api.NewServer(config, repository)
	if err != nil {
		log.Fatal(err)
	}

	err = server.StartServer(config.HTTPServerAddr)
	if err != nil {
		log.Fatal(err)
	}
}
