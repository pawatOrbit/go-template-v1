package httpserver

import "github.com/pawatOrbit/ai-mock-data-service/go/core/transport"


type Endpoint[T, R any] func() (fn transport.Service[T, R])

func NewEndpoint[T, R any](svc transport.Service[T, R], middlewares ...transport.EndpointMiddleware[T, R]) func() Endpoint[T, R] {
	return func() Endpoint[T, R] {
		return func() transport.Service[T, R] {
			var newSvc transport.Service[T, R]
			init := true
			for _, m := range middlewares {
				if init {
					newSvc = m(svc)
					init = false
				} else {
					newSvc = m(newSvc)
				}
			}

			if len(middlewares) == 0 {
				return svc
			}

			return newSvc
		}
	}
}
