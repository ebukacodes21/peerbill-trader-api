package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"

	// "peerbill-trader-server/api"
	"peerbill-trader-server/api"
	db "peerbill-trader-server/db/sqlc"
	_ "peerbill-trader-server/doc/statik"
	"peerbill-trader-server/gapi"
	"peerbill-trader-server/pb"
	"peerbill-trader-server/utils"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
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
	runDBMigration(config.MigrationURL, config.DBSource)
	go runGatewayServer(config, repository)
	runGrpcServer(config, repository)

}

func runDBMigration(url string, source string) {
	migration, err := migrate.New(url, source)
	if err != nil {
		log.Fatal(err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	log.Print("migration successful")
}

func runGatewayServer(config utils.Config, repository db.DatabaseContract) {
	server, err := gapi.NewServer(config, repository)
	if err != nil {
		log.Fatal(err)
	}

	options := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})
	grpcMux := runtime.NewServeMux(options)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterPeerBillTraderHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal(err)
	}

	// re route client request to the grpc gateway
	httpMux := http.NewServeMux()
	httpMux.Handle("/", grpcMux)

	// fserver := http.FileServer(http.Dir("./doc/swagger"))
	// httpMux.Handle("/swagger/", http.StripPrefix("/swagger/", fserver))

	// load from server
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	httpMux.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(statikFS)))

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
