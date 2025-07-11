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
) *BottomBar {
	bottomBar := &BottomBar{
		Default: *doodad.NewDefault(),
	}

	bottomBar.Gesturer = doodadGesturer

	return bottomBar
}

type BottomBar struct {
	doodad.Default
}

func (b *BottomBar) Setup() {
	b.Box.Computed(func(bb *box.Box) *box.Box {
		return bb.SetWidth(b.Parent().Layout().Width()).SetHeight(50).SetY(b.Parent().Layout().Height() - 50)
	})

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
			return bb.Copy(b.Box)
		}),
		BackgroundColor: colors.SemiTransparent(colors.Panel),
	})

	b.AddChild(mainStack)
	b.Children.Setup()
}
