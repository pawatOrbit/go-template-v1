package transport

import "context"

type Service[T, R any] func(ctx context.Context, req T) (resp R, err error)
