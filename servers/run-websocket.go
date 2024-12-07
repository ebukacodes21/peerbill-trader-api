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

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func RunWebSocketServer(group *errgroup.Group, ctx context.Context, config utils.Config, manager *socket.WebSocketManager) {
	httpMux := http.NewServeMux()

	handleWebSocket := func(w http.ResponseWriter, r *http.Request, subscription string) {
		// Upgrade the HTTP connection to a WebSocket connection
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Failed to upgrade connection:", err)
			return
		}
		defer conn.Close()

		// Add the new client to the manager
		manager.AddClient(conn, subscription)

		// Handle the WebSocket connection
		for {
			// keep alive 🤔
			_, _, err := conn.ReadMessage()
			if err != nil {
				// If there's an error (e.g., the client disconnects), remove the client
				manager.RemoveClient(conn)
				break
			}
		}
	}

	httpMux.HandleFunc("/ws/get-buy-orders", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(w, r, "get-buy-orders")
	})

	httpMux.HandleFunc("/ws/reject-buy-order", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(w, r, "reject-buy-order")
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
