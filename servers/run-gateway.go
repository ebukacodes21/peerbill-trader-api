package servers

import (
	"context"
	"log"
	"net/http"

	db "github.com/ebukacodes21/peerbill-trader-api/db/sqlc"
	"github.com/ebukacodes21/peerbill-trader-api/gapi"
	"github.com/ebukacodes21/peerbill-trader-api/pb"
	"github.com/ebukacodes21/peerbill-trader-api/utils"
	"github.com/ebukacodes21/peerbill-trader-api/worker"

	_ "github.com/ebukacodes21/peerbill-trader-api/doc/statik"

	"github.com/rakyll/statik/fs"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/encoding/protojson"
)

func RunGatewayServer(group *errgroup.Group, ctx context.Context, config utils.Config, repository db.DatabaseContract, td worker.TaskDistributor, tp worker.TaskProcessor) {
	server, err := gapi.NewServer(config, repository, td, tp)
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
		log.Print("Gateway server running on ", config.HTTPServerAddr)
		err = httpServer.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}

		return nil
	})

	group.Go(func() error {
		<-ctx.Done()
		log.Print("gateway gracefully shutting down...")

		err = httpServer.Shutdown(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		log.Print("goodbye")
		return nil
	})

}
