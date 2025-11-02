package doodad

import (
	"design-library/position/box"
	"fmt"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
)

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
	if c.Parent == nil {
		slog.Warn("Cannot setup children: parent is nil")
		return
	}

	for i := len(c.Doodads) - 1; i >= 0; i-- {
		d := c.Doodads[i]

		d.Setup()
		d.Reactions().Register(c.Parent.Gesturer(), d.Z())
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
		doodad.Reactions().Unregister()
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

func (c *Children) Remove(doodad Doodad) error {
	for i, d := range c.Doodads {
		if d == doodad {
			if err := d.Teardown(); err != nil {
				return fmt.Errorf("failed to teardown doodad: %w", err)
			}
			c.Doodads = append(c.Doodads[:i], c.Doodads[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("doodad not found in children")
}
