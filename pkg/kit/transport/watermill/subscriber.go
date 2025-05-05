package watermill

import (
	"context"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport"
	"github.com/go-kit/log"
)

type Subscriber struct {
	e            endpoint.Endpoint
	dec          DecodeRequestFunc
	enc          EncodeResponseFunc
	before       []RequestFunc
	after        []SubscriberResponseFunc
	errorEncoder ErrorEncoder
	finalizer    []SubscriberFinalizerFunc
	errorHandler transport.ErrorHandler
}

// NewSubscriber constructs a new subscriber
func NewSubscriber(
	e endpoint.Endpoint,
	dec DecodeRequestFunc,
	enc EncodeResponseFunc,
	options ...SubscriberOption,
) *Subscriber {
	s := &Subscriber{
		e:            e,
		dec:          dec,
		enc:          enc,
		errorHandler: transport.NewLogErrorHandler(log.NewNopLogger()),
	}

	for _, option := range options {
		option(s)
	}

	return s
}

func (s *Subscriber) Handle(msg *message.Message) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if len(s.finalizer) > 0 {
		defer func() {
			for _, f := range s.finalizer {
				f(ctx, msg)
			}
		}()
	}

	for _, f := range s.before {
		ctx = f(ctx, msg)
	}

	request, err := s.dec(ctx, msg)
	if err != nil {
		return err
	}

	response, err := s.e(ctx, request)
	if err != nil {
		return err
	}

	for _, f := range s.after {
		ctx = f(ctx)
	}

	return s.enc(ctx, msg, response)
}

// SubscriberOption sets an optional parameter for subscribers.
type SubscriberOption func(*Subscriber)

// ErrorEncoder is responsible for encoding an error to the subscriber reply.
// Users are encouraged to use custom ErrorEncoders to encode errors to
// their replies, and will likely want to pass and check for their own error
// types.
type ErrorEncoder func(ctx context.Context, err error, reply string)

// SubscriberFinalizerFunc can be used to perform work at the end of an request
// from a publisher, after the response has been written to the publisher. The principal
// intended use is for request logging.
type SubscriberFinalizerFunc func(ctx context.Context, msg *message.Message)

func EncodeResponse(ctx context.Context, _ *message.Message, response any) error {
	if e, ok := response.(endpoint.Failer); ok && e.Failed() != nil {
		return e.Failed()
	}
	return nil
}

func DecodeJSONMessage[T any](_ context.Context, msg *message.Message) (any, error) {
	var request T
	if err := json.Unmarshal(msg.Payload, &request); err != nil {
		return nil, err
	}
	return request, nil
}
