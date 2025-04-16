package world_map

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jhuggett/frontend/doodad"
	"github.com/jhuggett/frontend/game"
)

type PortDoodad struct {
	Port *game.Port

	img *ebiten.Image

	SpaceTranslator SpaceTranslator

	Doodads []doodad.Doodad
}

func (w *PortDoodad) Update() error {
	return nil
}

func (w *PortDoodad) Draw(screen *ebiten.Image) {

	// originX, originY := w.Origin()
	// scaleX, scaleY := w.Scale()

	// op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(
	// 	float64((w.Port.Location().X)*w.Port.Manager.TileSize()-originX),
	// 	float64((w.Port.Location().Y)*w.Port.Manager.TileSize()-originY),
	// )

	// op.GeoM.Scale(
	// 	scaleX,
	// 	scaleY,
	// )

	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(
		w.SpaceTranslator.FromWorldToScreen(
			w.SpaceTranslator.FromDataToWorld(
				float64(w.Port.Location().X),
				float64(w.Port.Location().Y),
			),
		),
	)

	op.GeoM.Scale(
		w.SpaceTranslator.ScreenScale(),
	)

	screen.DrawImage(w.img, op)
}

func (w *PortDoodad) Setup() error {

	w.img = ebiten.NewImage(
		w.Port.Manager.TileSize(),
		w.Port.Manager.TileSize(),
	)

	w.img.Fill(color.RGBA{
		R: 255,
		G: 0,
		B: 0,
		A: 255,
	})

	return nil
}
