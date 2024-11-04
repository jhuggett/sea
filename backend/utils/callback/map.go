package callback

import (
	"fmt"
	"log/slog"
)

type RegistryMap[T any] struct {
	registries map[string]*Registry[T]
}

func NewRegistryMap[T any]() *RegistryMap[T] {
	return &RegistryMap[T]{
		registries: map[string]*Registry[T]{},
	}
}

func createKey(id []any) string {
	key := ""
	for _, i := range id {
		key += fmt.Sprintf("%v", i)
	}
	return key
}

func (r *RegistryMap[T]) Register(id []any, callback func(T)) func() {
	key := createKey(id)

	slog.Debug("Registering map callback", "key", key)

	if r.registries[key] == nil {
		slog.Debug("Creating new registry", "key", key)
		r.registries[key] = NewRegistry[T]()
	}

	r.registries[key].Register(callback)

	return func() {
		r.registries[key].Register(callback)

		if len(r.registries[key].callbacks) == 0 {
			delete(r.registries, key)
		}
	}
}

func (r *RegistryMap[T]) Invoke(id []any, args T) {
	key := createKey(id)

	slog.Debug("Invoking registry map for", "key", key)

	if r.registries[key] == nil {
		return
	}
	r.registries[key].Invoke(args)
}
