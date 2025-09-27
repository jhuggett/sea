package reaction

type Reaction interface {
	ReactionType() ReactionType

	SetUnregister(func())
	Unregister()

	SetEnabled(enabled bool)

	TryPerform(event *Event, data any) error

	SetDepth(depth int)
	Depth() int
}

type basicEvent[T any] struct {
}

type basicReaction[T any] struct {
	Type      ReactionType
	Condition func(T) bool
	Callback  func(T)
	Enabled   bool

	unregister func()

	depth int
}

func (r *basicReaction[T]) SetDepth(depth int) {
	r.depth = depth
}

func (r *basicReaction[T]) Depth() int {
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

type Reactions struct {
	Reactions []Reaction
}

func (r *Reactions) Enable() {
	for _, reaction := range r.Reactions {
		reaction.SetEnabled(true)
	}
}

func (r *Reactions) Disable() {
	for _, reaction := range r.Reactions {
		reaction.SetEnabled(false)
	}
}

func (r *Reactions) Register(gesturer Gesturer, atDepth int) {
	for _, r := range r.Reactions {
		r.SetUnregister(gesturer.Register(r, atDepth))
	}
}

func (r *Reactions) Add(reactions ...Reaction) {
	r.Reactions = append(r.Reactions, reactions...)
}

func (r *Reactions) Unregister() {
	for _, reaction := range r.Reactions {
		if reaction != nil {
			reaction.Unregister()
		}
	}

	r.Reactions = nil
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

// func NewMouseUpReaction(doodad doodad.Doodad, do func(MouseUpEvent)) *basicReaction[MouseUpEvent] {
// 	return &basicReaction[MouseUpEvent]{
// 		Type:    MouseUp,
// 		Enabled: true,
// 		Condition: func(mue MouseUpEvent) bool {
// 			withinBounds := func(x, y int) bool {
// 				return x >= doodad.Layout().X() && x <= doodad.Layout().X()+doodad.Layout().Width() &&
// 					y >= doodad.Layout().Y() && y <= doodad.Layout().Y()+doodad.Layout().Height()
// 			}

// 			if !withinBounds(mue.X, mue.Y) {
// 				return false
// 			}
// 			return true
// 		},
// 		Callback: func(event MouseUpEvent) error {

// 			do(event)
// 			return nil
// 		},
// 	}
// }

// func NewMouseMovedWithinReaction(doodad doodad.Doodad, do func(MouseMoved)) *basicReaction[MouseMoved] {
// 	return &basicReaction[MouseMoved]{
// 		Type:    MouseMove,
// 		Enabled: true,
// 		Condition: func(mmf MouseMoved) bool {
// 			withinBounds := func(x, y int) bool {
// 				return x >= doodad.Layout().X() && x <= doodad.Layout().X()+doodad.Layout().Width() &&
// 					y >= doodad.Layout().Y() && y <= doodad.Layout().Y()+doodad.Layout().Height()
// 			}

// 			if !withinBounds(mmf.X, mmf.Y) {
// 				return false
// 			}
// 			return true
// 		},
// 		Callback: func(event MouseMoved) error {
// 			do(event)
// 			return nil
// 		},
// 	}
// }

// func NewMouseMovedWithoutReaction(doodad doodad.Doodad, do func(MouseMoved)) *basicReaction[MouseMoved] {
// 	return &basicReaction[MouseMoved]{
// 		Type:    MouseMove,
// 		Enabled: true,
// 		Condition: func(mmf MouseMoved) bool {
// 			withinBounds := func(x, y int) bool {
// 				return x >= doodad.Layout().X() && x <= doodad.Layout().X()+doodad.Layout().Width() &&
// 					y >= doodad.Layout().Y() && y <= doodad.Layout().Y()+doodad.Layout().Height()
// 			}

// 			if withinBounds(mmf.X, mmf.Y) {
// 				return false
// 			}
// 			return true
// 		},
// 		Callback: func(event MouseMoved) error {
// 			do(event)
// 			return nil
// 		},
// 	}
// }
