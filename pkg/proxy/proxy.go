package proxy

import (
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
)

func Retry(maxAttempts int, timeout time.Duration) func(endpoint.Endpoint) endpoint.Endpoint {
	return func(e endpoint.Endpoint) endpoint.Endpoint {
		endpointer := sd.FixedEndpointer{e}
		balancer := lb.NewRoundRobin(endpointer)

		return lb.Retry(maxAttempts, timeout, balancer)
	}
}
