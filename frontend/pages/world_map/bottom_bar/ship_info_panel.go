package bottom_bar

import (
	"design-library/button"
	"design-library/doodad"
	"design-library/label"
	"design-library/position/box"
	"design-library/stack"
	"fmt"

	"github.com/jhuggett/frontend/game"
	"github.com/jhuggett/sea/outbound"
)

func NewShipInfoPanel(manager *game.Manager) *ShipInfoPanel {
	shipInfoPanel := &ShipInfoPanel{
		Manager: manager,
	}
	return shipInfoPanel
}

type ShipInfoPanel struct {
	doodad.Default

	Manager *game.Manager

	ShipData outbound.ShipChangedReq
}

func (s *ShipInfoPanel) Setup() {

	if s.ShipData.ID == 0 {
		return // No ship data available yet
	}

	listStack := stack.New(stack.Config{
		Type: stack.Vertical,
	})

	s.AddChild(listStack)

	// Ship basic info
	listStack.AddChild(
		label.New(label.Config{
			Message: "Ship Info Panel",
		}),
	)

	listStack.AddChild(
		label.New(label.Config{
			Message: "Name: " + string(s.ShipData.ID),
		}),
	)

	// Position and status
	listStack.AddChild(
		label.New(label.Config{
			Message: "Position: X=" + fmt.Sprintf("%.2f", s.ShipData.X) + ", Y=" + fmt.Sprintf("%.2f", s.ShipData.Y),
		}),
	)

	dockedStatus := "No"
	if s.ShipData.IsDocked {
		dockedStatus = "Yes"
	}
	listStack.AddChild(
		label.New(label.Config{
			Message: "Docked: " + dockedStatus,
		}),
	)

	// Ship condition and performance
	listStack.AddChild(
		label.New(label.Config{
			Message: "State of Repair: " + fmt.Sprintf("%.1f%%", s.ShipData.StateOfRepair*100),
		}),
	)

	listStack.AddChild(
		label.New(label.Config{
			Message: "Est. Sailing Speed: " + fmt.Sprintf("%.2f knots", s.ShipData.EstimatedSailingSpeed),
		}),
	)

	// Manning info
	listStack.AddChild(
		label.New(label.Config{
			Message: "Manning: Min " + fmt.Sprintf("%d", s.ShipData.MinimumSafeManning) +
				", Max " + fmt.Sprintf("%d", s.ShipData.MaximumSafeManning),
		}),
	)

	// Cargo info
	listStack.AddChild(
		label.New(label.Config{
			Message: "Cargo: " + fmt.Sprintf("%.1f/%.0f weight, %.1f/%.0f space",
				s.ShipData.CurrentCargoWeight, float32(s.ShipData.RecommendedMaxCargoWeightCapacity),
				s.ShipData.CurrentCargoSpace, float32(s.ShipData.MaxCargoSpaceCapacity)),
		}),
	)

	// Fleet info
	fleetCount := len(s.ShipData.Fleet)
	listStack.AddChild(
		label.New(label.Config{
			Message: "Fleet Size: " + fmt.Sprintf("%d ship(s)", fleetCount),
		}),
	)

	repairButton := button.New(button.Config{
		OnClick: func(b *button.Button) {
			s.Manager.PlayerShip.Repair()
		},
		Config: label.Config{
			Message: "Repair Ship",
		},
	})
	s.AddChild(repairButton)

	s.Children().Setup()

	repairButton.Layout().Computed(func(b *box.Box) {
		b.AlignBottomWithin(s.Box).AlignRight(s.Box)
	})

	// TODO: handle unsub
	s.Manager.OnShipChangedCallback.Add(func(scr outbound.ShipChangedReq) error {
		s.ShipData = scr

		doodad.ReSetup(s)

		return nil
	})
}
