package world_map

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jhuggett/frontend/game"
)

type ShipDoodad struct {
	Ship *game.Ship

	img *ebiten.Image

	Origin func() (int, int)
	Scale  func() (float64, float64)
}

func (w *ShipDoodad) Update() error {
	return nil
}

func (w *ShipDoodad) Draw(screen *ebiten.Image) {
	originX, originY := w.Origin()
	scaleX, scaleY := w.Scale()

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(
		float64((w.Ship.Location().X)*w.Ship.Manager.TileSize()-originX),
		float64((w.Ship.Location().Y)*w.Ship.Manager.TileSize()-originY),
	)

	op.GeoM.Scale(
		scaleX,
		scaleY,
	)
	screen.DrawImage(w.img, op)
}

func (w *ShipDoodad) Setup() error {

	w.img = ebiten.NewImage(w.Ship.Manager.TileSize(), w.Ship.Manager.TileSize())
	w.img.Fill(color.RGBA{
		R: 0,
		G: 0,
		B: 255,
		A: 255,
	})

	return nil
}
