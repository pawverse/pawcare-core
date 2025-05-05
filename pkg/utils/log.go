package utils

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/transport/grpc"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

type contextKey string

const (
	RequestIdContextKey contextKey = "request_id"
)

func RequestIdHTTPToContext() kithttp.RequestFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		requestId := uuid.New().String()

		return context.WithValue(ctx, RequestIdContextKey, requestId)
	}
}

func RequestIdGRPCToContext() grpc.ServerRequestFunc {
	return func(ctx context.Context, md metadata.MD) context.Context {
		requestId := uuid.New().String()

		return context.WithValue(ctx, RequestIdContextKey, requestId)
	}
}
