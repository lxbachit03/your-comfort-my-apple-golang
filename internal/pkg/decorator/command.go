package decorator

import "context"

type CommandHandler[C any, R any] interface {
	Handle(ctx context.Context, cmd C) (R, error)
}

func ApplyCommandDecorators[C any, R any](handler CommandHandler[C, R]) CommandHandler[C, R] {
	return handler
}
