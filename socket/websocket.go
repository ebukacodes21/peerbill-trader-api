package socket

import (
	// "encoding/json"
	// "log"
	"net/http"
	// "peerbill-trader-api/pb"
	"peerbill-trader-api/utils"

	"github.com/gorilla/websocket"
	// "google.golang.org/grpc"
	// "google.golang.org/grpc/credentials/insecure"
)

func HandleWebSocketConnection(upgrader websocket.Upgrader, config utils.Config, w http.ResponseWriter, r *http.Request) {
	// conn, err := upgrader.Upgrade(w, r, nil)
	// if err != nil {
	// 	log.Println("Failed to upgrade connection:", err)
	// 	return
	// }
	// defer conn.Close()

	// // Connect to gRPC server securely (with TLS)
	// grpcConn, err := grpc.NewClient(config.GRPCServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// if err != nil {
	// 	log.Println("Failed to connect to gRPC server:", err)
	// 	return
	// }
	// defer grpcConn.Close()
	// client := pb.NewPeerbillTraderClient(grpcConn)

	// // Start the bidirectional stream using SubscribeRate
	// stream, err := client.SubscribeRate(r.Context())
	// if err != nil {
	// 	log.Println("Error starting SubscribeRate stream:", err)
	// 	return
	// }

	// // Goroutine for reading gRPC responses and sending them to the WebSocket
	// go func() {
	// 	for {
	// 		resp, err := stream.Recv() // Receive from the gRPC stream
	// 		if err != nil {
	// 			log.Println("Error receiving from gRPC stream:", err)
	// 			break
	// 		}

	// 		// Marshal the message to JSON
	// 		data, err := json.Marshal(resp)
	// 		if err != nil {
	// 			log.Println("Error marshaling response to JSON:", err)
	// 			break
	// 		}

	// 		// Send the response back to the WebSocket client
	// 		err = conn.WriteMessage(websocket.TextMessage, data) // Send the marshaled JSON response
	// 		if err != nil {
	// 			log.Println("Error writing to WebSocket:", err)
	// 			break
	// 		}
	// 	}
	// }()

	// // Continuously read messages from the WebSocket
	// for {
	// 	// Read a message from the WebSocket client
	// 	_, msg, err := conn.ReadMessage()
	// 	if err != nil {
	// 		log.Println("Error reading from WebSocket:", err)
	// 		break
	// 	}

	// 	// Use an anonymous struct to unmarshal the message
	// 	var message struct {
	// 		Fiat   string `json:"fiat"`
	// 		Crypto string `json:"crypto"`
	// 	}

	// 	// Unmarshal the JSON byte slice into the anonymous struct
	// 	err = json.Unmarshal(msg, &message)
	// 	if err != nil {
	// 		log.Println("Error unmarshaling JSON:", err)
	// 		continue
	// 	}

	// 	// Send the WebSocket message as a SubscribeRateRequest to the gRPC server
	// 	err = stream.Send(&pb.SubscribeRateRequest{
	// 		Crypto: message.Crypto,
	// 		Fiat:   message.Fiat,
	// 	})
	// 	if err != nil {
	// 		log.Println("Error sending message to gRPC server:", err)
	// 		continue
	// 	}
	// }

	// // When the WebSocket connection is closed, cancel the gRPC stream
	// conn.SetCloseHandler(func(code int, text string) error {
	// 	log.Println("WebSocket connection closed. Cancelling gRPC stream.")
	// 	stream.CloseSend() // Close the gRPC stream
	// 	return nil
	// })
}
