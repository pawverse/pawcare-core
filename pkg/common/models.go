package common

import "github.com/go-kit/kit/endpoint"

type EmbedError struct {
	Err error `json:"-"`
}

func NewEmbededError(err error) EmbedError {
	return EmbedError{Err: err}
}

func (b EmbedError) Failed() error {
	if b.Err == nil {
		return nil
	}

	return NewErrorResponse(b.Err)
}

var _ endpoint.Failer = (*EmbedError)(nil)
