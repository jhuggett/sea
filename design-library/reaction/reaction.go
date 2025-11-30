package reaction

type Reaction interface {
	ReactionType() ReactionType

	SetUnregister(func())
	Unregister()

	SetEnabled(enabled bool)
	IsEnabled() bool

	TryPerform(event *Event, data any) error

	SetDepth(depth int)
	Depth() int
}

func NewReaction[T any](
	reactionType ReactionType,
	condition func(T) bool,
	callback func(T),
) Reaction {
	return &basicReaction[T]{
		Type:      reactionType,
		Condition: condition,
		Callback:  callback,
		Enabled:   true,
	}
}
