package gapi

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Logger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	startTime := time.Now()
	result, err := handler(ctx, req)
	duration := time.Since(startTime)

	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}

	logger := log.Info()
	if err != nil {
		logger = log.Error().Err(err)
	}

	logger.Str("protocol", "grpc").Str("method", info.FullMethod).Int("status_code", int(statusCode)).Dur("duration", duration).Msg("recieved request")
	return result, err
}

type ResponseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

func (rr *ResponseRecorder) WriteHeader(statusCode int) {
	rr.statusCode = statusCode
	rr.ResponseWriter.WriteHeader(statusCode)
}

func (rr *ResponseRecorder) Write(body []byte) (int, error) {
	rr.body = body
	return rr.ResponseWriter.Write(body)
}

func HttpLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		startTime := time.Now()
		wr := &ResponseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		handler.ServeHTTP(wr, r)
		duration := time.Since(startTime)

		logger := log.Info()
		if wr.statusCode != http.StatusOK {
			logger = log.Error().Bytes("body", wr.body)
		}

		logger.
			Str("protocol", "http").
			Str("method", r.Method).
			Str("path", r.RequestURI).
			Int("status_code", wr.statusCode).
			Dur("duration", duration).
			Msg("received request")
	})
}
