package port_map

import (
	"design-library/button"
	"design-library/config"
	"design-library/doodad"
	"design-library/label"
	"design-library/position/box"
	"design-library/stack"

	"github.com/jhuggett/frontend/colors"
	"github.com/jhuggett/frontend/game"
)

var buildingMenuConfigurations = map[string]func(*BuildingMenu, *stack.Stack){
	"tavern": func(b *BuildingMenu, s *stack.Stack) {
		s.AddChild(
			label.New(label.Config{
				Message: "Tavern Menu Placeholder",
			}),
		)

		hireCrewModal := NewHireCrewModal(b.Building)
		b.AddChild(hireCrewModal)

		s.AddChild(button.New(button.Config{
			OnClick: func(but *button.Button) {
				hireCrewModal.Show()
			},
			Config: label.Config{
				Message: "Hire Crew",
				Layout:  box.Zeroed(),
			},
		}))
	},
	"market": func(b *BuildingMenu, s *stack.Stack) {
		s.AddChild(
			label.New(label.Config{
				Message: "Market Menu Placeholder",
			}),
		)
	},
	"shipyard": func(b *BuildingMenu, s *stack.Stack) {
		s.AddChild(
			label.New(label.Config{
				Message: "Shipyard Menu Placeholder",
			}),
		)
	},
}

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

	b.Layout().Computed(func(b *box.Box) {
		b.SetX(20)
		b.SetY(40)
	})

	panel := stack.New(stack.Config{
		BackgroundColor: colors.Panel,
		Padding:         config.Padding{10, 10, 10, 10},
	})
	b.AddChild(panel)

	panel.AddChild(
		label.New(label.Config{
			Message: b.Building.Name,
		}),
	)

	if configureFunc, ok := buildingMenuConfigurations[b.Building.Type]; ok {
		configureFunc(b, panel)
	} else {
		panic("no building menu configuration for building type: " + b.Building.Type)
	}

	b.Children().Setup()
}
