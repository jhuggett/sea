package doodad

import (
	"design-library/position/box"
	"fmt"

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

	AddChild(doodad Doodad)
}

type Children struct {
	Doodads []Doodad
}

func NewChildren(children ...[]Doodad) *Children {
	c := &Children{
		Doodads: []Doodad{},
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
	for _, doodad := range c.Doodads {
		doodad.Setup()
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
	Children *Children
	Gesturer Gesturer
	Box      *box.Box
}

func (t *Default) Update() error {
	return nil
}

func (t *Default) Draw(screen *ebiten.Image) {
	t.Children.Draw(screen)
}

func (t *Default) Teardown() error {
	// Teardown logic for the tabs can go here
	return nil
}

func (t *Default) Layout() *box.Box {
	return t.Box
}

func (t *Default) AddChild(doodad Doodad) {
	if t.Children == nil {
		t.Children = &Children{}
	}
	t.Children.add(doodad)
	t.Layout().AddDependent(doodad.Layout())
}

func NewDefault() *Default {
	return &Default{
		Children: NewChildren(),
		Gesturer: NewGesturer(),
		Box:      box.Zeroed(),
	}
}
