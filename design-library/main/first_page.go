package main

import (
	"design-library/button"
	"design-library/doodad"
	"design-library/label"
	"design-library/position/box"
	"design-library/stack"
	"image/color"
	"log/slog"
)

func NewFirstPage() (*firstPage, error) {
	page := &firstPage{
		Default: doodad.Default{
			Gesturer: doodad.NewGesturer(),
			Box:      box.New(box.Config{}),
			Children: &doodad.Children{},
		},
	}

	page.Setup()

	return page, nil
}

type firstPage struct {
	doodad.Default
}

func (p *firstPage) Update() error {
	p.Gesturer.Update()
	return nil
}

func (p *firstPage) SetWidthAndHeight(width, height int) {
	p.Box.SetDimensions(width, height)
	p.Box.Recalculate()
}

func (p *firstPage) Setup() {
	titleLabel := label.New(label.Config{
		Message: "First Page",
		Layout:  box.Zeroed(),

		FontSize: 36,
		BackgroundColor: color.RGBA{
			R: 255,
			G: 100,
			A: 20,
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
		OnClick: func() {
			slog.Info("Example button clicked")
		},
		Gesturer: p.Gesturer,
		Config: label.Config{
			Message: "Click Me",
			Layout:  box.Zeroed(),
			BackgroundColor: color.RGBA{
				R: 50,
				A: 255,
			},
			Padding: label.Padding{
				Left:   10,
				Right:  10,
				Top:    10,
				Bottom: 10,
			},
		},
	})
	// p.Children.Add(exampleButton)

	exampleButton2 := button.New(button.Config{
		OnClick: func() {
			slog.Info("No Click Me button clicked")
		},
		Gesturer: p.Gesturer,
		Config: label.Config{
			Message: "No Click Me",
			Layout:  box.Zeroed(),
			BackgroundColor: color.RGBA{
				R: 50,
				A: 255,
			},
			Padding: label.Padding{
				Left:   10,
				Right:  10,
				Top:    10,
				Bottom: 10,
			},
		},
	})

	// mainPanel.Children.Add(mainStack)

	mainStackChildren := &doodad.Children{
		Doodads: []doodad.Doodad{
			titleLabel,
			exampleButton,
			exampleButton2,
		},
	}

	mainStack := stack.New(stack.Config{
		Type: stack.Vertical,
		Layout: box.Computed(func(b *box.Box) *box.Box {
			boundingBox := box.Bounding(mainStackChildren.Boxes())
			return b.CopyDimensionsOf(boundingBox).CenterWithin(p.Box)
		}),
		Children:     mainStackChildren,
		SpaceBetween: 10,
	})

	// mainPanel.Children.Add(mainStack)
	// mainPanel.Layout().AddDependent(mainStack.Layout())

	// mainStack.Layout().Computed(func(b *box.Box) *box.Box {
	// 	return b.CenterWithin(p.Box)
	// })

	// mainStack.Children.DependOnLayout(mainStack.Layout())
	// mainStack.Layout().Recalculate()

	// mainPanel := panel.New(panel.Config{
	// 	Gesturer: p.Gesturer,
	// 	Layout: box.Computed(func(b *box.Box) *box.Box {
	// 		return b.SetDimensions(400, 300)
	// 	}),
	// 	Children: &doodad.Children{
	// 		Doodads: []doodad.Doodad{
	// 			mainStack,
	// 		},
	// 	},
	// })

	// mainPanel.Layout().Computed(func(b *box.Box) *box.Box {
	// 	return b.CenterWithin(p.Box)
	// })

	// p.Children.Add(mainPanel)

	p.AddChild(mainStack)
	p.Children.Setup()
}
