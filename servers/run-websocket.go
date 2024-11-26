package servers

import (
	"context"
	"log"
	"net/http"
	"peerbill-trader-api/socket"
	"peerbill-trader-api/utils"

	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
)

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

func RunWebSocketServer(group *errgroup.Group, ctx context.Context, config utils.Config) {
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
