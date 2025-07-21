package doodad

import (
	"design-library/position/box"
	"fmt"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
)

type Rectangle struct {
	Width  int
	Height int
}

type Doodad interface {
	Draw(screen *ebiten.Image)

	// Loads up any resources that need caching, such as images or fonts.
	Setup()

	// Releases any resources, unsubscribes from events, etc.
	Teardown() error

	Layout() *box.Box
	SetLayout(layout *box.Box)

	AddChild(doodads ...Doodad)
	Children() *Children
	SetChildren(children *Children)

	Parent() Doodad
	SetParent(parent Doodad)

	Gesturer() Gesturer
	SetGesturer(gesturer Gesturer)

	// this is what should be used to register gestures and return a list of unregister functions
	Gestures(gesturer Gesturer) []func()
	StoreUnregisterGestures(unregisters ...func()) // I don't like this
	StoreRegisterGestureFn(registerFn func())      // This is also stupid

	Hide()
	Show()
	IsVisible() bool
}

type Children struct {
	Parent  Doodad
	Doodads []Doodad
}

func NewChildren(parent Doodad, children ...[]Doodad) *Children {
	c := &Children{
		Doodads: []Doodad{},
		Parent:  parent,
	}
	for _, childGroup := range children {
		for _, doodad := range childGroup {
			c.add(doodad)
		}
	}
	return c
}

func (c *Children) Boxes() []*box.Box {
	boxes := make([]*box.Box, len(c.Doodads))
	for i, doodad := range c.Doodads {
		boxes[i] = doodad.Layout()
	}
	return boxes
}

func (c *Children) Draw(screen *ebiten.Image) {
	for _, doodad := range c.Doodads {
		doodad.Draw(screen)
	}
}

func (c *Children) Setup() {
	for i := len(c.Doodads) - 1; i >= 0; i-- {
		ii := i
		c.Doodads[ii].Setup()
		// This is so stupid
		c.Doodads[ii].StoreRegisterGestureFn(func() {
			c.Doodads[ii].StoreUnregisterGestures(c.Doodads[i].Gestures(c.Parent.Gesturer())...)
		})
		c.Doodads[ii].StoreUnregisterGestures(c.Doodads[i].Gestures(c.Parent.Gesturer())...)
	}
}

func (c *Children) add(doodad Doodad) {
	c.Doodads = append(c.Doodads, doodad)
}

func (c *Children) Teardown() error {
	for _, doodad := range c.Doodads {
		if err := doodad.Teardown(); err != nil {
			return err
		}
		doodad.Layout().ClearDependents()
	}
	return nil
}

func (c *Children) Clear() error {
	if err := c.Teardown(); err != nil {
		return fmt.Errorf("failed to teardown children: %w", err)
	}
	c.Doodads = []Doodad{}
	return nil
}

type Rectangular interface {
	Dimensions() Rectangle
}

type Default struct {
	children *Children
	Box      *box.Box

	parent   Doodad
	gesturer Gesturer

	hidden bool

	unregisterGestures []func()
	register           func()
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

func (t *Default) unregister() {
	for _, unregister := range t.unregisterGestures {
		unregister()
	}
	t.unregisterGestures = nil
}

func (t *Default) Teardown() error {
	t.unregister()
	t.Children().Clear()

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

		t.Children().add(doodad)
		doodad.SetParent(t)
		doodad.SetGesturer(t.Gesturer())
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

func (t *Default) Gesturer() Gesturer {
	return t.gesturer
}

func (t *Default) SetGesturer(gesturer Gesturer) {
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
	t.unregister()
	for _, child := range t.Children().Doodads {
		child.Hide()
	}
}

func (t *Default) Show() {
	t.hidden = false
	t.register()
	for _, child := range t.Children().Doodads {
		child.Show()
	}
}

func (t *Default) IsVisible() bool {
	return !t.hidden
}

func (t *Default) Gestures(gesturer Gesturer) []func() { return nil }

func (t *Default) StoreUnregisterGestures(unregisters ...func()) {
	t.unregisterGestures = unregisters
}

func (t *Default) StoreRegisterGestureFn(registerFn func()) {
	t.register = registerFn
}
