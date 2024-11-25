package gapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"peerbill-trader-api/pb"

	"google.golang.org/grpc"
)

type RateRequest struct {
	Fiat   string `json:"fiat"`
	Crypto string `json:"crypto"`
}

// Fetch the exchange rate from the given URL
func fetchExchangeRate(url string, rr *pb.SubscribeRateRequest) (float32, error) {
	// Construct the API URL with the crypto and fiat symbols
	apiURL := fmt.Sprintf("%sfsym=%s&tsyms=%s", url, rr.Crypto, rr.Fiat)

	// Fetch the rate from the API
	resp, err := http.Get(apiURL)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch exchange rate from API: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %v", err)
	}

	// Unmarshal the response into a map of fiat currencies to rates
	var data map[string]float32
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	// // Get the rate for the requested fiat currency (fiat is a string like "NGN")
	rateStr, ok := data[rr.Fiat]
	if !ok {
		return 0, fmt.Errorf("rate not found in response for fiat %s", rr.Fiat)
	}

	return rateStr, nil
}

// SubscribeRate is the method for bidirectional gRPC streaming
func (s *Server) SubscribeRate(src grpc.BidiStreamingServer[pb.SubscribeRateRequest, pb.SubscribeRateResponse]) error {
	// Define the base URL for the CryptoCompare API
	rateAPIURL := "https://min-api.cryptocompare.com/data/price?"

	// Loop to continuously process requests and send exchange rates
	for {
		// Receive the latest fiat/crypto pair from the client
		rr, err := src.Recv()
		if err != nil {
			log.Println("Error receiving from client:", err)
			return err // Return error if receiving fails
		}

		// Log received data for debugging
		log.Print("Received request for ", rr)

		// Use a goroutine to fetch the rate concurrently
		go func(rr *pb.SubscribeRateRequest) {
			// Fetch the current exchange rate for the received crypto/fiat pair
			rate, err := fetchExchangeRate(rateAPIURL, rr)
			if err != nil {
				log.Println("Error fetching exchange rate:", err)
				// Handle error, but don't block the loop
				return
			}

			// Send the rate back to the client via the gRPC stream
			log.Print("sending resp")
			err = src.Send(&pb.SubscribeRateResponse{Rate: rate})
			if err != nil {
				log.Println("Error sending response to client:", err)
				// Handle sending error, but don't block the loop
				return
			}
		}(rr)
	}
}
