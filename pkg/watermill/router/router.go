package router

import (
	"github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/pawverse/pawcare-core/pkg/watermill/log"
	"go.uber.org/zap"
)

func NewDefaultRouter(logger *zap.Logger) (*message.Router, error) {
	router, err := message.NewRouter(message.RouterConfig{}, log.NewLogger(logger))
	if err != nil {
		return nil, err
	}

	router.AddPlugin(plugin.SignalsHandler)
	router.AddMiddleware(
		middleware.CorrelationID,
		middleware.Recoverer,
		middleware.RandomFail(0.1),
		middleware.RandomPanic(0.1),
	)

	return router, nil
}

func generatePartitionKey(topic string, msg *message.Message) (string, error) {
	if key, ok := msg.Metadata["key"]; ok {
		return key, nil
	}

	if key, ok := msg.Context().Value("key").(string); ok {
		return key, nil
	}

	return "", nil
}

func NewDefaultPartitionKeyMarshaler() kafka.Marshaler {
	return kafka.NewWithPartitioningMarshaler(generatePartitionKey)
}
