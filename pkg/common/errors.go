package common

import (
	"encoding/json"
	"errors"
)

var (
	ErrCastRequest   = errors.New("cannot cast request")
	ErrCastResponse  = errors.New("cannot cast response")
	ErrParsingClaims = errors.New("error parsing claims")
	ErrUnauthorized  = errors.New("unauthorized")
)

type ErrorResponse struct {
	err error
}

func NewErrorResponse(err error) *ErrorResponse {
	return &ErrorResponse{
		err: err,
	}
}

func (e *ErrorResponse) MarshalJSON() ([]byte, error) {
	obj := struct {
		Error string `json:"error"`
	}{
		Error: e.Error(),
	}

	return json.Marshal(obj)
}

func (e *ErrorResponse) Unwrap() error {
	return e.err
}

func (e *ErrorResponse) Error() string {
	return e.err.Error()
}

var (
	_ error          = (*ErrorResponse)(nil)
	_ json.Marshaler = (*ErrorResponse)(nil)
)

func Err2Str(err error) string {
	if err == nil {
		return ""
	}

	return err.Error()
}
