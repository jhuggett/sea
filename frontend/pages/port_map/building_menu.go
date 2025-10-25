package port_map

import (
	"design-library/doodad"
	"design-library/label"
	"design-library/stack"

	"github.com/jhuggett/frontend/colors"
	"github.com/jhuggett/frontend/game"
)

func NewBuildingMenu() *BuildingMenu {
	return &BuildingMenu{}
}

type BuildingMenu struct {
	doodad.Default

	Building *game.Building
}

func (b *BuildingMenu) Setup() {
	if b.Building == nil {
		return
	}

	panel := stack.New(stack.Config{
		BackgroundColor: colors.Panel,
		Padding:         stack.Padding{10, 10, 10, 10},
		FitContents:     true,
	})
	b.AddChild(panel)

	panel.AddChild(
		label.New(label.Config{
			Message: b.Building.Name,
		}),
	)

	b.Children().Setup()
}
