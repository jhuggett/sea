package main

import (
	design_library "design-library"
	"design-library/doodad"
)

func NewThirdPage(
	app *design_library.App,
) *ThirdPage {
	page := &ThirdPage{
		Default: doodad.Default{},
		App:     app,
	}

	return page
}

type ThirdPage struct {
	doodad.Default

	App *design_library.App
}

func (p *ThirdPage) Setup() {
	nav := NewNavBar(p.App)
	p.AddChild(nav)

	p.Children().Setup()
}
