package main

import (
	design_library "design-library"
	"design-library/doodad"
	"design-library/label"
	"design-library/position/box"
	"design-library/stack"
	"image/color"
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
	nav := NewNavBar(p.App)
	p.AddChild(nav)

	contentStack := stack.New(stack.Config{
		BackgroundColor: color.RGBA{100, 200, 120, 100},
	})
	p.AddChild(contentStack)

	contentStack.AddChild(label.New(label.Config{
		Message: "This is the Second page",
	}))

	contentPane := box.Computed(func(b *box.Box) {
		b.Copy(p.Box).DecreaseWidth(nav.Box.Width()).MoveRight(nav.Box.Width())
	})

	contentStack.Layout().Computed(func(b *box.Box) {
		b.AlignTopWithin(contentPane).AlignLeftWithin(contentPane)
	})

	p.Children().Setup()

}
