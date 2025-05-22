package world_map

import (
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jhuggett/frontend/doodad"
)

type TileInformationDoodad struct {
	SpaceTranslator
	Gesturer doodad.Gesturer

	Doodads []doodad.Doodad

	Hidden bool

	HideTileCursor func()
	ShowTileCursor func()
}

func (w *TileInformationDoodad) Update() error {
	return nil
}

func (w *TileInformationDoodad) Draw(screen *ebiten.Image) {
	if w.Hidden {
		return
	}

	for _, doodad := range w.Doodads {
		doodad.Draw(screen)
	}
}

func (w *TileInformationDoodad) Setup() error {

	w.Hidden = true

	infoPanel := doodad.NewPanel(w.Gesturer)
	infoPanel.SetPosition(func() doodad.Position {
		return doodad.Position{
			X: 0,
			Y: 0,
		}
	})
	w.Doodads = append(w.Doodads, infoPanel)

	w.Gesturer.OnMouseUp(func(event doodad.MouseUpEvent) error {
		if !w.Hidden {
			return doodad.ErrStopPropagation
		}

		if event.Button != ebiten.MouseButtonRight {
			return nil
		}

		x, y := w.SpaceTranslator.FromWorldToData(w.SpaceTranslator.FromScreenToWorld(float64(event.X), float64(event.Y)))

		slog.Debug("WorldMapDoodad.OnClick", "x", x, "y", y)

		infoPanel.SetPosition(func() doodad.Position {
			panelX, panelY := w.SpaceTranslator.FromWorldToScreen(w.SpaceTranslator.FromDataToWorld(x, y))
			scaleX, scaleY := w.SpaceTranslator.ScreenScale()

			return doodad.Position{
				X: int(panelX * scaleX),
				Y: int(panelY * scaleY),
			}
		})

		w.Hidden = false
		w.HideTileCursor()

		return nil
	})

	w.Gesturer.OnMouseMove(func(x, y int) error {
		if !w.Hidden {
			return doodad.ErrStopPropagation
		}
		return nil
	})

	closeButton := doodad.NewButton(
		"X",
		func() {
			slog.Debug("WorldMapDoodad.CloseButton.OnClick")
			w.Hidden = true
			w.ShowTileCursor()
		},
		w.Gesturer,
	)
	closeButton.Setup()
	closeButton.SetPosition(func() doodad.Position {
		return doodad.Position{
			X: infoPanel.Position().X + 10,
			Y: infoPanel.Position().Y + 10,
		}
	})
	w.Doodads = append(w.Doodads, closeButton)

	return nil
}
