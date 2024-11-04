package callback

import "log/slog"

type Registry[T any] struct {
	callbacks []func(T)
}

func NewRegistry[T any]() *Registry[T] {
	return &Registry[T]{
		callbacks: []func(T){},
	}
}

func (r *Registry[T]) Register(callback func(T)) func() {
	r.callbacks = append(r.callbacks, callback)
	return func() {
		for i, cb := range r.callbacks {
			if functionsAreTheSame(cb, callback) {
				r.callbacks = append(r.callbacks[:i], r.callbacks[i+1:]...)
				break
			}
		}
	}
}

func (r *Registry[T]) Invoke(args T) {
	slog.Debug("Invoking registry")
	for _, cb := range r.callbacks {
		slog.Debug("Invoking callback")
		cb(args)
	}
}
