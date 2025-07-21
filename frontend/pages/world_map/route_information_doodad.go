package world_map

import (
	"design-library/button"
	"design-library/doodad"
	"design-library/label"
	"design-library/position/box"
	"design-library/stack"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jhuggett/frontend/colors"
	"github.com/jhuggett/frontend/game"
	"github.com/jhuggett/sea/outbound"
)

func NewRouteInformationDoodad(
	ship *game.Ship,
	gesturer doodad.Gesturer,
	positioner func(*box.Box) *box.Box,
) *RouteInformationDoodad {
	doodad := &RouteInformationDoodad{
		Ship:      ship,
		Default:   *doodad.NewDefault(),
		Postioner: positioner,
	}

	doodad.Gesturer = gesturer

	ship.Manager.OnShipMovedCallback.Add(func(smr outbound.ShipMovedReq) error {

		return nil
	})

	return doodad
}

type RouteInformationDoodad struct {
	SpaceTranslator SpaceTranslator

	Postioner func(*box.Box) *box.Box

	Ship *game.Ship

	doodad.Default
}

func (w *RouteInformationDoodad) Setup() {
	// setSailButton, err := button.New(button.Config{
	// 	Message: "Set Sail",
	// 	OnClick: func() {
	// 		_, err := w.Ship.SetSail()
	// 		if err != nil {
	// 			slog.Error("Failed to set sail", "error", err)
	// 			return
	// 		}
	// 	},
	// 	Gesturer: w.Gesturer,
	// 	Position: doodad.ZeroZero,
	// })
	// if err != nil {
	// 	slog.Error("Failed to create Set Sail button", "error", err)
	// }
	// w.Children.Add(setSailButton)

	panelChildren := doodad.NewChildren(
		[]doodad.Doodad{
			label.New(label.Config{
				Message: "Route Information",
			}),
			button.New(button.Config{
				OnClick: func() {
					_, err := w.Ship.SetSail()
					if err != nil {
						slog.Error("Failed to set sail", "error", err)
						return
					}
				},
				Gesturer: w.Gesturer,
				Config: label.Config{
					Message: "Set Sail",
				},
			}),
		},
	)

	panel := stack.New(stack.Config{
		Children: panelChildren,
		Layout: box.Computed(func(b *box.Box) *box.Box {
			boundingBox := box.Bounding(panelChildren.Boxes())
			return w.Postioner(b.CopyDimensionsOf(boundingBox))
		}),
		Padding: stack.Padding{
			Top:    10,
			Right:  10,
			Bottom: 10,
			Left:   10,
		},
		SpaceBetween:    10,
		Type:            stack.Horizontal,
		BackgroundColor: colors.Panel,
	})

	w.AddChild(panel)

	w.Children.Setup()
}

func (w *RouteInformationDoodad) Draw(screen *ebiten.Image) {
	if !w.Ship.HasRoute() {
		return
	}

	w.Children.Draw(screen)
}
