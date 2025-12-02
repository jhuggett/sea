package reaction

type basicReaction[T any] struct {
	Type      ReactionType
	Condition func(T) bool
	Callback  func(T)
	Enabled   bool

	unregister func()

	depth []int
}

func (r *basicReaction[T]) IsEnabled() bool {
	return r.Enabled
}

func (r *basicReaction[T]) SetDepth(depth []int) {
	r.depth = depth
}

func (r *basicReaction[T]) Depth() []int {
	return r.depth
}

func (r *basicReaction[T]) MeetsCondition(t T) bool {
	if r.Condition == nil {
		return true
	}
	return r.Condition(t)
}

func (r *basicReaction[T]) PerformCallback(t T) error {
	if r.Callback == nil {
		return nil
	}
	r.Callback(t)
	return nil
}

func (r *basicReaction[T]) TryPerform(event *Event, data any) error {
	if !r.Enabled {
		return nil
	}

	if t, ok := data.(T); ok {
		if !r.MeetsCondition(t) {
			return nil
		}
		return r.PerformCallback(t)
	}
	return nil
}

func (r *basicReaction[T]) SetUnregister(unregister func()) {
	r.unregister = unregister
}

func (r *basicReaction[T]) Unregister() {
	if r.unregister != nil {

		r.SetEnabled(false)

		r.unregister()
		r.unregister = nil
	}
}

func (r *basicReaction[T]) ReactionType() ReactionType {
	return r.Type
}

func (r *basicReaction[T]) SetEnabled(enabled bool) {
	r.Enabled = enabled
}
