package doodad

import (
	"design-library/position/box"
	"design-library/reaction"

	"github.com/hajimehoshi/ebiten/v2"
)

func MouseIsWithin[T reaction.PositionedEvent](doodad Doodad) func(event T) bool {
	return func(event T) bool {
		if !doodad.IsVisible() {
			return false
		}
		layout := doodad.Layout()
		x, y := event.XY()
		return x >= layout.X() && x <= layout.X()+layout.Width() &&
			y >= layout.Y() && y <= layout.Y()+layout.Height()
	}
}

func MouseIsOutside[T reaction.PositionedEvent](doodad Doodad) func(event T) bool {
	return func(event T) bool {
		if !doodad.IsVisible() {
			return false
		}

		layout := doodad.Layout()
		x, y := event.XY()

		return x < layout.X() || x > layout.X()+layout.Width() ||
			y < layout.Y() || y > layout.Y()+layout.Height()
	}
}

type Rectangle struct {
	Width  int
	Height int
}

type CachedDraw struct {
	Image *ebiten.Image
	X     int
	Y     int
	Op    *ebiten.DrawImageOptions
}

type Doodad interface {
	Draw(screen *ebiten.Image)

	CachedDraw() []*CachedDraw
	SetCachedDraw(cachedDraw ...*CachedDraw)

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

	Gesturer() reaction.Gesturer
	SetGesturer(gesturer reaction.Gesturer)

	Reactions() *reaction.Reactions
	SetReactions(*reaction.Reactions)

	// // this is what should be used to register gestures and return a list of unregister functions
	// Gestures(gesturer Gesturer) []func()
	// StoreUnregisterGestures(unregisters ...func()) // I don't like this
	// StoreRegisterGestureFn(registerFn func())      // This is also stupid

	Hide()
	Show()
	IsVisible() bool

	Z() []int
	SetZ(z []int)

	StatefulDoodads() map[string]Doodad
	AddStatefulChild(key string, create func() Doodad) Doodad
}

type Rectangular interface {
	Dimensions() Rectangle
}

func ReSetup(doodad Doodad) {
	doodad.Teardown()
	Setup(doodad)
}

func Setup(doodads ...Doodad) {
	for _, doodad := range doodads {
		doodad.Setup()
		doodad.Layout().Recalculate()
		doodad.Reactions().Register(doodad.Gesturer(), doodad.Z())
	}
}
