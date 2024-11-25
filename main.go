package main

import (
	"context"

	// "crypto/tls"
	"database/sql"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"
	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"golang.org/x/sync/errgroup"

	db "peerbill-trader-api/db/sqlc"
	_ "peerbill-trader-api/doc/statik"
	"peerbill-trader-api/gapi"
	"peerbill-trader-api/mail"
	"peerbill-trader-api/pb"
	"peerbill-trader-api/socket"
	"peerbill-trader-api/utils"
	"peerbill-trader-api/worker"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"

	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

/*
*
signals to look out for a graceful shutdown
*/
var signals = []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGINT}

/*
*
set up websocket upgrader
used to upgrade http conn
to websocket. the Upgrade method
will accept rw,r and return a websocket
connection | err
*/
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocket handler function

func runWebSocketServer(group *errgroup.Group, ctx context.Context, config utils.Config) {
	// Define WebSocket handler route
	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		socket.HandleWebSocketConnection(upgrader, config, w, r)
	})

	// Set up your HTTP server
	httpServer := &http.Server{
		Handler: httpMux,
		Addr:    config.WebsocketAddr,
	}

	// Start the WebSocket server in a goroutine
	group.Go(func() error {
		log.Print("WebSocket server listening on ", config.WebsocketAddr)
		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return nil
	})

	// Gracefully shutdown WebSocket server
	group.Go(func() error {
		<-ctx.Done()
		log.Print("Gracefully shutting down WebSocket server...")

		err := httpServer.Shutdown(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		log.Print("WebSocket server shutdown completed.")
		return nil
	})
}

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), signals...)
	defer stop()

	runDBMigration(config.MigrationURL, config.DBSource)
	repository := db.NewRepository(conn)

	// message broker
	redisOption := asynq.RedisClientOpt{
		Addr: config.REDISServerAddr,
	}

	group, ctx := errgroup.WithContext(ctx)

	taskDistributor := worker.NewRedisTaskDistributor(redisOption)
	runGatewayServer(group, ctx, config, repository, taskDistributor)
	runTaskProcessor(group, ctx, redisOption, repository, config)
	runGrpcServer(group, ctx, config, repository, taskDistributor)
	runWebSocketServer(group, ctx, config)

	// wait bfr exiting main fn
	err = group.Wait()
	if err != nil {
		log.Fatal(err)
	}

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

func runGatewayServer(group *errgroup.Group, ctx context.Context, config utils.Config, repository db.DatabaseContract, td worker.TaskDistributor) {
	server, err := gapi.NewServer(config, repository, td)
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
	err = pb.RegisterPeerbillTraderHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal(err)
	}

	// re route client request to the grpc gateway
	httpMux := http.NewServeMux()
	httpMux.Handle("/", grpcMux)

	// load from server
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	httpMux.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(statikFS)))

	c := cors.New(cors.Options{
		AllowedOrigins: config.AllowedOrigins,
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders: []string{
			"Authorization",
			"Content-Type",
		},
	})
	handler := c.Handler(gapi.HttpLogger(httpMux))

	httpServer := &http.Server{
		Handler: handler,
		Addr:    config.HTTPServerAddr,
	}

	group.Go(func() error {
		log.Print("listening...", config.HTTPServerAddr)
		err = httpServer.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}

		return nil
	})

	group.Go(func() error {
		<-ctx.Done()
		log.Print("gracefully shutting down...")

		err = httpServer.Shutdown(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		log.Print("http gateway server shutdown.. goodbye")
		return nil
	})

}

func runTaskProcessor(group *errgroup.Group, ctx context.Context, options asynq.RedisClientOpt, repository db.DatabaseContract, config utils.Config) {
	log.Print("running processor")
	mailer := mail.NewGmailSender(config.EmailSender, config.EmailAddress, config.EmailPassword)
	taskProcessor := worker.NewRedisTaskProcessor(options, repository, mailer)

	err := taskProcessor.Start()
	if err != nil {
		log.Fatal("failed to process tasks")
	}

	group.Go(func() error {
		<-ctx.Done()
		log.Print("gracefully shutting down...")

		taskProcessor.Shutdown()
		log.Print("task processor shutdown.. goodbye")

		return nil
	})
}

func runGrpcServer(group *errgroup.Group, ctx context.Context, config utils.Config, repository db.DatabaseContract, td worker.TaskDistributor) {
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
