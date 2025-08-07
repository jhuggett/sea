package world_map

import (
	"design-library/doodad"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/jhuggett/frontend/game"
)

type ContinentDoodad struct {
	Continent *game.Continent

	doodad.Default

	LargestPointX int
	LargestPointY int

	SmallestPointX int
	SmallestPointY int

	SpaceTranslator SpaceTranslator

	Image *ebiten.Image
}

func (w *ContinentDoodad) Draw(screen *ebiten.Image) {
	// originX, originY := w.Origin()
	// scaleX, scaleY := w.Scale()

	// // Draw the continent
	// op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(
	// 	float64((w.SmallestPointX*w.Continent.TileSize())-originX),
	// 	float64((w.SmallestPointY*w.Continent.TileSize())-originY),
	// )
	// op.GeoM.Scale(
	// 	scaleX,
	// 	scaleY,
	// )
	// screen.DrawImage(w.Image, op)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(
		w.SpaceTranslator.FromWorldToScreen(
			w.SpaceTranslator.FromDataToWorld(
				float64(w.SmallestPointX),
				float64(w.SmallestPointY),
			),
		),
	)
	op.GeoM.Scale(
		w.SpaceTranslator.ScreenScale(),
	)
	screen.DrawImage(w.Image, op)

	w.Children().Draw(screen)
}

func (w *ContinentDoodad) Setup() {

	for _, point := range w.Continent.Points() {
		if point.X < w.SmallestPointX {
			w.SmallestPointX = point.X
		}
		if point.Y < w.SmallestPointY {
			w.SmallestPointY = point.Y
		}
		if point.X > w.LargestPointX {
			w.LargestPointX = point.X
		}
		if point.Y > w.LargestPointY {
			w.LargestPointY = point.Y
		}
	}

	a := uint8(50)
	b := uint8(100)
	// lowestColor := color.RGBA{
	// 	R: 23 + a,
	// 	G: 18 + a,
	// 	B: 13 + a,
	// 	A: 255,
	// }
	// highestColor := color.RGBA{
	// 	R: 58 + b,
	// 	G: 45 + b,
	// 	B: 25 + b,
	// 	A: 255,
	// }

	lowestColor := color.RGBA{
		R: 50 + a,
		G: 50 + a,
		B: 50 + a,
		A: 255,
	}
	highestColor := color.RGBA{
		R: 100 + b,
		G: 100 + b,
		B: 100 + b,
		A: 255,
	}

	img := ebiten.NewImage(
		(w.LargestPointX-w.SmallestPointX+1)*int(w.Continent.TileSize()),
		(w.LargestPointY-w.SmallestPointY+1)*int(w.Continent.TileSize()),
	)
	for _, point := range w.Continent.Points() {
		e := (point.Elevation - .5) * 2
		vector.DrawFilledRect(
			img,
			float32((point.X-w.SmallestPointX)*w.Continent.TileSize()),
			float32((point.Y-w.SmallestPointY)*w.Continent.TileSize()),
			float32(w.Continent.TileSize()),
			float32(w.Continent.TileSize()),
			color.RGBA{
				R: uint8(float64(highestColor.R-lowestColor.R)*e) + lowestColor.R,
				G: uint8(float64(highestColor.G-lowestColor.G)*e) + lowestColor.G,
				B: uint8(float64(highestColor.B-lowestColor.B)*e) + lowestColor.B,
				A: 255,
			},
			false,
		)
	}

	w.Image = img

	// for _, port := range w.Continent.Ports {

	// 	portDoodad := &PortDoodad{
	// 		Port:            port,
	// 		SpaceTranslator: w.SpaceTranslator,
	// 	}

	// 	err := portDoodad.Setup()
	// 	if err != nil {
	// 		return fmt.Errorf("failed to setup port widget: %w", err)
	// 	}
	// 	w.Doodads = append(w.Doodads, portDoodad)
	// }

	for _, port := range w.Continent.Ports {
		w.AddChild(NewPortDoodad(
			port,
			w.SpaceTranslator,
		))
	}

	w.Children().Setup()
}
