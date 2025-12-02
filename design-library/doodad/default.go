package doodad

import (
	"design-library/position/box"
	"design-library/reaction"

	"github.com/hajimehoshi/ebiten/v2"
)

type Default struct {
	children *Children
	Box      *box.Box

	parent    Doodad
	gesturer  reaction.Gesturer
	reactions *reaction.Reactions

	hidden bool

	z []int

	cachedDraw []*CachedDraw

	actionOnTeardown []func()

	statefulDoodads map[string]Doodad
}

func (t *Default) DoOnTeardown(actions ...func()) {
	t.actionOnTeardown = append(t.actionOnTeardown, actions...)
}

func (t *Default) Z() []int {
	return t.z
}

func (t *Default) SetZ(z []int) {
	t.z = z
}

func (t *Default) Update() error {
	return nil
}

func (t *Default) Draw(screen *ebiten.Image) {
	if t.hidden {
		return
	}

	if cached := t.CachedDraw(); cached != nil {
		for _, draw := range cached {
			op := &ebiten.DrawImageOptions{}
			if draw.Op != nil {
				op = draw.Op
			}
			op.GeoM.Translate(float64(draw.X), float64(draw.Y))

			x, y := t.Layout().XY()
			op.GeoM.Translate(float64(x), float64(y))

			screen.DrawImage(draw.Image, op)
		}
	}

	// t.Children().Draw(screen)
}

func (t *Default) Teardown() error {
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
		// parentType := fmt.Sprintf("%T", t)
		// childType := fmt.Sprintf("%T", doodad)
		// fmt.Printf("🔗 Adding child: \x1b[36m%s\x1b[0m to parent: \x1b[32m%s\x1b[0m\n", childType, parentType)

		// // Still log structured information
		// slog.Info(
		// 	"Adding child to parent",
		// 	"parent", parentType,
		// 	"child", childType,
		// )

		if doodad.Children() == nil {
			doodad.SetChildren(NewChildren(doodad))
		}

		if doodad.Layout() == nil {
			doodad.SetLayout(box.Computed(func(b *box.Box) {
				b.Copy(t.Layout())
			}))
		}
		if !t.Layout().HasDependent(doodad.Layout()) {
			t.Layout().AddDependent(doodad.Layout())
		}

		if doodad.Reactions() == nil {
			doodad.SetReactions(&reaction.Reactions{})
		}

		parentZ := t.Z()

		if len(parentZ) == 0 {
			parentZ = []int{0}
		}

		childZ := make([]int, len(parentZ))
		copy(childZ, parentZ)
		childZ[len(childZ)-1] += 1

		doodad.SetZ(childZ)

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

func (t *Default) CachedDraw() []*CachedDraw {
	return t.cachedDraw
}

func (t *Default) SetCachedDraw(cachedDraw ...*CachedDraw) {
	t.cachedDraw = cachedDraw
}

func (t *Default) StatefulDoodads() map[string]Doodad {
	if t.statefulDoodads == nil {
		t.statefulDoodads = make(map[string]Doodad)
	}
	return t.statefulDoodads
}

func (t *Default) AddStatefulChild(key string, create func() Doodad) Doodad {
	child := Stateful(t, key, create)

	t.AddChild(child)

	return child
}

func Stateful[D Doodad](d Doodad, key string, create func() D) D {
	statefulMap := d.StatefulDoodads()

	if existing, exists := statefulMap[key]; exists {
		return existing.(D)
	}

	newDoodad := create()
	statefulMap[key] = newDoodad
	return newDoodad
}
