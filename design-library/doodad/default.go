package doodad

import (
	"design-library/position/box"
	"design-library/reaction"
	"fmt"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
)

type Default struct {
	children *Children
	Box      *box.Box

	parent    Doodad
	gesturer  reaction.Gesturer
	reactions *reaction.Reactions

	hidden bool

	z int

	actionOnTeardown []func()

	// unregisterGestures []func()
	// register           func()
}

func (t *Default) DoOnTeardown(actions ...func()) {
	t.actionOnTeardown = append(t.actionOnTeardown, actions...)
}

func (t *Default) Z() int {
	return t.z
}

func (t *Default) SetZ(z int) {
	t.z = z
}

func (t *Default) Update() error {
	return nil
}

func (t *Default) Draw(screen *ebiten.Image) {
	if t.hidden {
		return
	}
	t.Children().Draw(screen)
}

// func (t *Default) unregister() {
// 	for _, unregister := range t.unregisterGestures {
// 		unregister()
// 	}
// 	t.unregisterGestures = nil
// }

func (t *Default) Teardown() error {
	// t.unregister()

	for _, action := range t.actionOnTeardown {
		action()
	}

	t.Children().Clear()
	t.Reactions().Unregister()
	return nil
}

func (t *Default) Layout() *box.Box {
	return t.Box
}

func (t *Default) SetLayout(layout *box.Box) {
	t.Box = layout
}

func (t *Default) AddChild(doodads ...Doodad) {
	if t.Children() == nil {
		t.SetChildren(NewChildren(t))
	}

	for _, doodad := range doodads {
		// Colorized logging using ASCII color codes
		// Print colorized message about child addition
		parentType := fmt.Sprintf("%T", t)
		childType := fmt.Sprintf("%T", doodad)
		fmt.Printf("🔗 Adding child: \x1b[36m%s\x1b[0m to parent: \x1b[32m%s\x1b[0m\n", childType, parentType)

		// Still log structured information
		slog.Info(
			"Adding child to parent",
			"parent", parentType,
			"child", childType,
		)

		if doodad.Children() == nil {
			doodad.SetChildren(NewChildren(doodad))
		}

		if doodad.Layout() == nil {
			doodad.SetLayout(box.Computed(func(b *box.Box) {
				b.Copy(t.Layout())
			}))
		}
		t.Layout().AddDependent(doodad.Layout())

		if doodad.Reactions() == nil {
			doodad.SetReactions(&reaction.Reactions{})
		}

		doodad.SetZ(t.Z() + 1)

		t.Children().add(doodad)
		doodad.SetParent(t)
		doodad.SetGesturer(t.Gesturer())

		if !t.IsVisible() {
			doodad.Hide()
		}
	}
}

func NewDefault(parent Doodad) *Default {
	return &Default{
		Box:      box.Zeroed(),
		children: NewChildren(parent),
	}
}

func (t *Default) Setup() {}

func (t *Default) Parent() Doodad {
	return t.parent
}

func (t *Default) SetParent(parent Doodad) {
	t.parent = parent
}

func (t *Default) Gesturer() reaction.Gesturer {
	return t.gesturer
}

func (t *Default) SetGesturer(gesturer reaction.Gesturer) {
	t.gesturer = gesturer
}

func (t *Default) Children() *Children {
	return t.children
}

func (t *Default) SetChildren(children *Children) {
	t.children = children
}

func (t *Default) Hide() {
	t.hidden = true
	t.reactions.Disable()
	// t.unregister()
	for _, child := range t.Children().Doodads {
		child.Hide()
	}
}

func (t *Default) Show() {
	t.hidden = false
	t.reactions.Enable()
	// t.register()
	for _, child := range t.Children().Doodads {
		child.Show()
	}
}

func (t *Default) IsVisible() bool {
	return !t.hidden
}

// func (t *Default) Gestures(gesturer Gesturer) []func() { return nil }

// func (t *Default) StoreUnregisterGestures(unregisters ...func()) {
// 	t.unregisterGestures = unregisters
// }

// func (t *Default) StoreRegisterGestureFn(registerFn func()) {
// 	t.register = registerFn
// }

func (t *Default) Reactions() *reaction.Reactions {
	return t.reactions
}

func (t *Default) SetReactions(reactions *reaction.Reactions) {
	if t.reactions != nil {
		t.reactions.Unregister()
	}
	t.reactions = reactions
}
