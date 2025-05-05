package common

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/pawverse/pawcare-core/pkg/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func DecodeNoBodyRequest(_ context.Context, r *http.Request) (any, error) {
	return nil, nil
}

func DecodePathParameters[T any](_ context.Context, r *http.Request) (any, error) {
	vars := mux.Vars(r)
	result, err := utils.MapToStruct[T](vars)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func DecodeJSONRequest[T any](_ context.Context, r *http.Request) (any, error) {
	var request T
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func KafkaDecodeJSONMessage[T any](_ context.Context, msg *message.Message) (any, error) {
	var request T
	if err := json.Unmarshal(msg.Payload, &request); err != nil {
		return nil, err
	}
	return request, nil
}

func KafkaEncodeResponse(ctx context.Context, response any) error {
	if e, ok := response.(endpoint.Failer); ok && e.Failed() != nil {
		return e.Failed()
	}
	return nil
}

func GRPCDecodeNoBody(_ context.Context, req any) (any, error) {
	return nil, nil
}

func GRPCDecodeToObject[T any](_ context.Context, request any) (any, error) {
	req, ok := request.(T)
	if !ok {
		return nil, ErrCastRequest
	}

	return req, nil
}

// EncodeJSONResponse is a EncodeResponseFunc that serializes the response as a
// JSON object to the ResponseWriter. Many JSON-over-HTTP services can use it as
// a sensible default. If the response implements Headerer, the provided headers
// will be applied to the response. If the response implements StatusCoder, the
// provided StatusCode will be used instead of 200.
func EncodeJSONResponse(ctx context.Context, w http.ResponseWriter, response any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if headerer, ok := response.(kithttp.Headerer); ok {
		for k, values := range headerer.Headers() {
			for _, v := range values {
				w.Header().Add(k, v)
			}
		}
	}
	code := http.StatusOK
	if sc, ok := response.(kithttp.StatusCoder); ok {
		code = sc.StatusCode()
	}
	w.WriteHeader(code)
	if code == http.StatusNoContent {
		return nil
	}

	if failer, ok := response.(endpoint.Failer); ok {
		if failer.Failed() != nil {
			kithttp.DefaultErrorEncoder(ctx, failer.Failed(), w)
			return nil
		}
	}

	return json.NewEncoder(w).Encode(response)
}

func ErrorEncoder(err2Status func(error) int) func(_ context.Context, err error, w http.ResponseWriter) {
	return func(_ context.Context, err error, w http.ResponseWriter) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(err2Status(err))
		json.NewEncoder(w).Encode(NewErrorResponse(err))
	}
}

func HTTPLoggingServerOptions(logger *zap.Logger) []kithttp.ServerOption {
	return []kithttp.ServerOption{
		kithttp.ServerBefore(func(ctx context.Context, r *http.Request) context.Context {
			clientIp := r.RemoteAddr
			method := r.Method
			path := r.URL.Path
			contentLength := r.ContentLength
			userAgent := r.UserAgent()
			requestId := ctx.Value(utils.RequestIdContextKey).(string)

			logger.
				Info("Incoming request",
					zap.String("request_id", requestId),
					zap.String("method", method),
					zap.String("path", path),
					zap.String("client_ip", clientIp),
					zap.String("user_agent", userAgent),
					zap.Int64("content_length", contentLength),
				)

			return ctx
		}),

		kithttp.ServerFinalizer(func(ctx context.Context, code int, r *http.Request) {
			method := r.Method
			path := r.URL.Path
			clientIp := r.RemoteAddr
			requestId := ctx.Value(utils.RequestIdContextKey).(string)
			userAgent := r.UserAgent()

			logger.
				Info("Outgoing response",
					zap.String("request_id", requestId),
					zap.String("method", method),
					zap.String("path", path),
					zap.String("client_ip", clientIp),
					zap.String("user_agent", userAgent),
					zap.Int("status_code", code),
				)
		}),
	}
}

func GRPCLoggingServerOptions(logger *zap.Logger) []kitgrpc.ServerOption {
	return []kitgrpc.ServerOption{
		kitgrpc.ServerBefore(func(ctx context.Context, md metadata.MD) context.Context {
			requestId := ctx.Value(utils.RequestIdContextKey).(string)
			logs := []zapcore.Field{zap.String("request_id", requestId)}
			if userAgents, ok := md["user-agent"]; ok {
				logs = append(logs, zap.String("user_agent", strings.Join(userAgents, ";")))
			}

			if p, ok := peer.FromContext(ctx); ok && p != nil {
				logs = append(logs, zap.Stringer("client_ip", p.Addr))
			}

			logger.Info("Incoming request", logs...)

			return ctx
		}),
		kitgrpc.ServerFinalizer(func(ctx context.Context, err error) {
			requestId := ctx.Value(utils.RequestIdContextKey).(string)
			fields := []zapcore.Field{zap.String("request_id", requestId)}

			if md, ok := metadata.FromIncomingContext(ctx); ok {
				if userAgents, ok := md["user-agent"]; ok {
					fields = append(fields, zap.String("user_agent", strings.Join(userAgents, ";")))
				}
			}

			if p, ok := peer.FromContext(ctx); ok && p != nil {
				fields = append(fields, zap.Stringer("client_ip", p.Addr))
			}

			logger.Info("Outgoing response", fields...)
		}),
	}
}

type logErrorHandler struct {
	logger *zap.Logger
}

func NewLogErrorHandler(logger *zap.Logger) transport.ErrorHandler {
	return &logErrorHandler{logger}
}

func (h *logErrorHandler) Handle(ctx context.Context, err error) {
	h.logger.Error(err.Error())
}
