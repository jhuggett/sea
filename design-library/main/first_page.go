package main

import (
	design_library "design-library"
	"design-library/button"
	"design-library/config"
	"design-library/doodad"
	"design-library/label"
	"design-library/position/box"
	"design-library/stack"
	"image/color"
	"log/slog"
)

func NewFirstPage(
	app *design_library.App,
) *firstPage {
	page := &firstPage{
		Default: doodad.Default{},
		App:     app,
	}

	return page
}

type firstPage struct {
	doodad.Default

	App *design_library.App
}

func (p *firstPage) Setup() {

	navBar := NewNavBar(p.App)
	p.AddChild(navBar)

	titleLabel := label.New(label.Config{
		Message:  "First Page",
		FontSize: 36,
		BackgroundColor: color.RGBA{
			R: 255,
			G: 100,
			A: 250,
		},
		ForegroundColor: color.RGBA{
			G: 255,
		},
		Padding: label.Padding{
			Top:    10,
			Right:  20,
			Bottom: 10,
			Left:   20,
		},
	})

	exampleButton := button.New(button.Config{
		OnClick: func(b *button.Button) {
			slog.Info("Example button clicked")
		},
		Config: label.Config{
			Message: "Click Me",
			BackgroundColor: color.RGBA{
				R: 50,
				A: 100,
			},
			Padding: label.Padding{
				Left:   10,
				Right:  10,
				Top:    10,
				Bottom: 10,
			},
		},
	})

	exampleButton2 := button.New(button.Config{
		OnClick: func(b *button.Button) {
			slog.Info("No Click Me button clicked")
		},
		Config: label.Config{
			Message: "No Click Me",
			BackgroundColor: color.RGBA{
				R: 50,
				A: 100,
			},
			Padding: label.Padding{
				Left:   10,
				Right:  10,
				Top:    10,
				Bottom: 10,
			},
		},
	})

	anotherLabel := label.New(label.Config{
		Message: "This is another label",
	})

	yetAnotherLabel := label.New(label.Config{
		Message: "This is yet another label",
	})

	mainStack := stack.New(stack.Config{
		SpaceBetween: 10,
		Padding:      config.EqualPadding(20),
		BackgroundColor: color.RGBA{
			R: 100,
			G: 150,
			B: 170,
			A: 255,
		},
	})

	p.AddChild(mainStack)

	mainStack.AddChild(
		titleLabel,
		exampleButton,
		exampleButton2,
		anotherLabel,
		yetAnotherLabel,
	)

	toggleMessage := "Hide"

	toggleButton := button.New(button.Config{
		OnClick: func(b *button.Button) {
			if toggleMessage == "Hide" {
				toggleMessage = "Show"
				mainStack.Hide()
			} else {
				toggleMessage = "Hide"
				mainStack.Show()
			}
			b.SetMessage(toggleMessage)
		},
		Config: label.Config{
			Message: toggleMessage,
			BackgroundColor: color.RGBA{
				R: 100,
				G: 150,
				B: 100,
				A: 255,
			},
			Padding: label.Padding{
				Top:    10,
				Right:  20,
				Bottom: 10,
				Left:   20,
			},
		},
	})

	p.AddChild(toggleButton)

	contentPane := box.Computed(func(b *box.Box) {
		b.Copy(p.Box).DecreaseWidth(navBar.Box.Width()).MoveRight(navBar.Box.Width())
	})

	mainStack.Layout().Computed(func(b *box.Box) {
		b.AlignTopWithin(contentPane).AlignLeftWithin(contentPane)
	})

	p.Children().Setup()

	toggleButton.Layout().Computed(func(b *box.Box) {
		b.AlignRight(p.Layout()).AlignBottom(p.Layout())
	})
}
