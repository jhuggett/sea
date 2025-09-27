package bottom_bar

import (
	"design-library/button"
	"design-library/doodad"
	"design-library/label"
	"design-library/position/box"
	"design-library/stack"

	"github.com/jhuggett/frontend/colors"
	"github.com/jhuggett/frontend/game"
)

func NewBottomBar(manager *game.Manager) *BottomBar {
	bottomBar := &BottomBar{
		Manager: manager,
	}

	return bottomBar
}

type BottomBar struct {
	doodad.Default

	Manager *game.Manager

	panel *Panel
}

func (b *BottomBar) Setup() {

	var activePanel string

	shipInfoPanelButton := button.New(button.Config{
		OnClick: func(bb *button.Button) {
			if b.panel.IsVisible() && activePanel == "ship-info" {
				b.panel.Hide()
				bb.SetMessage("Ship Info")
				activePanel = ""
			} else {
				b.panel.Show()
				b.panel.SetContents([]doodad.Doodad{
					NewShipInfoPanel(b.Manager),
				})
				bb.SetMessage("[Ship Info]")
				activePanel = "ship-info"
			}
		},
		Config: label.Config{
			Message: "Ship Info",
		},
	})

	shipInventoryPanelButton := button.New(button.Config{
		OnClick: func(bb *button.Button) {
			if b.panel.IsVisible() && activePanel == "inventory" {
				b.panel.Hide()
				bb.SetMessage("Inventory")
				activePanel = ""
			} else {
				b.panel.Show()
				b.panel.SetContents([]doodad.Doodad{
					NewShipInventoryPanel(b.Manager),
				})
				bb.SetMessage("[Inventory]")
				activePanel = "inventory"
			}
		},
		Config: label.Config{
			Message: "Inventory",
		},
	})

	mainStack := stack.New(stack.Config{
		Type: stack.Horizontal,
		Children: doodad.NewChildren(
			b,
			[]doodad.Doodad{
				shipInfoPanelButton,
				shipInventoryPanelButton,
			},
		),
		Layout: box.Computed(func(bb *box.Box) {
			bb.SetWidth(b.Parent().Layout().Width()).SetHeight(50).SetY(b.Parent().Layout().Height() - 50)
		}),
		BackgroundColor: colors.SemiTransparent(colors.Panel),
	})

	b.AddChild(mainStack)

	b.panel = NewPanel()
	b.AddChild(b.panel)

	b.Children().Setup()

	b.panel.Layout().Computed(func(b *box.Box) {
		b.SetWidth(450).SetHeight(250).MoveAbove(mainStack.Box).MoveUp(20)
	})

	b.panel.Hide()
}
