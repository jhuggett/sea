package world_map

import (
	"design-library/doodad"
	"design-library/label"
	"design-library/position/box"
	"design-library/stack"

	"github.com/jhuggett/frontend/colors"
)

func NewBottomBar(
	doodadGesturer doodad.Gesturer,
	layout *box.Box,
) *BottomBar {
	bottomBar := &BottomBar{
		Default: *doodad.NewDefault(),
	}

	bottomBar.Gesturer = doodadGesturer
	if layout != nil {
		bottomBar.Box = layout
	}

	return bottomBar
}

type BottomBar struct {
	doodad.Default
}

func (b *BottomBar) Setup() {

	exampleLabel := label.New(label.Config{
		Message: "Bottom Bar",
	})

	mainStack := stack.New(stack.Config{
		Type: stack.Horizontal,
		Children: &doodad.Children{
			Doodads: []doodad.Doodad{
				exampleLabel,
			},
		},
		Layout: box.Computed(func(bb *box.Box) *box.Box {
			return bb.SetWidth(b.Box.Width()).SetHeight(50).SetY(b.Box.Height() - 50)
		}),
		BackgroundColor: colors.SemiTransparent(colors.Panel),
	})

	b.AddChild(mainStack)
	b.Children.Setup()
}
