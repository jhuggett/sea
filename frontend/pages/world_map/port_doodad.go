package world_map

import (
	"design-library/doodad"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jhuggett/frontend/game"
)

func NewPortDoodad(
	port *game.Port,
	spaceTranslator SpaceTranslator,
) *PortDoodad {
	portDoodad := &PortDoodad{
		Port:            port,
		SpaceTranslator: spaceTranslator,
		Default:         *doodad.NewDefault(),
	}

	return portDoodad
}

type PortDoodad struct {
	Port *game.Port

	img *ebiten.Image

	SpaceTranslator SpaceTranslator

	doodad.Default
}

func (w *PortDoodad) Teardown() error {
	return nil
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

func (w *PortDoodad) Setup() {

	w.img = ebiten.NewImage(
		w.Port.Manager.TileSize(),
		w.Port.Manager.TileSize(),
	)

	// Draw a red square to represent the port
	squareColor := color.RGBA{
		R: 255,
		G: 0,
		B: 0,
		A: 255, // Red color for the square
	}

	tileSize := w.Port.Manager.TileSize()
	centerX := tileSize / 2
	centerY := tileSize / 2

	// Draw the square
	squareSize := tileSize / 2 // Slightly smaller than the tile size
	for x := centerX - squareSize/2; x <= centerX+squareSize/2; x++ {
		for y := centerY - squareSize/2; y <= centerY+squareSize/2; y++ {
			w.img.Set(x, y, squareColor)
		}
	}
}
