package decorator

import "context"

type QueryHandler[Q any, R any] interface {
	Handle(ctx context.Context, q Q) (R, error)
}

func ApplyQueryDecorators[Q any, R any](handler QueryHandler[Q, R]) QueryHandler[Q, R] {
	return handler
}
