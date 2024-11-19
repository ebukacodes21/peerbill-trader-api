package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	_ "github.com/lib/pq"

	// "peerbill-trader-server/api"
	"peerbill-trader-server/api"
	db "peerbill-trader-server/db/sqlc"
	"peerbill-trader-server/gapi"
	"peerbill-trader-server/pb"
	"peerbill-trader-server/utils"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
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
	go runGatewayServer(config, repository)
	runGrpcServer(config, repository)

}

func runGatewayServer(config utils.Config, repository db.DatabaseContract) {
	server, err := gapi.NewServer(config, repository)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	grpcMux := runtime.NewServeMux()
	err = pb.RegisterPeerBillTraderHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal(err)
	}

	// re route client request to the grpc gateway
	httpMux := http.NewServeMux()
	httpMux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.HTTPServerAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("listening...", config.HTTPServerAddr)
	err = http.Serve(listener, httpMux)
	if err != nil {
		log.Fatal(err)
	}
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
	log.Print("listening...", config.GRPCServerAddr)
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
