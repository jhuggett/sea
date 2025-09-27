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

	Subbed bool
}

func (s *ShipInventoryPanel) Setup() {
	// Subscribe to inventory updates if not already subscribed
	if !s.Subbed {
		// Add callback for ship inventory changes
		s.Manager.OnShipInventoryChangedCallback.Add(func(sicr outbound.ShipInventoryChangedReq) error {
			s.InventoryData = sicr
			doodad.ReSetup(s)
			return nil
		})
		s.Subbed = true
	}

	// Main container stack
	mainStack := stack.New(stack.Config{
		Type: stack.Vertical,
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
	})
	mainStack.AddChild(headerStack)

	// Add headers
	headerStack.AddChild(
		label.New(label.Config{
			Message: "Item Name",
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
		Type: stack.Vertical,
	})
	mainStack.AddChild(itemsContainer)

	// Add each item to the list
	for _, item := range s.InventoryData.Inventory.Items {
		itemRow := stack.New(stack.Config{
			Type: stack.Horizontal,
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
}
