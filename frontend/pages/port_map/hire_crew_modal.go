package port_map

import (
	"design-library/button"
	"design-library/doodad"
	"design-library/label"
	"design-library/position/box"
	"design-library/stack"

	"github.com/jhuggett/frontend/colors"
	"github.com/jhuggett/frontend/game"
)

func NewHireCrewModal(
	building *game.Building,
) *HireCrewModal {

	// Fetch potential crew members from the port
	potentialCrew, err := building.GetHireablePeople()
	if err != nil {
		panic(err)
	}

	return &HireCrewModal{
		Building:      building,
		PotentialCrew: potentialCrew,
	}
}

type HireCrewModal struct {
	Building *game.Building
	doodad.Default

	PotentialCrew []*game.HireablePerson
}

func (h *HireCrewModal) Setup() {

	panel := stack.New(stack.Config{
		BackgroundColor: colors.Panel,
		Padding:         stack.Padding{10, 10, 10, 10},
		FitContents:     true,
		SpaceBetween:    20,
	})

	h.AddChild(panel)

	panel.AddChild(label.New(label.Config{
		Message:  "Hire Crew Members",
		FontSize: 20,
		Padding:  label.Padding{10, 0, 20, 0},
	}))

	for _, crewMember := range h.PotentialCrew {
		row := stack.New(stack.Config{
			Type:         stack.Horizontal,
			FitContents:  true,
			SpaceBetween: 10,
		})

		panel.AddChild(row)

		nameLabel := label.New(label.Config{
			Message: crewMember.FirstName + " " + crewMember.LastName,
		})
		row.AddChild(nameLabel)

		hireButton := button.New(button.Config{
			OnClick: func(b *button.Button) {
				err := h.Building.HireCrewMember(crewMember)
				if err != nil {
					panic(err)
				}
			},
			Config: label.Config{
				Message: "Hire",
			},
		})
		row.AddChild(hireButton)

	}

	closeButton := button.New(button.Config{
		OnClick: func(b *button.Button) {
			h.Hide()
		},
		Config: label.Config{
			Message: "Close",
		},
	})
	panel.AddChild(closeButton)

	h.Children().Setup()

	panel.Layout().Computed(func(b *box.Box) {
		b.SetWidth(400).CenterWithin(h.Box)
	})

	closeButton.Layout().Computed(func(b *box.Box) {
		b.AlignTopWithin(panel.Box).AlignRightWithin(panel.Box)
	})

	h.Box.Recalculate()
}
