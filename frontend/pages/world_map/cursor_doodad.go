package world_map

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jhuggett/frontend/doodad"
)

type CursorDoodad struct {
	MouseX, MouseY int

	SpaceTranslator
	Gesturer doodad.Gesturer

	img *ebiten.Image

	Hidden bool
}

func (w *CursorDoodad) Update() error {
	return nil
}

func (w *CursorDoodad) Draw(screen *ebiten.Image) {
	if w.Hidden {
		return
	}

	op := &ebiten.DrawImageOptions{}

	x, y := w.SpaceTranslator.FromScreenToWorld(
		float64(w.MouseX),
		float64(w.MouseY),
	)

	x, y = w.SpaceTranslator.FromWorldToData(x, y)

	x, y = Floor(x, y)

	x, y = w.SpaceTranslator.FromDataToWorld(x, y)

	x, y = w.SpaceTranslator.FromWorldToScreen(x, y)

	xScale, yScale := w.SpaceTranslator.ScreenScale()

	op.GeoM.Translate(x, y)
	op.GeoM.Scale(xScale, yScale)

	screen.DrawImage(w.img, op)
}

func (w *CursorDoodad) Setup() error {

	width, height := w.SpaceTranslator.TileSize()
	w.img = ebiten.NewImage(int(width), int(height))
	w.img.Fill(color.Black)

	w.Gesturer.OnMouseMove(func(x, y int) error {
		w.MouseX = x
		w.MouseY = y

		return nil
	})

	return nil
}
