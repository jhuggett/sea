package world_map

import (
	"design-library/doodad"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jhuggett/frontend/game"
)

type RouteInformationDoodad struct {
	SpaceTranslator SpaceTranslator
	Gesturer        doodad.Gesturer
	Ship            *game.Ship

	Children doodad.Children
}

func (w *RouteInformationDoodad) Teardown() error {
	w.Children.Teardown()
	return nil
}

func (w *RouteInformationDoodad) Update() error {
	return nil
}

func (w *RouteInformationDoodad) Draw(screen *ebiten.Image) {
	w.Children.Draw(screen)
}

func (w *RouteInformationDoodad) Setup() error {
	// setSailButton, err := button.New(button.Config{
	// 	Message: "Set Sail",
	// 	OnClick: func() {
	// 		slog.Debug("Set Sail button clicked")
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

	return nil
}
