package main

import (
	design_library "design-library"
	"design-library/doodad"
)

func NewSecondPage(
	app *design_library.App,
) *SecondPage {
	page := &SecondPage{
		Default: doodad.Default{},
		App:     app,
	}

	return page
}

type SecondPage struct {
	doodad.Default

	App *design_library.App
}

func (p *SecondPage) Setup() {
	nav := NewNavBar(
		p.App,
	)
	p.AddChild(nav)

	p.Children().Setup()
}
