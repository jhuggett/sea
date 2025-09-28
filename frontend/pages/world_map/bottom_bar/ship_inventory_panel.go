package bottom_bar

import (
	"design-library/doodad"
	"design-library/label"
	"design-library/position/box"
	"design-library/stack"
	"fmt"

	"github.com/jhuggett/frontend/colors"
	"github.com/jhuggett/frontend/game"
	"github.com/jhuggett/sea/outbound"
)

func NewShipInventoryPanel(manager *game.Manager) *ShipInventoryPanel {
	shipInventoryPanel := &ShipInventoryPanel{
		Manager: manager,
	}
	return shipInventoryPanel
}

type ShipInventoryPanel struct {
	doodad.Default

	Manager *game.Manager

	InventoryData outbound.ShipInventoryChangedReq

	Subbed string
}

func (s *ShipInventoryPanel) Setup() {
	// Subscribe to inventory updates if not already subscribed
	if s.Subbed == "" {
		// Add callback for ship inventory changes
		s.Subbed = s.Manager.OnShipInventoryChangedCallback.Add(func(sicr outbound.ShipInventoryChangedReq) error {
			s.InventoryData = sicr
			doodad.ReSetup(s)
			return nil
		})

	}

	// Main container stack
	mainStack := stack.New(stack.Config{
		FitContents: true,
		Type:        stack.Vertical,
	})
	s.AddChild(mainStack)

	// Title
	mainStack.AddChild(
		label.New(label.Config{
			Message: "Ship Inventory",
		}),
	)

	// If no inventory data yet, show a message
	if s.InventoryData.ShipID == 0 || len(s.InventoryData.Inventory.Items) == 0 {
		mainStack.AddChild(
			label.New(label.Config{
				Message: "No inventory data available yet",
			}),
		)
		s.Children().Setup()

		go s.Manager.PlayerShip.TriggerShipInventoryRequest()
		return
	}

	// Display inventory ID
	mainStack.AddChild(
		label.New(label.Config{
			Message: fmt.Sprintf("Inventory ID: %d", s.InventoryData.Inventory.ID),
		}),
	)

	// Create header for items table
	headerStack := stack.New(stack.Config{
		Type:            stack.Horizontal,
		BackgroundColor: colors.SemiTransparent(colors.Panel),
		FitContents:     true,
		SpaceBetween:    10,
		Padding: stack.Padding{
			Top: 15,
		},
	})
	mainStack.AddChild(headerStack)

	// Add headers
	headerStack.AddChild(
		label.New(label.Config{
			Message: "Name",
			Layout:  box.Zeroed().SetWidth(200),
		}),
	)

	headerStack.AddChild(
		label.New(label.Config{
			Message: "Amount",
			Layout:  box.Zeroed().SetWidth(100),
		}),
	)

	// Create a scrollable list for items
	itemsContainer := stack.New(stack.Config{
		Type:        stack.Vertical,
		FitContents: true,
	})
	mainStack.AddChild(itemsContainer)

	// Add each item to the list
	for _, item := range s.InventoryData.Inventory.Items {
		itemRow := stack.New(stack.Config{
			Type:         stack.Horizontal,
			FitContents:  true,
			SpaceBetween: 10,
		})
		itemsContainer.AddChild(itemRow)

		// Item name
		itemRow.AddChild(
			label.New(label.Config{
				Message: item.Name,
				Layout:  box.Zeroed().SetWidth(200),
			}),
		)

		// Item amount
		itemRow.AddChild(
			label.New(label.Config{
				Message: fmt.Sprintf("%.1f", item.Amount),
				Layout:  box.Zeroed().SetWidth(100),
			}),
		)
	}

	s.Children().Setup()

	s.Layout().Recalculate()
}

func (s *ShipInventoryPanel) Teardown() error {

	if s.Subbed != "" {
		s.Manager.OnShipInventoryChangedCallback.Remove(s.Subbed)
		s.Subbed = ""
	}

	s.Children().Teardown()

	return nil
}
