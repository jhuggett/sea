package main

import (
	design_library "design-library"
	"design-library/button"
	"design-library/doodad"
	"design-library/label"
	"design-library/position/box"
	"design-library/stack"
	"image/color"
	"log/slog"
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

	contentPane := box.Computed(func(b *box.Box) {
		b.Copy(p.Box).DecreaseWidth(195).MoveRight(195)
	})

	contentStack := stack.New(stack.Config{
		Layout:          contentPane,
		Type:            stack.Vertical,
		BackgroundColor: color.RGBA{R: 75, G: 50, B: 75, A: 255},
	})
	p.AddChild(contentStack)

	contentStack.AddChild(button.New(button.Config{
		OnClick: func(b *button.Button) {

		},
		Config: label.Config{
			Message: "This is the second page",
		},
	}))

	container := &Container{
		Default: *doodad.NewDefault(p),
	}

	container.Layout().Computed(func(b *box.Box) {
		b.Copy(contentStack.Layout())
	})

	contentStack.AddChild(container)

	container.Contents = []doodad.Doodad{
		button.New(button.Config{
			OnClick: func(b *button.Button) {
				container.SetContents([]doodad.Doodad{
					button.New(button.Config{
						OnClick: func(b *button.Button) {
							slog.Info("Going to main page")
						},
						Config: label.Config{
							Message: "Go to the main page",
						},
					}),
				})
			},
			Config: label.Config{
				Message: "Second Page",
			},
		}),
	}

	p.Children().Setup()
}

type Container struct {
	doodad.Default

	Contents []doodad.Doodad
}

func (c *Container) Setup() {
	for _, content := range c.Contents {
		c.AddChild(content)
	}
	c.Children().Setup()
}

func (c *Container) SetContents(contents []doodad.Doodad) {
	c.Children().Clear()

	// Update the Contents field after clearing children
	c.Contents = contents

	c.Setup()

	c.Box.Recalculate()
}
