package callback

import (
	"errors"

	"github.com/google/uuid"
)

type CallbackRegistryRecord[T any] struct {
	Callback T
	UUID     string
}

type CallbackRegistry[T any] struct {
	callbacks []CallbackRegistryRecord[T]
}

func (r *CallbackRegistry[T]) Add(callback T) string {
	uuid := uuid.NewString()
	r.callbacks = append(r.callbacks, CallbackRegistryRecord[T]{Callback: callback, UUID: uuid})
	return uuid
}

func (r *CallbackRegistry[T]) Remove(uuid string) {
	for i, cb := range r.callbacks {
		if cb.UUID == uuid {
			r.callbacks = append(r.callbacks[:i], r.callbacks[i+1:]...)
			break
		}
	}
}

var ErrStopPropagation = errors.New("stop propagation")
var ErrUnregister = errors.New("unregister")

func (r *CallbackRegistry[T]) InvokeEndToStart(do func(T) error) {
	for i := len(r.callbacks) - 1; i >= 0; i-- {
		cb := r.callbacks[i]
		err := do(cb.Callback)
		if err != nil {
			if errors.Is(err, ErrStopPropagation) {
				return
			}
			if errors.Is(err, ErrUnregister) {
				r.Remove(cb.UUID)
			}
		}
	}
}

func (r *CallbackRegistry[T]) Register(callback T) func() {
	uuid := r.Add(callback)
	return func() {
		r.Remove(uuid)
	}
}
