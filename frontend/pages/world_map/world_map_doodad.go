package world_map

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jhuggett/frontend/doodad"
	"github.com/jhuggett/frontend/game"
)

type WorldMapDoodad struct {
	WorldMap *game.WorldMap

	Doodads []doodad.Doodad

	Background *ebiten.Image

	SpaceTranslator SpaceTranslator
	Gesturer        doodad.Gesturer

	SmallestContinentPointX int
	SmallestContinentPointY int
}

func (w *WorldMapDoodad) Update() error {
	// w.WorldMap.Update()
	return nil
}

func (w *WorldMapDoodad) Draw(screen *ebiten.Image) {

	// originX, originY := w.Origin()
	// scaleX, scaleY := w.Scale()

	// // Draw the background
	// op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(
	// 	float64(w.SmallestContinentPointX)*float64(w.WorldMap.Manager.TileSize())-float64(originX),
	// 	float64(w.SmallestContinentPointY)*float64(w.WorldMap.Manager.TileSize())-float64(originY),
	// )

	// op.GeoM.Scale(
	// 	scaleX,
	// 	scaleY,
	// )
	// screen.DrawImage(w.Background, op)

	// Draw the background
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(
		w.SpaceTranslator.FromWorldToScreen(
			w.SpaceTranslator.FromDataToWorld(
				float64(w.SmallestContinentPointX),
				float64(w.SmallestContinentPointY),
			),
		),
	)
	op.GeoM.Scale(
		w.SpaceTranslator.ScreenScale(),
	)
	screen.DrawImage(w.Background, op)

	// Draw the doodads
	for _, doodad := range w.Doodads {
		doodad.Draw(screen)
	}
}

func (w *WorldMapDoodad) Setup() error {

	if w.WorldMap == nil {
		return fmt.Errorf("world map is nil")
	}

	for _, continent := range w.WorldMap.Continents {
		continentDoodad := &ContinentDoodad{
			Continent:       continent,
			SpaceTranslator: w.SpaceTranslator,
		}

		err := continentDoodad.Setup()
		if err != nil {
			return fmt.Errorf("failed to setup continent widget: %w", err)
		}

		w.Doodads = append(w.Doodads, continentDoodad)
	}

	largestX := 01
	largestY := 0

	for _, continent := range w.WorldMap.Continents {
		if continent.LargestX > largestX {
			largestX = continent.LargestX
		}
		if continent.LargestY > largestY {
			largestY = continent.LargestY
		}
		if continent.SmallestX < w.SmallestContinentPointX {
			w.SmallestContinentPointX = continent.SmallestX
		}
		if continent.SmallestY < w.SmallestContinentPointY {
			w.SmallestContinentPointY = continent.SmallestY
		}
	}

	w.Background = ebiten.NewImage(
		int(float64(largestX+1)*float64(w.WorldMap.Manager.TileSize()))-int(float64(w.SmallestContinentPointX)*float64(w.WorldMap.Manager.TileSize())),
		int(float64(largestY+1)*float64(w.WorldMap.Manager.TileSize()))-int(float64(w.SmallestContinentPointY)*float64(w.WorldMap.Manager.TileSize())),
	)

	w.Background.Fill(color.RGBA{
		R: 220,
		G: 202,
		B: 127,
		A: 255,
	})

	return nil
}

/*

mapImg := ebiten.NewImage(
		int(float64(page.LargestPointX-page.SmallestPointX+1)*page.TileSize),
		int(float64(page.LargestPointY-page.SmallestPointY+1)*page.TileSize),
	)

	mapImg.Fill(color.RGBA{
		R: 220,
		G: 202,
		B: 127,
		A: 255,
	})

	page.MapImage = mapImg

*/
