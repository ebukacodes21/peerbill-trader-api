package main

import (
	"context"
	"encoding/json"

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

	// "peerbill-server/api"
	"peerbill-server/api"
	db "peerbill-server/db/sqlc"
	_ "peerbill-server/doc/statik"
	"peerbill-server/gapi"
	"peerbill-server/mail"
	"peerbill-server/pb"
	"peerbill-server/utils"
	"peerbill-server/worker"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"

	// "google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

var signals = []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGINT}
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocket handler function
func handleWebSocketConnection(config utils.Config, w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}
	defer conn.Close()

	// Connect to gRPC server securely (with TLS)
	grpcConn, err := grpc.NewClient(config.GRPCServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("Failed to connect to gRPC server:", err)
		return
	}
	defer grpcConn.Close()
	client := pb.NewPeerbillClient(grpcConn)

	// Start the bidirectional stream using SubscribeRate
	stream, err := client.SubscribeRate(r.Context())
	if err != nil {
		log.Println("Error starting SubscribeRate stream:", err)
		return
	}

	// Goroutine for reading gRPC responses and sending them to the WebSocket
	go func() {
		for {
			resp, err := stream.Recv() // Receive from the gRPC stream
			if err != nil {
				log.Println("Error receiving from gRPC stream:", err)
				break
			}

			// Marshal the message to JSON
			data, err := json.Marshal(resp)
			if err != nil {
				log.Println("Error marshaling response to JSON:", err)
				break
			}

			// Send the response back to the WebSocket client
			err = conn.WriteMessage(websocket.TextMessage, data) // Send the marshaled JSON response
			if err != nil {
				log.Println("Error writing to WebSocket:", err)
				break
			}
		}
	}()

	// Continuously read messages from the WebSocket
	for {
		// Read a message from the WebSocket client
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading from WebSocket:", err)
			break
		}

		// Use an anonymous struct to unmarshal the message
		var message struct {
			Fiat   string `json:"fiat"`
			Crypto string `json:"crypto"`
		}

		// Unmarshal the JSON byte slice into the anonymous struct
		err = json.Unmarshal(msg, &message)
		if err != nil {
			log.Println("Error unmarshaling JSON:", err)
			continue
		}

		// Send the WebSocket message as a SubscribeRateRequest to the gRPC server
		err = stream.Send(&pb.SubscribeRateRequest{
			Crypto: message.Crypto,
			Fiat:   message.Fiat,
		})
		if err != nil {
			log.Println("Error sending message to gRPC server:", err)
			continue
		}
	}

	// When the WebSocket connection is closed, cancel the gRPC stream
	conn.SetCloseHandler(func(code int, text string) error {
		log.Println("WebSocket connection closed. Cancelling gRPC stream.")
		stream.CloseSend() // Close the gRPC stream
		return nil
	})
}

func runWebSocketServer(group *errgroup.Group, ctx context.Context, config utils.Config) {
	// Define WebSocket handler route
	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocketConnection(config, w, r)
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
	err = pb.RegisterPeerbillHandlerServer(ctx, grpcMux, server)
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
	pb.RegisterPeerbillServer(grpcServer, server)
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
