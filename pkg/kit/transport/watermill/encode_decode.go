package watermill

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
)

type DecodeRequestFunc func(context.Context, *message.Message) (request any, err error)

type EncodeRequestFunc func(context.Context, *message.Messages, any) error

type EncodeResponseFunc func(context.Context, *message.Message, any) error

type DecodeResponseFunc func(context.Context, *message.Message) (response any, err error)
