package main

import (
	design_library "design-library"
	"design-library/button"
	"design-library/doodad"
	"design-library/label"
	"design-library/position/box"
	"design-library/stack"
	"image/color"
)

func NewNavBar(
	app *design_library.App,
) *NavBar {
	navBar := &NavBar{
		App: app,
	}

	return navBar
}

type NavBar struct {
	doodad.Default

	App *design_library.App
}

func (n *NavBar) Setup() {
	titleLabel := label.New(label.Config{
		Message:  "Design Library",
		FontSize: 24,
		ForegroundColor: color.RGBA{
			R: 0,
			G: 0,
			B: 0,
			A: 255,
		},
	})

	firstPageButton := button.New(button.Config{
		OnClick: func(b *button.Button) {
			n.App.Replace(NewFirstPage(n.App))
		},
		Config: label.Config{
			Message: "First Page",
		},
	})

	secondPageButton := button.New(button.Config{
		OnClick: func(b *button.Button) {
			n.App.Replace(NewSecondPage(n.App))
		},
		Config: label.Config{
			Message: "Second Page",
		},
	})

	thirdPageButton := button.New(button.Config{
		OnClick: func(b *button.Button) {
			n.App.Replace(NewThirdPage(n.App))
		},
		Config: label.Config{
			Message: "Third Page",
		},
	})

	mainStack := stack.New(stack.Config{
		Type: stack.Vertical,
		BackgroundColor: color.RGBA{
			R: 100,
			G: 150,
			B: 100,
			A: 255,
		},
		Layout: box.Computed(func(b *box.Box) {
			b.SetDimensions(titleLabel.Box.Width(), n.Layout().Height())
		}),
		Padding: stack.Padding{
			Top:    10,
			Right:  20,
			Bottom: 10,
			Left:   20,
		},
		SpaceBetween: 10,
	})
	n.AddChild(mainStack)

	mainStack.AddChild(
		titleLabel,
		firstPageButton,
		secondPageButton,
		thirdPageButton,
	)

	n.Children().Setup()
}
